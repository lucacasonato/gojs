package main

import (
	"github.com/lucacasonato/gojs/coordinator"
	httphandler "github.com/lucacasonato/gojs/handlers/http"
)

func main() {
	err := coordinator.SetupWorkers(4, "test.js")
	if err != nil {
		panic(err.Error())
	}

	panic(httphandler.Start())
}
