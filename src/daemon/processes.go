package daemon

import (
	"db"
	"log"
	"math/rand"
	"time"
)

var (
	done        = make(chan bool)
	connections = []*db.Handler{}
)

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

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getDeviceMetricsForDevice(ID int) int {
	return random(0, 0)
}

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigurationFromProperties())
	if err != nil {
		panic(err)
	}
	return handler
}
