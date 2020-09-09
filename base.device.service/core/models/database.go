package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var rawdb *gorm.DB

//InitializeConnections to mysql databases(device and backoffice database)
func InitializeConnections(connectionString string, mode string) error {
	var rawerr error
	rawdb, rawerr = gorm.Open("mysql", connectionString)
	switch mode {
	case "release":
		rawdb = rawdb.LogMode(false)
		break
	default:
		rawdb = rawdb.LogMode(true)
	}

	if rawerr != nil {
		panic("Error connecting to raw database:" + rawerr.Error())
	}
	return rawerr
}
