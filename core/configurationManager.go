package core

import (
	"container/list"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"queclink-go/queclinkreport"

	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/report"
)

//QueclinkConfigurationManager struct to manage device configuration
type QueclinkConfigurationManager struct {
	core.ConfigurationManager
	DeviceType int
}

//Delete all old configuration/ Leave untouched othed configuration types and add configuration at the end of list
func (manager *QueclinkConfigurationManager) devideConfiguration() {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	cfgID := manager.ConfigurationModel.ID
	if cfgID == 0 || len(manager.ConfigurationModel.Command) == 0 {
		return
	}
	items := list.New()
	for c := manager.ConfigurationItems.Front(); c != nil; c = c.Next() {
		if c.Value.(core.IConfigurationItem).GetType() != "config" {
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

		messageType := strings.Split(cfgs[i], "=")
		if len(messageType) == 0 {
			continue
		}

		deviceConfigs := report.ReportConfiguration.(*queclinkreport.ReportConfiguration)
		deviceConfig, err := deviceConfigs.FindDeviceType(manager.DeviceType)
		if err != nil {
			log.Fatalf("[Device]processHbdMessage. Not found configuration for device type %v", manager.DeviceType)
			return
		}

		ack, err := deviceConfig.FindAckMessagesType(messageType[0])
		if err != nil || ack == nil {
			log.Fatalf("[ConfigurationManager]devideConfiguration. Not found ack configuration for device type %v", manager.Device.GetIdentity())
			continue

		}

		item := &core.ConfigurationItem{}
		item.ID = cfgID
		item.Type = "config"
		item.Command = cfgs[i]
		item.MessageType = fmt.Sprintf("%v", ack.ID)
		item.SentAt = min
		items.PushBack(item)
	}
	manager.ConfigurationItems = items
}

//Load get configuration from database
func (manager *QueclinkConfigurationManager) Load() {
	manager.ConfigurationModel, _ = models.FindDeviceConfigByIdentity(manager.Device.GetIdentity())
	manager.devideConfiguration()
	manager.LastLoadConfig = time.Now().UTC()
}

//UpdateMessageSent update current configuration state
func (manager *QueclinkConfigurationManager) UpdateMessageSent(message report.IMessage) {
	manager.ConfigurationManager.UpdateMessageSent(message)
	manager.SendConfiguration()
}

//InitializeConfigurationManager initialize instance of configuration manager
func InitializeConfigurationManager(deviceType int, device core.IDevice) *QueclinkConfigurationManager {
	manager := &QueclinkConfigurationManager{}
	manager.Device = device
	manager.Mutex = &sync.Mutex{}
	manager.ConfigurationItems = list.New()
	manager.DeviceType = deviceType
	return manager
}
