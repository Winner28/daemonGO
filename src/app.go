package main

import (
	"db"
	"fmt"
)

func main() {
	//daemon.StartDaemon()
	handler := getHandler()

	handler.UpdateMetrics(2, map[int]int{
		1: 100,
		2: 200,
		3: 300,
		4: 2400,
		5: 3000,
	})
	handler.MonitorMetrics(1, 2)
	fmt.Println(handler.DeviceMetricExists(2))
}

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
