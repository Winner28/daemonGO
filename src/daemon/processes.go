package daemon

import (
	"db"
	"encoding/json"
	"log"
	"math/rand"
	"model"
	"os"
	"time"
)

var (
	done          = make(chan bool)
	connections   = []*db.Handler{}
	deviceMetrics = make(map[int]model.DeviceMetricsRange)
)

func init() {
	for _, value := range getDeviceMetricsRangesFromProperties() {
		deviceMetrics[value.DeviceID] = value
	}
}

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

// registerIncomingMetricsFromDevices регистрируюет приходящие метрики с устройств.
// Диапазон интересующих устройств указывается явно
func registerMetricsFromDevices(from, to int) {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		log.Println("Registering metrics from devices...")
		// симулируется ситуацию отправки метрик с определенного устройства
		for ID := from; ID < to; ID++ {
			metrics := make(map[int]int)
			for i := 1; i <= 5; i++ {
				metrics[i] = getDeviceMetricsForDevice(ID)
			}
			log.Println("Incoming mettics:", ID, metrics)
			handler.UpdateMetrics(ID, metrics)
		}
		time.Sleep(5 * time.Second)
	}
}

// monitorIncomingMetricsOfDevices мониторит incoming метрики
func monitorIncomingMetricsOfDevices(from, to int) {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		time.Sleep(3 * time.Second)
		handler.MonitorMetrics(from, to)
	}
}

func getDeviceMetricsForDevice(ID int) int {
	devMetr := deviceMetrics[ID]
	return random(devMetr.Min, devMetr.Max)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
