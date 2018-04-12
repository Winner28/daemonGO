package model

import (
	"strconv"
	"time"
)

// User struct represents a `users` table in the database
type User struct {
	ID    int
	Name  string
	Email string
}

// Device struct represents a `devices` table in the database
type Device struct {
	ID     int
	Name   string
	UserID int
}

// DeviceMetrics struct represents a `device_metrics` table in the database
type DeviceMetrics struct {
	ID        int
	DeviceID  int `json:"device_id"`
	Metric1   int `json:"metric1"`
	Metric2   int `json:"metric2"`
	Metric3   int `json:"metric3"`
	Metric4   int `json:"metric4"`
	Metric5   int `json:"metric5"`
	LocalTime time.Time
}

// DeviceAlerts struct represents a `device_alerts` table in the database
type DeviceAlerts struct {
	id       int
	deviceID int
	message  string
}

// DeviceMetricsRange to this struct json parses from [/config/config.metrics_range.json] file
type DeviceMetricsRange struct {
	DeviceID int `json:"device_id"`
	Min      int `json:"min"`
	Max      int `json:"max"`
}

// ToString method returns string-representation of User
func (user User) ToString() string {
	toReturn := "User information:\n"
	toReturn += "User ID: " + strconv.Itoa(user.ID)
	toReturn += "\nUser Name: " + user.Name
	toReturn += "\nUser Email: " + user.Email

	return toReturn
}

// ToString method returns string-representation of Device
func (device Device) ToString() string {
	toReturn := "\nDevice information: \n"
	toReturn += "\nDevice ID: " + strconv.Itoa(device.ID)
	toReturn += "\nDevice Name: " + device.Name

	return toReturn
}
