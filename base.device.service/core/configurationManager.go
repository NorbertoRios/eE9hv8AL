package core

import (
	"container/list"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/rabbit"
	"queclink-go/base.device.service/report"
	"github.com/streadway/amqp"
)

//IConfigurationItem interface for configuration item
type IConfigurationItem interface {
	GetID() int32
	GetCommand() string
	GetMessageType() string
	SetMessageType(messageType string)
	GetType() string
	GetSentAt() time.Time
	SetSentAt(time.Time)
}

//ConfigurationItem configuration item
type ConfigurationItem struct {
	ID          int32
	Command     string
	MessageType string
	Type        string
	SentAt      time.Time
}

//GetID returns command database id
func (c *ConfigurationItem) GetID() int32 {
	return c.ID
}

//GetCommand returns command
func (c *ConfigurationItem) GetCommand() string {
	return c.Command
}

//GetMessageType returns command MessageType
func (c *ConfigurationItem) GetMessageType() string {
	return c.MessageType
}

//GetType returns type "config", "api_immobilizer_request", "api_request"
func (c *ConfigurationItem) GetType() string {
	return c.Type
}

//GetSentAt returns last sent timestamp
func (c *ConfigurationItem) GetSentAt() time.Time {
	return c.SentAt
}

//SetSentAt set sent command timestamp
func (c *ConfigurationItem) SetSentAt(timestamp time.Time) {
	c.SentAt = timestamp
}

//SetMessageType set type of message
func (c *ConfigurationItem) SetMessageType(messageType string) {
	c.MessageType = messageType
}

//OnCommandSent set callback for sent command
func (c *ConfigurationItemAPI) OnCommandSent(callback func(callbackId string, code string, result bool)) {
	c.onCommandSent = callback
}

//SendResponse for API call request
func (c *ConfigurationItemAPI) SendResponse(code string, result bool) bool {
	if c.onCommandSent != nil {
		c.onCommandSent(c.CallbackID, code, result)
		return true
	}
	return false
}

//ConfigurationItemAPI api request configuration item
type ConfigurationItemAPI struct {
	ConfigurationItem
	CallbackID    string
	TTL           int
	CreationTime  time.Time
	onCommandSent func(callbackId string, code string, result bool)
}

//ConfigurationItemImmobilizerAPI configuration item
type ConfigurationItemImmobilizerAPI struct {
	ConfigurationItemAPI
	SafetyOption bool
}

//IConfigurationManager interface for canfiguration manager
type IConfigurationManager interface {
	GetModel() *models.DeviceConfig
	GetItems() *list.List
	UpdateMessageSent(message report.IMessage)
	SyncDeviceConfig()
	SendConfiguration()
	GetLocateCommand() string
	Load()
	AddCommand(confItem IConfigurationItem, toFirstPosition bool) bool
	GetSelf() IConfigurationManager
	SetSelf(self IConfigurationManager)
	Initialize(IDevice)
	RemoveDuplicates(confItem IConfigurationItem)
}

//ConfigurationManager device configuration manager
type ConfigurationManager struct {
	ConfigurationModel *models.DeviceConfig
	ConfigurationItems *list.List
	Device             IDevice
	Mutex              *sync.Mutex
	LastLoadConfig     time.Time
	SendingItem        IConfigurationItem
	Self               IConfigurationManager
}

//GetSelf returns interface of self instance
func (manager *ConfigurationManager) GetSelf() IConfigurationManager {
	return manager.Self
}

//SetSelf assign self
func (manager *ConfigurationManager) SetSelf(self IConfigurationManager) {
	manager.Self = self
}

//GetLocateCommand returns device locate command
func (manager *ConfigurationManager) GetLocateCommand() string {
	panic("Not implemented exception ConfigurationManager.GetLocateCommand")
}

//GetModel return configuration model
func (manager *ConfigurationManager) GetModel() *models.DeviceConfig {
	return manager.ConfigurationModel
}

//GetItems return parsed configuration items
func (manager *ConfigurationManager) GetItems() *list.List {
	return manager.ConfigurationItems
}

//Delete all old configuration/ Leave untouched othed configuration types and add configuration at the end of list
func (manager *ConfigurationManager) devideConfiguration() {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	cfgID := manager.ConfigurationModel.ID
	if cfgID == 0 || len(manager.ConfigurationModel.Command) == 0 {
		return
	}
	items := list.New()
	for c := manager.ConfigurationItems.Front(); c != nil; c = c.Next() {
		if c.Value.(IConfigurationItem).GetType() != "config" {
			items.PushBack(c.Value)
		}
	}
	configuration := manager.ConfigurationModel.Command

	re := regexp.MustCompile(`(\n)|(\r\n)`)
	cfgs := re.Split(configuration, -1)
	min, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	for i := range cfgs {
		if len(cfgs[i]) == 0 {
			continue
		}
		if strings.Contains(cfgs[i], "+XT:1010,") {
			continue
		}

		messageType := strings.Split(cfgs[i], ",")
		if len(messageType) == 0 {
			continue
		}
		item := &ConfigurationItem{
			ID:          cfgID,
			Type:        "config",
			Command:     cfgs[i],
			MessageType: messageType[1],
			SentAt:      min,
		}
		items.PushBack(item)
	}
	manager.ConfigurationItems = items
}

//SendConfiguration sends configuration to device
func (manager *ConfigurationManager) SendConfiguration() {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	if manager.ConfigurationItems.Len() == 0 {
		return
	}
	if manager.SendingItem != nil && time.Now().UTC().Sub(manager.SendingItem.GetSentAt()).Seconds() < 10 {
		return
	}

	cmd := manager.ConfigurationItems.Front().Value.(IConfigurationItem)
	log.Println("[DeviceConfigurationManager]Sending configuration item:", cmd.GetCommand(), "DevId:", manager.Device.GetIdentity())
	if manager.ConfigurationItems.Len() == 0 {
		return
	}
	manager.SendingItem = cmd
	manager.SendingItem.SetSentAt(time.Now().UTC())
	manager.Device.SendString(manager.SendingItem.GetCommand())
}

//UpdateCommandStatus update status of command db/api/etc
func (manager *ConfigurationManager) UpdateCommandStatus(confItem IConfigurationItem) {
	if confItem == nil {
		return
	}
	if confItem.GetType() == "config" && confItem.GetID() != 0 {
		manager.updateCommandDBStatus(confItem)
		return
	}

	if confItem.GetType() == "api_request" {
		if confItem.(*ConfigurationItemAPI).SendResponse("Done", true) {
			manager.Device.SendSystemMessageS("Request "+confItem.GetCommand()+" is done", config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		}
	}

	if confItem.GetType() == "api_immobilizer_request" {
		if confItem.(*ConfigurationItemImmobilizerAPI).SendResponse("Done", true) {
			manager.Device.SendSystemMessageS("Request "+confItem.GetCommand()+" is done", config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		}
	}
}

func (manager *ConfigurationManager) updateCommandDBStatus(confItem IConfigurationItem) {

	if confItem == nil {
		return
	}

	for c := manager.ConfigurationItems.Front(); c != nil; c = c.Next() {
		iConfig := c.Value.(IConfigurationItem)
		if iConfig.GetID() == confItem.GetID() {
			return
		}

	}
	manager.ConfigurationModel.UpdateSentConfiguration()

	manager.Device.SendSystemMessageS("Configuration(id:"+fmt.Sprintf("%v", confItem.GetID())+
		") uploaded to device:"+manager.Device.GetIdentity(),
		config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
}

//Load get configuration from database
func (manager *ConfigurationManager) Load() {
	manager.ConfigurationModel, _ = models.FindDeviceConfigByIdentity(manager.Device.GetIdentity())
	manager.devideConfiguration()
	manager.LastLoadConfig = time.Now().UTC()
}

//Count of configuration waiting to be send
func (manager *ConfigurationManager) Count() int {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	return manager.ConfigurationItems.Len()
}

//Initialize configuration manager
func (manager *ConfigurationManager) Initialize(device IDevice) {
	manager.Device = device
	manager.ConfigurationItems = list.New()
	manager.Self.Load()
}

//UpdateMessageSent update current configuration state
func (manager *ConfigurationManager) UpdateMessageSent(message report.IMessage) {

	manager.Mutex.Lock()

	defer func() {
		if r := recover(); r != nil {
			manager.Mutex.Unlock()
			log.Println("Recovered in UpdateMessageSent:", r)
		}
		manager.SendConfiguration()
	}()

	var confItem IConfigurationItem
	var el *list.Element

	for c := manager.ConfigurationItems.Front(); c != nil; c = c.Next() {
		iConfig := c.Value.(IConfigurationItem)
		cMessageType := strings.ToUpper(iConfig.GetMessageType())
		mMessageType := strconv.FormatInt(int64(message.EventCode()), 10)
		if cMessageType == mMessageType {
			confItem = c.Value.(IConfigurationItem)
			el = c
			break
		}
	}

	if confItem == nil {
		log.Println("[DeviceConfigurationManager]Configuration not found for reply:", message.GetStringValue("Message", ""))
	}
	if el != nil {
		manager.ConfigurationItems.Remove(el)
	}

	manager.Mutex.Unlock()

	if confItem != nil {
		manager.UpdateCommandStatus(confItem)
		if manager.SendingItem != nil &&
			manager.SendingItem.GetCommand() == confItem.GetCommand() {
			manager.SendingItem = nil
		}
	}
}

//SyncDeviceConfig synchronize device configuration
func (manager *ConfigurationManager) SyncDeviceConfig() {
	if time.Now().UTC().Sub(manager.LastLoadConfig).Seconds() >= config.Config.GetBase().SyncConfigurationInterval {
		manager.SendingItem = nil
		manager.Self.Initialize(manager.Device)
	}
	manager.Self.SendConfiguration()
}

//RemoveDuplicates commands by command string
func (manager *ConfigurationManager) RemoveDuplicates(confItem IConfigurationItem) {
	var el *list.Element
	for c := manager.ConfigurationItems.Front(); c != nil; c = c.Next() {
		iConfig := c.Value.(IConfigurationItem)
		if strings.ToUpper(iConfig.GetCommand()) == confItem.GetCommand() {
			el = c
			break
		}
	}
	if el != nil {
		manager.ConfigurationItems.Remove(el)
		manager.Self.RemoveDuplicates(confItem)
	}
}

//AddCommand to configuration manager
func (manager *ConfigurationManager) AddCommand(confItem IConfigurationItem, toFirstPosition bool) bool {

	if confItem.GetMessageType() == "" {
		messageType := strings.Split(confItem.GetCommand(), ",")
		if len(messageType) == 0 {
			return false
		}
		confItem.SetMessageType(messageType[0])
	}

	manager.Self.RemoveDuplicates(confItem)
	if toFirstPosition {
		manager.ConfigurationItems.PushFront(confItem)
	} else {
		manager.ConfigurationItems.PushBack(confItem)
	}
	return true
}

//InitializeConfigurationManager initialize instance of configuration manager
func InitializeConfigurationManager(device *Device) *ConfigurationManager {
	return &ConfigurationManager{
		Device:             device,
		Mutex:              &sync.Mutex{},
		ConfigurationItems: list.New(),
	}
}

//SendConfigurationAPIResponse for API call request
func SendConfigurationAPIResponse(callbackID string, code string, result bool) {

	response := &APIResponse{
		CallbackID: callbackID,
		Success:    result,
		Code:       code,
	}
	sResponse := response.Marshal()
	log.Println("[ConfigurationManager]", "SendCallbackToFacade. Sends callback ", sResponse)
	rabbit.RabbitConnection.Publish(sResponse, config.Config.GetBase().Rabbit.FacadeCallbackExchange, config.Config.GetBase().Rabbit.FacadeCallbackRoutingKey, make(amqp.Table, 0), 10)
}
