package queclinkreport

//FuelSensor struct for fuel value
type FuelSensor struct {
	SensorType int32
	Percentage int32
	Volume     int32
}

//AcSensor 1wire sensor data
type AcSensor struct {
	OWireDeviceID     string
	OWireDeviceType   int32
	DeviceDataLength  int32
	OneWireDeviceData float32
}
