package model

import "time"

// User struct represents a user table in database
type User struct {
	id    int
	name  string
	email string
}

// Device struct represents a devices table in database
type Device struct {
	id     int
	name   string
	userID int
}

// DeviceMetrics struct represents a device_metrics table in database
type DeviceMetrics struct {
	ID         int
	DeviceID   int `json:"device_id"`
	Metric1    int `json:"metric1"`
	Metric2    int `json:"metric2"`
	Metric3    int `json:"metric3"`
	Metric4    int `json:"metric4"`
	Metric5    int `json:"metric5"`
	LocalTime  time.Time
	ServerTime time.Time
}

// DeviceAlerts struct represents a device_alerts table in database
type DeviceAlerts struct {
	id       int
	deviceID int
	message  string
}

// DeviceMetricsRange range
type DeviceMetricsRange struct {
	DeviceID int `json:"device_id"`
	Min      int `json:"min"`
	Max      int `json:"max"`
}
