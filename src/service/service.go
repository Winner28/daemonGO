package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"model"
	"os"
	"strconv"
	"time"
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

	insertDeviceMetrics = `INSERT INTO device_metrics (id, device_id, metric_1, metric_2, 
													  metric_3,  metric_4,  metric_5) 
						   VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	selectExists = `SELECT exists(SELECT device_id FROM device_metrics WHERE device_id=$1);`
)

// Statements for working with users, devices, and device_alerts tables
const (
	selectUser        = `SELECT id, name, email FROM users WHERE id=$1;`
	insertDeviceAlert = `INSERT INTO device_alerts (id, device_id, message) 
						  VALUES ($1, $2, $3) RETURNING id;`
	selectDevice = `SELECT id, user_id, name FROM devices WHERE id=$1`
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

	// Last id in device_alerts store
	lastDeviceAlertID int
)

func init() {
	lastDeviceAlertID = 1
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

// UpdateMetrics updates device_metrics of given DEVICE_ID and metrics values
// that contains in passed metrics map (key = metricsID, value = metricValue
// (TODO: swap with array [5]int))
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

// MonitorMetrics from database and adjusts it to the optimal metrics values for device
// [metrics] map contains optimal metrics values for specified device
// * if something goes bad and metrics from database is greater than the optimal ones we:
// ** log it into file [/src/log]
// ** formulate error message [user-device-device_metrics] that writes into device_alerts table
// ** check redis-cache, if we already have message for that device => overwrite old one msg
func (handler *Handler) MonitorMetrics(from int, to int) {
	log.Println("Monitoring metrics!")
	for ID := from; ID <= to; ID++ {
		row := handler.Db.QueryRow(selectDeviceMetrics, ID)
		var metr model.DeviceMetrics
		err := row.Scan(&metr.ID, &metr.DeviceID, &metr.Metric1,
			&metr.Metric2, &metr.Metric3, &metr.Metric4, &metr.Metric5)
		badMetrics, ok := checkMetrics(metr)
		if !ok {
			log.Println("Looks like we got metrics trouble on Device_ID:", metr.DeviceID)
			handler.alertDeviceMetricsError(metr, badMetrics)
		}
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

// Alert message to device_alerts table, message contains information about:
// concrete user, device and the device_metrics that out of range
// [] badMetrics represents a numbers of bad metrics
func (handler *Handler) alertDeviceMetricsError(deviceMetrics model.DeviceMetrics, badMetrics []int) {
	row := handler.Db.QueryRow(selectDevice, deviceMetrics.DeviceID)
	var device model.Device
	var user model.User
	err := row.Scan(&device.ID, &device.UserID, &device.Name)

	switch err {
	case sql.ErrNoRows:
		log.Println("We dont have such device")
	case nil:
		log.Println("Selecting Device with ID:", device.ID)
	default:
		panic(err)
	}
	row = handler.Db.QueryRow(selectUser, device.UserID)
	err = row.Scan(&user.ID, &user.Name, &user.Email)

	switch err {
	case sql.ErrNoRows:
		log.Println("We dont have such device")
	case nil:
		log.Println("Selecting User with Name:", user.Name)
	default:
		panic(err)
	}

	message := createErrorMessage(user, device, deviceMetrics, badMetrics)
	log.Println(message)
	/* handler.Db.QueryRow(insertDeviceAlert, lastID, device.ID, message).Scan(&lastDeviceAlertID)
	*(&lastDeviceAlertID)++ */
}

// Specified error message for device_metrics, that out of range
func createErrorMessage(user model.User, device model.Device, deviceMetrics model.DeviceMetrics, badMetrics []int) string {
	message := "TROUBLE WITH METRICS!\n" + user.ToString() + device.ToString()
	message += "\n\nFor Device Metrics ID: " + strconv.Itoa(device.ID)
	for _, value := range badMetrics {
		switch value {
		case 1:
			message += fmt.Sprintf("\n*Metric_1*  Prefered: %v, ACTUAL: %v", metrics[device.ID].Metric1, deviceMetrics.Metric1)
			break
		case 2:
			message += fmt.Sprintf("\n*Metric_2*  Prefered: %v, ACTUAL: %v", metrics[device.ID].Metric2, deviceMetrics.Metric2)
			break
		case 3:
			message += fmt.Sprintf("\n*Metric_3*  Prefered: %v, ACTUAL: %v", metrics[device.ID].Metric3, deviceMetrics.Metric2)
			break
		case 4:
			message += fmt.Sprintf("\n*Metric_4*  Prefered: %v, ACTUAL: %v", metrics[device.ID].Metric4, deviceMetrics.Metric2)
			break
		case 5:
			message += fmt.Sprintf("\n*Metric_5*  Prefered: %v, ACTUAL: %v", metrics[device.ID].Metric5, deviceMetrics.Metric2)
			break
		}
	}
	message += fmt.Sprint("\n\nReported Time:", time.Now().Format("Mon Jan _2 15:04:05 2006"), "\n\n")

	return message
}

// Checking if we already got device_metric with such device_id
func (handler *Handler) deviceMetricExists(ID int) bool {
	var exists bool
	err := handler.Db.QueryRow(selectExists, ID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking if row exists")
	}

	return exists
}

// checking if metrics are out of specified range
// (range specified in /config/config.metrics.json file)
func checkMetrics(device model.DeviceMetrics) ([]int, bool) {
	badMetrics := make([]int, 5)
	ok := true
	metric := metrics[device.DeviceID]
	if device.Metric1 > metric.Metric1 {
		ok = false
		badMetrics = append(badMetrics, 1)
	}
	if device.Metric2 > metric.Metric2 {
		ok = false
		badMetrics = append(badMetrics, 2)
	}
	if device.Metric3 > metric.Metric3 {
		ok = false
		badMetrics = append(badMetrics, 3)
	}
	if device.Metric4 > metric.Metric4 {
		ok = false
		badMetrics = append(badMetrics, 4)
	}
	if device.Metric5 > metric.Metric5 {
		ok = false
		badMetrics = append(badMetrics, 5)
	}

	return badMetrics, ok
}
