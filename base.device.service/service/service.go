package service

import (
	"log"
	"os"
	"time"

	"queclink-go/base.device.service/utils"

	"github.com/mitchellh/panicwrap"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"queclink-go/base.device.service/api"
	"queclink-go/base.device.service/api/controller"
	_ "queclink-go/base.device.service/api/docs"
	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/rabbit"
	"queclink-go/base.device.service/report"
)

//Instance of device manager
var Instance IService

//IService interface for service
type IService interface {
	Initialize()
	ConfigureReport(filename string, config report.IReportConfiguration) *Base
	ConfigureConfig(filename string, config config.IConfiguration) *Base
	ConfigureDeviceManager(dm core.IDeviceManager) *Base
	ConfigureParser(parser report.IParser) *Base
	ConfigureApi(controller controller.IApiController) *Base
	GetAPIServer() *api.Server
	Start()
}

//Base collect service features to start
type Base struct {
	UDPServer1                  *comm.UDPServer
	UDPServer2                  *comm.UDPServer
	ReportConfigurationFilename string
	ReportConfiguration         report.IReportConfiguration
	ConfigurationFilename       string
	Configuration               config.IConfiguration
	DeviceManager               core.IDeviceManager
	TCPServer                   *comm.TCPServer
	APIServer                   *api.Server
	Parser                      report.IParser
	APIController               controller.IApiController
}

//InitializeInstance initialize service instance
func InitializeInstance(instance IService) {
	Instance = instance
}

//Initialize service
func (service *Base) Initialize() {
	models.InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString,
		config.Config.GetBase().LoggerServerMode)
	rabbit.InitializeRabbitConnection(&rabbit.Credentials{
		Host:     config.Config.GetBase().Rabbit.Host,
		Port:     config.Config.GetBase().Rabbit.Port,
		Username: config.Config.GetBase().Rabbit.Username,
		Password: config.Config.GetBase().Rabbit.Password,
	})
	service.DeviceManager.SetParser(service.Parser)
	core.InitializeDeviceManager(service.DeviceManager)
	service.APIServer = api.StartNewAPIServer(service.APIController,
		config.Config.GetBase().WebAPIPort,
		config.Config.GetBase().LoggerServerMode)
	log.Println("Started web API server on:", service.APIServer.Port)
	service.UDPServer1 = comm.NewUDPServer(config.Config.GetBase().UDPHost1, config.Config.GetBase().UDPPort)
	service.UDPServer1.OnPacket(service.DeviceManager.OnUDPPacket)
	//service.UDPServer2 = comm.NewUDPServer(config.Config.GetBase().UDPHost2, config.Config.GetBase().UDPPort)
	//service.UDPServer2.OnPacket(service.DeviceManager.OnUDPPacket)
	service.TCPServer = comm.NewTCPServer("", config.Config.GetBase().TCPPort)
	service.TCPServer.OnNewClient(service.DeviceManager.ReceivedNewConnection)
}

//ConfigureReport sets filename and configuration instance
func (service *Base) ConfigureReport(filename string, config report.IReportConfiguration) *Base {
	service.ReportConfigurationFilename = filename
	service.ReportConfiguration = config
	return service
}

//ConfigureConfig sets filename and configuration instance
func (service *Base) ConfigureConfig(filename string, config config.IConfiguration) *Base {
	service.ConfigurationFilename = filename
	service.Configuration = config
	return service
}

//ConfigureDeviceManager configure device menager for service
func (service *Base) ConfigureDeviceManager(dm core.IDeviceManager) *Base {
	service.DeviceManager = dm
	return service
}

//ConfigureParser configure instanse for parser
func (service *Base) ConfigureParser(parser report.IParser) *Base {
	service.Parser = parser
	return service
}

//ConfigureApi ...
func (service *Base) ConfigureApi(controller controller.IApiController) *Base {
	service.APIController = controller
	return service
}

//GetAPIServer returns instance
func (service *Base) GetAPIServer() *api.Server {
	return service.APIServer
}

//Start service
func (service *Base) Start() {
	go service.UDPServer1.Listen()
	//go service.UDPServer2.Listen()
	service.TCPServer.Listen()
	service.APIServer = nil
}

//InitLogWrapper init log wrapper
func InitLogWrapper() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(&lumberjack.Logger{
		Filename:   utils.GetAbsFilePath(os.Args[0], "/logs/console.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 10,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
	})
}

func initPanicWrapper() {
	panicwrap.BasicWrap(panicHandler)
}

func panicHandler(output string) {
	log.Fatalf("[%v] error >> %v", time.Now(), output)
}
