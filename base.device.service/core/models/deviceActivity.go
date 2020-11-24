package models

import (
	"errors"
	"time"

	"queclink-go/base.device.service/report"

	"gorm.io/gorm"
)

//DeviceActivity model
type DeviceActivity struct {
	Identity           string    `gorm:"column:daiDeviceIdentity"`
	MessageTime        time.Time `gorm:"column:daiLastMessageTime"`
	LastUpdateTime     time.Time `gorm:"column:daiLastUpdateTime"`
	LastMessageID      uint64    `gorm:"column:daiLastMessageId"`
	LastMessage        string    `gorm:"column:daiLastMessage"`
	Serializedsoftware string    `gorm:"column:daiSoftware"`

	Software     *Software `gorm:"-" sql:"-"`
	BatteryLevel float32   `gorm:"-" sql:"-"`
	SatFix       int32     `gorm:"-" sql:"-"`
	Ignition     string    `gorm:"-" sql:"-"`
	TimeStamp    time.Time `gorm:"-" sql:"-"`
	GPSTimeStamp time.Time `gorm:"-" sql:"-"`
	Latitude     float32   `gorm:"-" sql:"-"`
	Longitude    float32   `gorm:"-" sql:"-"`
	Odometer     int32     `gorm:"-" sql:"-"`
	PowerState   string    `gorm:"-" sql:"-"`
	Relay        byte      `gorm:"-" sql:"-"`
	DTC          *DTCCodes `gorm:"-" sql:"-"`
	FuelLevel    float32   `gorm:"-" sql:"-"`
}

//TableName for DeviceActivity model
func (DeviceActivity) TableName() string {
	return "ats.tblDeviceActivityInfo"
}

//Save device activity to database
func (activity *DeviceActivity) Save() error {
	activity.LastUpdateTime = time.Now().UTC()
	activity.Serializedsoftware = activity.Software.Marshal()
	db := rawdb.Model(&activity).Where("daiDeviceIdentity=?", activity.Identity).Updates(&activity)
	if db.RowsAffected == 0 {
		db = rawdb.Create(&activity)
	}
	return db.Error
}

//BeforeUpdate prepare DeviceActivity for save
func (activity *DeviceActivity) BeforeUpdate() (err error) {
	activity.Serializedsoftware = activity.Software.Marshal()
	return nil
}

//BeforeCreate prepare DeviceActivity for save
func (activity *DeviceActivity) BeforeCreate() (err error) {
	activity.BeforeUpdate()
	return nil
}

//BeforeSave prepare DeviceActivity for save
func (activity *DeviceActivity) BeforeSave() (err error) {
	activity.BeforeUpdate()
	return nil
}

//AfterUpdate unmarshal string to struct
func (activity *DeviceActivity) AfterUpdate() (err error) {
	activity.AfterFind()
	return nil
}

//AfterSave unmarshal string to struct
func (activity *DeviceActivity) AfterSave() (err error) {
	activity.AfterFind()
	return nil
}

//AfterCreate unmarshal string to struct
func (activity *DeviceActivity) AfterCreate() (err error) {
	activity.AfterFind()
	return nil
}

//AfterFind unmarshal string to struct
func (activity *DeviceActivity) AfterFind() (err error) {
	activity.Software, _ = UnMarshalSoftware(activity.Serializedsoftware)
	return nil
}

func (activity *DeviceActivity) unmarshal() error {
	if activity.Software == nil {
		activity.Software = &Software{}
	}

	strMsg := activity.LastMessage
	if strMsg == "" {
		return errors.New("Last Message is empty for device:" + activity.Identity)
	}
	message, err := report.UnMarshalMessage(strMsg)
	if err != nil {
		return err
	}
	if v, f := message.GetValue("BatteryPercentage"); f {
		activity.BatteryLevel = float32(v.(float64))
	}

	if v, f := message.GetValue("Satellites"); f {
		activity.SatFix = int32(v.(float64))
	}

	if v, f := message.GetValue("IgnitionState"); f {
		if v.(float64) == 1 {
			activity.Ignition = "On"
		} else {
			activity.Ignition = "Off"
		}
	}

	if v, f := message.GetValue("TimeStamp"); f {
		if v1, e := time.Parse("2006-01-02T15:04:05Z", v.(string)); e == nil {
			activity.GPSTimeStamp = v1
		}
	}

	if v, f := message.GetValue("ReceivedTime"); f {
		if v1, e := time.Parse("2006-01-02T15:04:05Z", v.(string)); e == nil {
			activity.TimeStamp = v1
		}
	}

	if v, f := message.GetValue("Latitude"); f {
		activity.Latitude = float32(v.(float64))
	}

	if v, f := message.GetValue("Longitude"); f {
		activity.Longitude = float32(v.(float64))
	}

	if v, f := message.GetValue("Odometer"); f {
		activity.Odometer = int32(v.(float64))
	}

	if v, f := message.GetValue("PowerState"); f {
		if ps, valid := v.(string); valid {
			activity.PowerState = ps
		} else {
			if ps, valid := v.(float64); valid {
				switch ps {
				case 1:
					activity.PowerState = "Unknown"
					break
				case 2:
					activity.PowerState = "Sleep Mode"
					break
				case 3:
					activity.PowerState = "Powered"
					break
				case 4:
					activity.PowerState = "Power Lost"
					break
				case 5:
					activity.PowerState = "Backup battery"
					break
				case 6:
					activity.PowerState = "Power off"
					break
				}
			}
		}
	}
	if v, f := message.GetValue("Relay"); f {
		activity.Relay = byte(v.(float64))
	}

	activity.DTC = NewDTCCodes()
	if v, f := message.GetValue("DTCCode"); f {
		dtcs := []string{}
		for _, dtc := range v.([]interface{}) {
			dtcs = append(dtcs, dtc.(string))
		}
		activity.DTC.LoadCodes(dtcs)
	}

	if v, f := message.GetValue("FuelLevel"); f {
		activity.FuelLevel = float32(v.(float64))
	}

	return nil
}

//FindDeviceActivityInfo lookup not sent device configuration
func FindDeviceActivityInfo(identity string) (*DeviceActivity, bool) {
	d := &DeviceActivity{}
	err := rawdb.Where("daiDeviceIdentity=?", identity).Find(d).Error
	bErr := errors.Is(err, gorm.ErrRecordNotFound)
	d.unmarshal()
	if d.DTC == nil {
		d.DTC = NewDTCCodes()
	}
	return d, !bErr
}
