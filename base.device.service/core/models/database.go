package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var rawdb *gorm.DB

//InitializeConnections to mysql databases(device and backoffice database)
func InitializeConnections(connectionString string) error {
	var rawerr error
	rawdb, rawerr = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if rawerr != nil {
		panic("Error connecting to raw database:" + rawerr.Error())
	}
	return rawerr
}
