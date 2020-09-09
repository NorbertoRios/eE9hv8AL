package models

import (
	"testing"
	"time"
)

func TestUnmarshalActivity(t *testing.T) {
	jsonMessage := "{\"MessageEvent\":6001,\"MessageType\":0,\"UniqueId\":\"355922064286244\",\"sid\":34521956,\"Data\":{\"MessageEvent\":6001,\"UniqueId\":\"355922064286244\",\"EventCode\":6001,\"UTCDate\":\"2018/11/14\",\"UTCTime\":\"09:04:06\",\"Latitude\":50.97954,\"Longitude\":-113.95818,\"Altitude\":1046.0,\"Speed\":9.4325438231229786,\"Acceleration\":0.0,\"Deceleration\":0.0,\"Rpm\":1199,\"Heading\":115.0,\"Satellites\":16,\"HDOP\":1.0,\"MileageFromLastReset\":68586852.11999999,\"FuelConsumption\":0.0,\"Supply\":14100,\"SignalStrength\":17.0,\"GPSStatus\":1,\"GPSLostLockTime\":0,\"FuelLevel\":100.0,\"AccXforce\":\"-39\",\"AccYforce\":\"-863\",\"AccZforce\":\"-445\",\"Ack\":\"057##\",\"MessageType\":0,\"DevId\":\"xirgo_355922064286244\",\"ReceivedTime\":\"2018-11-14T09:04:10.201878Z\",\"IgnitionState\":1,\"TimeStamp\":\"2018-11-14T09:04:06Z\",\"Odometer\":68586852,\"GpsValidity\":1,\"Reason\":6,\"Type\":\"Location\",\"DTCCode\":[\"P0300\",\"P0303\"],\"PowerState\":3,\"PrevSourceId\":34521955},\"ts\":null,\"time\":\"2018-11-14T09:04:06Z\",\"Events\":[]}"
	activity := &DeviceActivity{
		Identity:           "xirgo_355922064286244",
		LastMessage:        string(jsonMessage),
		LastMessageID:      34521956,
		MessageTime:        time.Now().UTC(),
		LastUpdateTime:     time.Now().UTC(),
		DTC:                NewDTCCodes(),
		Serializedsoftware: "{\"FirmwareVersion\":\"AAz1-1168AA5\",\"OBDFirmwareVersion\":\"0b01\"}"}
	if err := activity.unmarshal(); err != nil {
		t.Error("[TestUnmarshalActivity]Unmarshal activity error:", err.Error())
	}
}
