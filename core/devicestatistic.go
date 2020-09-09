package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport"
	"queclink-go/queclinkreport/gv55/evt55"
)

//DeviceStatistic gathers device behavior
type DeviceStatistic struct {
	HbdCount                 uint
	HbdDateTime              time.Time
	PdpCount                 uint
	PdpDateTime              time.Time
	LocationMessagesCount    uint
	LocationMessagesDateTime time.Time
	StartCount               uint
	StartDateTime            time.Time
	LoggedStart              bool
	CreatedAt                time.Time
	RequestCount             uint
	StartedAt                time.Time
}

//UpdateCounter of device bihavior
func (statistic *DeviceStatistic) UpdateCounter(device core.IDevice, message *queclinkreport.QueclinkMessage) {
	if statistic.CreatedAt.Add(2 * time.Hour).Before(time.Now().UTC()) {
		statistic.clear()
	}
	if statistic.StartedAt.AddDate(0, 0, 1).Before(time.Now().UTC()) {
		statistic.RequestCount = 0
		statistic.LoggedStart = false
	}

	switch message.MessageType() {
	case "+HBD", "+BBD":
		{
			statistic.HbdCount++
			statistic.HbdDateTime = time.Now().UTC()
			if statistic.HbdCount > 4 && statistic.LocationMessagesCount == 0 &&
				statistic.CreatedAt.Add(40*time.Minute).Before(time.Now().UTC()) {
				if statistic.RequestCount >= 10 {
					statistic.clear()
					return
				}
				statistic.RequestCount++
				statistic.SendClearBuffer(device)
				statistic.SendRestart(device)
			}
			break
		}
	case "+EVT", "+BVT":
		switch message.EventCode() {
		case evt55.GTPDP:
			{
				statistic.PdpCount++
				statistic.PdpDateTime = time.Now().UTC()
				break
			}
		case evt55.GTPNA:
			{
				statistic.StartDateTime = time.Now().UTC()
				statistic.StartCount++
				if statistic.StartCount >= 12 && !statistic.LoggedStart {
					statistic.LoggedStart = true
				}
				break
			}
		default:
			{
				statistic.LocationMessagesCount++
				statistic.LocationMessagesDateTime = time.Now().UTC()
				break
			}
		}
		break
	case "+RSP", "+BSP":
		{
			statistic.LocationMessagesCount++
			statistic.LocationMessagesDateTime = time.Now().UTC()
			break
		}
	}
}

func (statistic *DeviceStatistic) clear() {
	statistic.HbdCount = 0
	statistic.PdpCount = 0
	statistic.LocationMessagesCount = 0
	statistic.CreatedAt = time.Now().UTC()
	statistic.LocationMessagesDateTime = utils.MinTimeStamp()
	statistic.HbdDateTime = utils.MinTimeStamp()
	statistic.PdpDateTime = utils.MinTimeStamp()
	statistic.StartCount = 0
	statistic.StartDateTime = utils.MinTimeStamp()
}

//SendClearBuffer request to device via sms
func (statistic *DeviceStatistic) SendClearBuffer(device core.IDevice) {
	log.Println("[DeviceStatistic] Try send SendClearBuffer request.  device:", device.GetIdentity(), " to:", config.Config.GetBase().DeviceFacadeHost)

	resp, err := http.PostForm(fmt.Sprintf("%v/device/command", config.Config.GetBase().DeviceFacadeHost), url.Values{
		"identity":        {device.GetIdentity()},
		"delivery_method": {"sms_only"},
		"command":         {"AT+GTRTO=gv55,D,,,,,,FFFF$"},
	})
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("[DeviceStatistic] SendClearBuffer Invalid response for:", device.GetIdentity(), " to:", config.Config.GetBase().DeviceFacadeHost)
	}
}

//SendRestart request to device via sms
func (statistic *DeviceStatistic) SendRestart(device core.IDevice) {
	log.Println("[DeviceStatistic] Try send SendRestart request.  device:", device.GetIdentity(), " to:", config.Config.GetBase().DeviceFacadeHost)

	resp, err := http.PostForm(fmt.Sprintf("%v/device/command", config.Config.GetBase().DeviceFacadeHost), url.Values{
		"identity":        {device.GetIdentity()},
		"delivery_method": {"sms_only"},
		"command":         {"AT+GTRTO=gv55,3,,,,,,FFFF$"},
	})
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("[DeviceStatistic] SendRestart Invalid response for:", device.GetIdentity(), " to:", config.Config.GetBase().DeviceFacadeHost)
	}
}

//NewDeviceStatistic returns device statistic instance
func NewDeviceStatistic() *DeviceStatistic {
	return &DeviceStatistic{
		CreatedAt:                time.Now().UTC(),
		StartedAt:                time.Now().UTC(),
		LoggedStart:              false,
		LocationMessagesDateTime: utils.MinTimeStamp(),
		HbdDateTime:              utils.MinTimeStamp(),
		PdpDateTime:              utils.MinTimeStamp(),
		StartDateTime:            utils.MinTimeStamp(),
	}
}
