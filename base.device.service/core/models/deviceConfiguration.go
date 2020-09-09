package models

import (
	"log"
	"time"

	"queclink-go/base.device.service/utils"
)

//DeviceConfig struct for device configuration
type DeviceConfig struct {
	ID        int32          `gorm:"column:cfgId;primary_key"`
	DevID     int32          `gorm:"column:devid; default:1"`
	Identity  string         `gorm:"column:devIdentity"`
	Command   string         `gorm:"column:cfgCommand"`
	CreatedAt time.Time      `gorm:"column:cfgCreated_at"`
	SentAt    utils.NullTime `gorm:"column:cfgSent_at"`
}

//TableName for DeviceConfig model
func (DeviceConfig) TableName() string {
	return "ats.tblDeviceConfig"
}

//Save device state cache to database
func (config *DeviceConfig) Save() {
	if rawdb.Model(&config).Where("devIdentity=? and isnull(cfgSent_At)", config.Identity).Update(&config).RowsAffected == 0 {
		rawdb.Create(&config)
	}
}

//DeleteAll not sent records for device
func (config *DeviceConfig) DeleteAll() {
	rawdb.Model(&config).Where("devIdentity=? and isnull(cfgSent_At)", config.Identity).Delete(&config)
}

//UpdateSentConfiguration update configuration was sent
func (config *DeviceConfig) UpdateSentConfiguration() {
	rawdb.Model(&config).Update("cfgSent_at", time.Now().UTC())
	log.Println("[Configuration]UpdateSentConfiguration: Configuration with id=", config.ID, " for device=", config.Identity, " updated")
}

//FindDeviceConfigByIdentity lookup not sent device configuration
func FindDeviceConfigByIdentity(identity string) (*DeviceConfig, bool) {
	d := &DeviceConfig{}
	err := rawdb.Where("devIdentity=? and isnull(cfgSent_At)", identity).Find(d).RecordNotFound()
	return d, !err
}
