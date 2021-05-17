package main

import (
	"log"
	"net/http"
)

type ServerApplication struct {
	config *Configurations
}

type HandlerFunc func(*ServerApplication) func(http.ResponseWriter, *http.Request)

func CreateServer(config *Configurations) *ServerApplication {
	server := &ServerApplication{
		config: config,
	}

	return server
}

func UnhandledErrors(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (a *ServerApplication) Run() error {
	log.Printf("INFO: Server ready and listening on %s", a.config.Server.ConnectionString())

	return http.ListenAndServe(a.config.Server.ConnectionString(), nil)
}

func (a *ServerApplication) AddRoute(pattern string, handler HandlerFunc) {
	http.HandleFunc(pattern, handler(a))
}
