package db

import (
	"encoding/json"
	"log"
	"model"
	"os"
)

var metrics map[int]model.DeviceMetrics

func init() {
	metrics = make(map[int]model.DeviceMetrics, 1000)
	for _, value := range parseMetricsJSON() {
		metrics[value.DeviceID] = value
	}
}

func parseMetricsJSON() []model.DeviceMetrics {

	file, err := os.Open("./config/config.metrics.json")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var devM []model.DeviceMetrics
	if err = decoder.Decode(&devM); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return devM
}

// UpdateMetrics s
func (handler *Handler) UpdateMetrics(ID int, metrics map[int]int) {
	//update metrick
}

// MonitorMetrics m
func (handler *Handler) MonitorMetrics(from int, to int) {
	log.Println("Monitoring metrics!")
	for i := from; i < to; i++ {
		log.Print("Monitoring metrics: ", i)
	}
}
