package main

import (
	"log"
	"os"

	qcore "queclink-go/core"
	"queclink-go/qapi/qcontroller"
	"queclink-go/queclinkreport"

	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/service"
)

func buildReportConfiguration() {
	if err := report.LoadReportConfiguration(os.Args[0], "/config/initializer/ReportConfiguration.xml", &queclinkreport.ReportConfiguration{}); err != nil {
		panic("Error while report config initializing: %v" + err.Error())
	}
}

func buildCredentialsConfiguration() {
	if err := config.Initialize(os.Args[0], "/config/initializer/credentials.example.json"); err != nil {
		panic("Error while credentials config initializing: " + err.Error())
	}
}

func initializeDeviceManagerCallbacks(manager core.IDeviceManager) core.IDeviceManager {
	switch manager.(type) {
	case *qcore.QueclinkDM:
		{
			dm, _ := manager.(*qcore.QueclinkDM)
			dm.DeviceManager.InitializeUDPDeviceCallback = dm.InitializeUDPDevice
			dm.DeviceManager.InitializeDeviceCallback = dm.InitializeDevice
			return dm
		}
	case *qcore.Queclink600DM:
		{
			dm, _ := manager.(*qcore.Queclink600DM)
			dm.DeviceManager.InitializeUDPDeviceCallback = dm.InitializeUDPDevice
			dm.DeviceManager.InitializeDeviceCallback = dm.InitializeDevice
			return dm
		}
	default:
		{
			panic("Unexpected device manager type")
		}
	}
}

func buildServiceInstance() *service.Base {
	buildCredentialsConfiguration()
	buildReportConfiguration()
	serviceInstance := &service.Base{}
	service.InitializeInstance(serviceInstance)
	var dm core.IDeviceManager
	serviceName := config.Config.GetBase().Service
	switch serviceName {
	case "gv55":
		{
			dm = &qcore.QueclinkDM{}
		}
	case "gv600":
		{
			dm = &qcore.Queclink600DM{}
		}
	default:
		{
			panic("Unexpected service type: " + serviceName)
		}
	}
	if dm == nil {
		panic("Device manager is nil")
	}
	manager := initializeDeviceManagerCallbacks(dm)
	serviceInstance.ConfigureDeviceManager(manager).
		ConfigureParser(&queclinkreport.Parser{}).
		ConfigureApi(&qcontroller.QController{}).
		Initialize()
	log.Println("Service instance successfully build")
	return serviceInstance
}

func main() {
	service.InitLogWrapper()
	serviceInstance := buildServiceInstance()
	serviceInstance.Start()
}
