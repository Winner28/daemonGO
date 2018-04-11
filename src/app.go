package main

import (
	"daemon"
	"service"
)

func main() {
	daemon.StartDaemon()

}

func getHandler() *service.Handler {
	handler, err := service.GetConnection(service.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
