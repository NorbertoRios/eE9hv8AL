package config

import (
	"log"

	"github.com/tkanos/gonfig"
	"queclink-go/base.device.service/utils"
)

//IConfiguration interface for service configuration
type IConfiguration interface {
	GetBase() *Configuration
}

//Configuration of service
type Configuration struct {
	Service                          string
	MysqDeviceMasterConnectionString string
	DisconnectUnmanagedDevceTimeout  int
	UpgradeFirmwareCommand           string
	SyncConfigurationInterval        float64
	TCPPort                          int
	UDPHost1                         string
	UDPHost2                         string
	UDPPort                          int
	WebAPIPort                       int
	WorkersCount                     int
	ParsingWorkersCount              int
	Rabbit                           *RabbitConfiguration
	RefreshDtcIntervalH              int
	DeviceFacadeHost                 string
	LoggerServerMode                 string
	Storage                          string
}

//GetBase return base configuraation
func (c *Configuration) GetBase() *Configuration {
	return c
}

//Config represent service configuration
var Config IConfiguration

//Initialize configuration
func Initialize(dir, fileDest string) error {
	if Config == nil {
		Config = &Configuration{}
	}
	absFileName := utils.GetAbsFilePath(dir, fileDest)
	log.Println("Loading service configuration from:", absFileName)
	err := gonfig.GetConf(absFileName, Config)
	if err != nil {
		panic("Unable to load configuration from:" + absFileName)
	}
	return err
}

//RabbitConfiguration settings for rabbit connection
type RabbitConfiguration struct {
	Host                     string
	Port                     int
	Username                 string
	Password                 string
	SystemExchange           string
	SystemRoutingKey         string
	FacadeCallbackExchange   string
	FacadeCallbackRoutingKey string
	OrleansDebugRoutingKey   string
}
