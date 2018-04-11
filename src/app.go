package main

import (
	"daemon"
	"db"
)

func main() {
	daemon.StartDaemon()

}

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
