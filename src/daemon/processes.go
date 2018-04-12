package daemon

import (
	"encoding/json"
	"log"
	"math/rand"
	"model"
	"os"
	"service"
	"time"
)

var (
	// Contains an open connections to database. If our Daemon-app will be interrupted
	// We gonna run through connections and close each
	connections = []*service.Handler{}

	// Map with the next structure: [key = deviceMetrics.DeviceID, value = deviceMetricsRange]
	// Initialized in init method with a help of getDeviceMetricsRangesFromProperties() function
	// Contains DeviceID and [min:max] of metrics. (uses to simulate sending metrics process)
	deviceMetrics = make(map[int]model.DeviceMetricsRange)
)

func init() {
	for _, value := range getDeviceMetricsRangesFromProperties() {
		deviceMetrics[value.DeviceID] = value
	}
}

// This method fills up deviceMetrics that contains range ([min:max]) for Specified ID.
// Parses from a file that lies under "./config/config.metrics_range.json"
func getDeviceMetricsRangesFromProperties() []model.DeviceMetricsRange {
	file, err := os.Open("./config/config.metrics_range.json")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var devMetrics []model.DeviceMetricsRange
	if err = decoder.Decode(&devMetrics); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return devMetrics
}

// Designed to register device metrics in the database
// Each device has an optimal value for a certain metric (lies under /config/config.metrics.json)
// We simulate the sending of the metric from the device with a help of deviceMetrics map,
// that for [key = deviceID] gonna return deviceMetricsRange, that contains a min value,
// and a max value of metrics (values that lies under max is always bigger than the optimal metrics).
// In this range[min:max] we gonna generate specific metric and sent it to db.
func registerIncomingMetricsFromDevices(from, to int) {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		log.Println("Registering metrics from devices...")
		for ID := from; ID <= to; ID++ {
			metrics := make(map[int]int)
			for i := 1; i <= 5; i++ {
				metrics[i] = getDeviceMetricsForDevice(ID)
			}
			log.Println("Updating metrics:", ID, metrics)
			handler.UpdateMetrics(ID, metrics)
		}
		time.Sleep(5 * time.Second)
	}
}

// Rutine, that every 3 Seconds monitor metrics values of devices
func monitorIncomingMetricsOfDevices(from, to int) {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		time.Sleep(3 * time.Second)
		handler.MonitorMetrics(from, to)
	}
}

// Simulate getting device metrics for specified Device_ID
func getDeviceMetricsForDevice(ID int) int {
	devMetr := deviceMetrics[ID]
	return random(devMetr.Min, devMetr.Max)
}

// Help of sim :)
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// Function that returns Handler.
func getHandler() *service.Handler {
	handler, err := service.GetConnection(service.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
