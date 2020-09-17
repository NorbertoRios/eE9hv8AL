package main

import (
	"fmt"
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

func initializeDeviceManager(service string) core.IDeviceManager {
	switch service {
	case "gv55":
		{
			dm := &qcore.QueclinkDM{}
			dm.DeviceManager.InitializeUDPDeviceCallback = dm.InitializeUDPDevice
			dm.DeviceManager.InitializeDeviceCallback = dm.InitializeDevice
			log.Println("GV55 device manager initialized")
			return dm
		}
	case "gv600":
		{
			dm := &qcore.Queclink600DM{}
			dm.DeviceManager.InitializeUDPDeviceCallback = dm.InitializeUDPDevice
			dm.DeviceManager.InitializeDeviceCallback = dm.InitializeDevice
			log.Println("GV600 device manager initialized")
			return dm
		}
	default:
		{
			panic(fmt.Sprintf("Unexpected device manager type %v", service))
		}
	}
}

func buildServiceInstance() *service.Base {
	buildCredentialsConfiguration()
	buildReportConfiguration()
	serviceInstance := &service.Base{}
	service.InitializeInstance(serviceInstance)
	manager := initializeDeviceManager(config.Config.GetBase().Service)
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
