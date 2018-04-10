package daemon

import (
	"db"
	"fmt"
	"log"
	"time"

	"github.com/sevlyar/go-daemon"
)

var (
	values      = make(chan []int)
	done        = make(chan bool)
	connections = []*db.Handler{}
)

// StartDaemon starts Daemon
func StartDaemon() {
	fmt.Println("Daemon started")
	fmt.Println("You can kill a daemon with help of pid file, that contains pid of process")
	fmt.Println("(pid file is gonna be created after the app has been started)")
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

	go registerIncomingMetricsFromDevices()
	go monitorIncomingMetricsOfDevices()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	for _, conn := range connections {
		conn.CloseConnection()
	}
	log.Println("Daemon process is killed")

}

// registerIncomingMetricsFromDevices регистрируюет приходящие метрики с устройств
func registerIncomingMetricsFromDevices() {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		time.Sleep(5 * time.Second)
		//generatesRandomMetrick -> update Metrick
		done <- true
	}
}

// monitorIncomingMetricsOfDevices мониторит incoming метрики
func monitorIncomingMetricsOfDevices() {
	handler := getHandler()
	connections = append(connections, handler)
	for {
		if <-done {
			handler.ReadMetrics(values)
		}
	}
}

func getHandler() *db.Handler {
	config := db.GetConfig()
	config.User = "postgres"
	config.Database = "test_app"
	config.Host = " 127.0.0.1"
	config.Password = "12345678"
	config.Port = "5432"

	handler, err := db.GetConnection(config)
	if err != nil {
		panic(err)
	}
	return handler
}
