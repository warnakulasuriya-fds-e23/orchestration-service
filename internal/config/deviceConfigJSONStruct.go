package config

type DevicesConfigJSON struct {
	EnrollmentDevices []string                 `json:"enrollmentDevices"`
	DeviceDetails     map[string]DeviceDetails `josn:"deviceDetails"`
}

type DeviceDetails struct {
	Floor string `json:"floor"`
	Door  string `json:"door"`
}
