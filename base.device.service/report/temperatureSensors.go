package report

type TemperatureSensors struct {
	Sensor1 *TemperatureSensor `json:",omitempty"`
	Sensor2 *TemperatureSensor `json:",omitempty"`
	Sensor3 *TemperatureSensor `json:",omitempty"`
	Sensor4 *TemperatureSensor `json:",omitempty"`
}

type TemperatureSensor struct {
	Id               string
	TemperatureValue float32
}
