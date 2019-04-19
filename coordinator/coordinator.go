package coordinator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/ry/v8worker2"
)

// Worker that executes JS code
type Worker struct {
	*v8worker2.Worker
}

// Workers is a collection of all registered workers
var Workers []*Worker

// Message to/from JS worker
type Message struct {
	ID        string      `json:"id"`
	Namespace string      `json:"namespace"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
}

// HandlerFunc deals with an incoming message
type HandlerFunc func(*Message) (*Message, error)

// handlers connects namespaces and handlers
var handlers = make(map[string]HandlerFunc)
var handlerLock = &sync.RWMutex{}

// SetupWorkers that can handle code
func SetupWorkers(count int, file string) error {
	for i := 0; i < count; i++ {
		worker := &Worker{
			Worker: v8worker2.New(handleWorker),
		}

		loaded := make(map[string]bool)

		err := loadCode(worker.Worker, file, file, loaded)
		if err != nil {
			log.Fatalln(err)
		}

		Workers = append(Workers, worker)
	}

	return nil
}

var regexpImportFrom = regexp.MustCompile(`(?m)import(?:["'\s]*([\w*{}\n, ]+)from\s*)?["'\s]*([@\w/_.-]+)["'\s].*`)

func loadCode(worker *v8worker2.Worker, file string, name string, loaded map[string]bool) error {
	fileToLoad := file
	if !strings.HasSuffix(fileToLoad, ".js") {
		fileToLoad += ".js"
	}

	js, err := ioutil.ReadFile(fileToLoad)
	if err != nil {
		return err
	}

	for _, match := range regexpImportFrom.FindAllStringSubmatch(string(js), -1) {
		path := match[2]
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(file), path)
		}

		if loaded[match[2]] == false {
			err = loadCode(worker, path, match[2], loaded)
			if err != nil {
				return err
			}

			loaded[match[2]] = true
		}
	}

	err = worker.LoadModule(name, string(js), func(moduleName, referrerName string) int {
		fmt.Println(moduleName, referrerName)
		return 0
	})
	if err != nil {
		return err
	}

	return nil
}

// GetWorker to do some processing with
func GetWorker() *Worker {
	if len(Workers) == 1 {
		return Workers[0]
	} else if len(Workers) > 1 {
		return Workers[rand.Intn(len(Workers)-1)]
	} else {
		log.Fatalln("no workers available")
	}

	return nil
}

// SendMessage to any worker
func SendMessage(id string, namespace string, typ string, data interface{}) error {
	req, err := json.Marshal(Message{
		ID:        id,
		Namespace: namespace,
		Type:      typ,
		Data:      data,
	})
	if err != nil {
		return err
	}

	err = GetWorker().SendBytes(req)
	if err != nil {
		return err
	}

	return nil
}

// AddHandler for a certain namespace
func AddHandler(namespace string, handler HandlerFunc) {
	handlerLock.Lock()
	defer handlerLock.Unlock()
	handlers[namespace] = handler
}

func handleWorker(message []byte) []byte {
	m := new(Message)
	err := json.Unmarshal(message, &m)
	if err != nil {
		log.Fatalln(err)
	}

	var handler HandlerFunc
	var ok bool
	handlerLock.RLock()
	if handler, ok = handlers[m.Namespace]; !ok {
		return (transformError(fmt.Errorf("no handler registered for namespace %s", m.Namespace)))
	}
	handlerLock.RUnlock()

	resp, err := handler(m)
	if err != nil {
		return (transformError(err))
	}

	d, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return d
}

type responseError struct {
	Error string `json:"error"`
}

func transformError(err error) []byte {
	d, err := json.Marshal(responseError{
		Error: err.Error(),
	})
	if err != nil {
		panic(err)
	}

	return d
}
