package daemon

import (
	"db"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/sevlyar/go-daemon"
)

var (
	done        = make(chan bool)
	connections = []*db.Handler{}
)

// StartDaemon starts Daemon
func StartDaemon() {
	fmt.Println("Daemon started")
	fmt.Println("You can kill a daemon with help of pid file, that contains pid of process")
	fmt.Println("(Pid file is gonna be created after the app has been started)")
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	go registerMetricsFromDevices(1, 100)
	go monitorIncomingMetricsOfDevices(1, 100)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	for _, conn := range connections {
		conn.CloseConnection()
	}
	log.Println("Daemon process is killed")

}

// registerIncomingMetricsFromDevices регистрируюет приходящие метрики с устройств, диапазон который указывается явно
func registerMetricsFromDevices(from, to int) {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		log.Println("Registering metrics from devices...")
		// симулирую ситуацию отправки метрик с определенного устройства
		for ID := from; ID < to; ID++ {
			metrics := make(map[int]int)
			for i := 1; i <= 5; i++ {
				metrics[i] = random(10, 230)
			}
			log.Println("Incoming mettics:", ID, metrics)
			handler.UpdateMetrics(ID, metrics)
		}
		time.Sleep(5 * time.Second)
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
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

func getHandler() *db.Handler {
	handler, err := db.GetConnection(db.GetConfigurationFromProperties())
	if err != nil {
		panic(err)
	}

	return handler
}
