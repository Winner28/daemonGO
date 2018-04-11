package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"model"
	"os"
)

// Statements for working with device_metrics table
const (
	updateDeviceMetrics = `UPDATE device_metrics SET 
					   metric_1 = $1, metric_2 = $2,
					   metric_3 = $3, metric_4 = $4, metric_5 = $5 
					   WHERE device_id=$6`

	selectDeviceMetrics = `SELECT id, device_id, metric_1,  metric_2,  metric_3,  metric_4,  metric_5
					   FROM device_metrics 
					   WHERE id=$1;`

	selectExists = `SELECT exists(SELECT device_id FROM device_metrics WHERE device_id=$1)`

	insertDeviceMetrics = `INSERT INTO device_metrics (id, device_id, metric_1, metric_2, 
													  metric_3,  metric_4,  metric_5) 
						   VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
)

var (
	// Contains the optimal metrics for a given devices
	// key = deviceID, value = device
	// optimal metrics are set in the config "./config/config.metrics.json"
	metrics map[int]model.DeviceMetrics

	// Contains the information about if the metrics with such device_id contains in database
	// (is used to not turn each time in db)
	// key = deviceID, value = true/false
	contains map[int]bool

	// Last metric_id in device_metrics store
	lastID int
)

func init() {
	contains = make(map[int]bool)
	metrics = make(map[int]model.DeviceMetrics, 1000)
	for _, value := range parseMetricsJSON() {
		metrics[value.DeviceID] = value
	}
}

// Parses optimal metrics for devices from a file which lies under "./config/config.metrics.json"
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

// UpdateMetrics updates device_metrics of given ID and metrics values
func (handler *Handler) UpdateMetrics(ID int, metrics map[int]int) {
	if !contains[ID] {
		if !handler.deviceMetricExists(ID) {
			lastID = ID
			var id int
			err := handler.Db.QueryRow(insertDeviceMetrics, lastID, ID, metrics[1], metrics[2],
				metrics[3], metrics[4], metrics[5]).Scan(&id)
			if err != nil {
				panic(err)
			}
			return
		}
		contains[ID] = true
	}
	res, err := handler.Db.Exec(updateDeviceMetrics, metrics[1], metrics[2],
		metrics[3], metrics[4], metrics[5], ID)
	if err != nil {
		log.Panic(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}
}

// MonitorMetrics from database
func (handler *Handler) MonitorMetrics(from int, to int) {
	log.Println("Monitoring metrics!")
	for ID := from; ID <= to; ID++ {
		row := handler.Db.QueryRow(selectDeviceMetrics, ID)
		var metr model.DeviceMetrics
		err := row.Scan(&metr.ID, &metr.DeviceID, &metr.Metric1,
			&metr.Metric2, &metr.Metric3, &metr.Metric4, &metr.Metric5)
		switch err {
		case sql.ErrNoRows:
			log.Println("We dont have metrics for such device ID")
		case nil:
			fmt.Println(metr)
		default:
			panic(err)

		}
	}
}

// Uses to check if we already got device_metric with such device_id
func (handler *Handler) deviceMetricExists(ID int) bool {
	var exists bool
	err := handler.Db.QueryRow(selectExists, ID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking if row exists")
	}

	return exists
}
