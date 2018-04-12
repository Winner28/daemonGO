package daemon

import (
	"fmt"
	"log"

	"github.com/sevlyar/go-daemon"
)

// StartDaemon starts Daemon process
func StartDaemon() {
	fmt.Println("Daemon started")
	fmt.Println("You can kill a daemon with a help of a pid file, that contains pid of process")
	fmt.Println("(Pid file is gonna be created after the app has been started)")
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[daemon-process]"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	go registerIncomingMetricsFromDevices(1, 1)
	go monitorIncomingMetricsOfDevices(1, 1)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	for _, conn := range connections {
		conn.CloseConnection()
	}
	log.Println("Daemon process is killed")
}
