package model

// UsersRepository rep
type UsersRepository interface {
	GetByID(ID int)
}

// DevicesRepository rep
type DevicesRepository interface {
	GetByID(ID int)
}

// DeviceMetricsRepository rep
type DeviceMetricsRepository interface {
	GetByID(ID int)
	GetByDeviceID(ID int)
	Create(device Device)
	Update(ID int, device Device)
}

// DeviceAlertsRepository rep
type DeviceAlertsRepository interface {
	Create()
}

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
	id       int
	DeviceID int `json:"device_id"`
	Metric1  int `json:"metric1"`
	Metric2  int `json:"metric2"`
	Metric3  int `json:"metric3"`
	Metric4  int `json:"metric4"`
	Metric5  int `json:"metric5"`
}

// DeviceAlerts struct represents a device_alerts table in database
type DeviceAlerts struct {
	id       int
	deviceID int
	message  string
}

func (deviceMetrics DeviceMetrics) GetByID(ID int) {

}

func (deviceMetrics DeviceMetrics) GetByDeviceID(ID int) {

}

func (deviceMetrics DeviceMetrics) Create(ID int) {

}

func (deviceMetrics DeviceMetrics) Update(ID int, device Device) {

}
