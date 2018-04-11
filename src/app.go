package main

import "db"

func main() {
	//daemon.StartDaemon()
	handler := getHandler()

	handler.UpdateMetrics(1, map[int]int{
		1: 100,
		2: 200,
		3: 300,
		4: 400,
	})
	handler.MonitorMetrics(1, 1)
}

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
