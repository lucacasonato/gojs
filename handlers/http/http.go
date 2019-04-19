package httphandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/lucacasonato/gojs/coordinator"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

type request struct {
	Handler int            `json:"handler"`
	Request requestRequest `json:"request"`
}

type requestRequest struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type response struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

var responses = make(map[string]chan *response)
var responseLock = sync.RWMutex{}

func init() {
	coordinator.AddHandler("http", func(m *coordinator.Message) (*coordinator.Message, error) {
		switch m.Type {
		case "route":
			return nil, addRoute(m)
		case "handle":
			return nil, handleResponse(m)
		case "fetch":
			return fetch(m)
		default:
			return nil, fmt.Errorf("the namespace http does not have type %s", m.Type)
		}
	})
}

func addRoute(m *coordinator.Message) error {
	var data map[string]interface{}
	var ok bool
	if data, ok = m.Data.(map[string]interface{}); !ok {
		return fmt.Errorf("data supplied in registration message for http_route is not object")
	}

	var route string
	if route, ok = data["route"].(string); !ok {
		return fmt.Errorf("route supplied in registration message for http_route is not string")
	}

	var h float64
	if h, ok = data["handler"].(float64); !ok {
		return fmt.Errorf("handler supplied in registration message for http_route is not string")
	}

	handler := int(h)

	r.HandleFunc(route, handle(handler))

	return nil
}

func handle(handler int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		setResponse(id)

		request := &request{
			Handler: handler,
			Request: requestRequest{
				Headers: r.Header,
				Body:    string(body),
			},
		}

		err = coordinator.SendMessage(id, "http", "handle", request)
		if err != nil {
			panic(err)
		}

		resp := <-getResponse(id)
		clearResponse(id)

		for k, v := range resp.Headers {
			for _, s := range v {
				w.Header().Set(k, s)
			}
		}
		w.WriteHeader(resp.Status)
		w.Write([]byte(resp.Body))
	}
}

func handleResponse(m *coordinator.Message) error {
	var ok bool
	var resp map[string]interface{}
	if resp, ok = m.Data.(map[string]interface{}); !ok {
		return fmt.Errorf("data supplied in response message for http_response is not object")
	}

	var statusCode float64
	if statusCode, ok = resp["statusCode"].(float64); !ok {
		return fmt.Errorf("statusCode supplied in data of response message for http_response is not number")
	}

	var h map[string]interface{}
	if h, ok = resp["headers"].(map[string]interface{}); !ok {
		return fmt.Errorf("headers supplied in data of response message for http_response is not object")
	}

	var body string
	if body, ok = resp["body"].(string); !ok {
		return fmt.Errorf("body supplied in data of response message for http_response is not string")
	}

	var headers map[string][]string
	for k, v := range h {
		if headers[k], ok = v.([]string); !ok {
			return fmt.Errorf("header value supplied in header data of response message for http_response is not string[]")
		}
	}

	sendResponse(m.ID, &response{
		Status:  int(statusCode),
		Headers: headers,
		Body:    body,
	})

	return nil
}

// Start listenting for http requests
func Start() error {
	return http.ListenAndServe(":8080", r)
}

func getResponse(id string) chan *response {
	//fmt.Println("r", id)
	responseLock.RLock()
	defer responseLock.RUnlock()
	return responses[id]
}

func sendResponse(id string, resp *response) {
	//fmt.Println("w", id)
	responseLock.Lock()
	defer responseLock.Unlock()
	responses[id] <- resp
}

func setResponse(id string) {
	//fmt.Println("s", id)
	responseLock.Lock()
	defer responseLock.Unlock()
	responses[id] = make(chan *response, 1)
}

func clearResponse(id string) {
	//fmt.Println("c", id)
	responseLock.Lock()
	defer responseLock.Unlock()
	delete(responses, id)
}

/* ---------------------------------- */
/* HTTP Fetch                         */
/* ---------------------------------- */

type fetchInit struct {
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type fetchResponse struct {
	OK      bool                `json:"ok"`
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	URL     string              `json:"url"`
	Body    string              `json:"body"`
}

func fetch(m *coordinator.Message) (*coordinator.Message, error) {
	var ok bool
	var resp map[string]interface{}
	if resp, ok = m.Data.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("data supplied in response message for http_fetch is not object")
	}

	var resource string
	if resource, ok = resp["resource"].(string); !ok {
		return nil, fmt.Errorf("resource supplied in data of fetch request for http_fetch is not string")
	}

	var init fetchInit
	var initMap map[string]interface{}
	if initMap, ok = resp["init"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("init supplied in data of fetch request for http_fetch is not object")
	}

	if init.Body, ok = initMap["body"].(string); !ok {
		return nil, fmt.Errorf("body supplied in init of data of fetch request for http_fetch is not string")
	}

	var h map[string]interface{}
	if h, ok = initMap["headers"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("body supplied in init of data of fetch request for http_fetch is not fetchInit")
	}

	for k, v := range h {
		if init.Headers[k], ok = v.([]string); !ok {
			return nil, fmt.Errorf("header value supplied in header data of fetch request for http_fetch is not string[]")
		}
	}

	if init.Method, ok = initMap["method"].(string); !ok {
		return nil, fmt.Errorf("method supplied in init of data of fetch request for http_fetch is not fetchInit")
	}

	var b bytes.Buffer
	b.WriteString(init.Body)
	req, err := http.NewRequest(init.Method, resource, &b)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var location string
	if res.Header.Get("Location") != "" {
		loc, err := res.Location()
		if err != nil {
			return nil, err
		}

		location = loc.String()
	}

	return &coordinator.Message{
		ID:        "",
		Namespace: "http",
		Type:      "fetch",
		Data: fetchResponse{
			OK:      res.StatusCode >= 200 && res.StatusCode < 300,
			Status:  res.StatusCode,
			Headers: res.Header,
			URL:     location,
			Body:    string(body),
		},
	}, nil
}
