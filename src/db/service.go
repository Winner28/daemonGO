package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"model"
	"os"
)

const (
	updateStatement = `UPDATE device_metrics SET 
	metric_1 = $1, metric_2 = $2,
	metric_3 = $3, metric_4 = $4, metric_5 = $5 
	WHERE device_id=$6`

	selectStatement = `SELECT id, device_id, metric_1,  metric_2,  metric_3,  metric_4,  metric_5 FROM device_metrics WHERE id=$1;`
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
	res, err := handler.Db.Exec(updateStatement, metrics[1], metrics[2],
		metrics[3], metrics[4], metrics[5], ID)
	if err != nil {
		log.Panic(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}
}

// MonitorMetrics m
func (handler *Handler) MonitorMetrics(from int, to int) {
	log.Println("Monitoring metrics!")
	for ID := from; ID <= to; ID++ {
		row := handler.Db.QueryRow(selectStatement, ID)
		var metr model.DeviceMetrics
		err := row.Scan(&metr.ID, &metr.DeviceID, &metr.Metric1,
			&metr.Metric2, &metr.Metric3, &metr.Metric4, &metr.Metric5)
		fmt.Println(metr)
		switch err {
		case sql.ErrNoRows:
			fmt.Println("We dont have metrics for such device ID")
			return
		case nil:
			log.Println("We dont have metrics for such device")
		default:
			panic(err)
		}
	}
}
