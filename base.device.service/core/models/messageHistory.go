package models

import (
	"encoding/json"
	"log"
	"time"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"

	"github.com/go-sql-driver/mysql"
)

//MessageHistory struct
type MessageHistory struct {
	ID              uint64    `gorm:"column:ID;primary_key"`
	DevID           string    `gorm:"column:DevId"`
	EntryData       []byte    `gorm:"column:EntryData"`
	ParsedEntryData []byte    `gorm:"column:ParsedEntryData"`
	Time            time.Time `gorm:"column:Time"`
	RecievedTime    time.Time `gorm:"column:RecievedTime"`
	ReportClass     string    `gorm:"column:ReportClass"`
	ReportType      int32     `gorm:"column:ReportType"`
	Reason          int32     `gorm:"column:Reason"`
	Latitude        float32   `gorm:"column:Latitude"`
	Longitude       float32   `gorm:"column:Longitude"`
	Speed           float32   `gorm:"column:Speed"`
	ValidFix        byte      `gorm:"column:ValidFix"`
	Altitude        float32   `gorm:"column:Altitude"`
	Heading         float32   `gorm:"column:Heading"`
	IgnitionState   byte      `gorm:"column:IgnitionState"`
	Odometer        int32     `gorm:"column:Odometer"`
	Satellites      int32     `gorm:"column:Satellites"`
	Supply          int32     `gorm:"column:Supply"`
	GPIO            byte      `gorm:"column:GPIO"`
	Relay           byte      `gorm:"column:Relay"`
}

//TableName for DeviceActivity model
func (h *MessageHistory) TableName() string {
	tableName := GetMessageHistoryTableName(h.DevID)
	return "raw_data." + tableName
}

//Save message to raw history table
func (h *MessageHistory) save() (uint64, error) {
	err := rawdb.Create(h).Error
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)
		if ok && merr.Number == 1146 {
			cerr := CreateMessageHistoryTable(GetMessageHistoryTableName(h.DevID))
			if cerr == nil {
				err = rawdb.Create(h).Error
			}
		}
	}
	return h.ID, err
}

//Save message to raw history table
func (h *MessageHistory) Save(message report.IMessage) (uint64, error) {
	if v, f := message.GetValue("RawData"); f {
		switch v.(type) {
		case string:
			{
				h.EntryData = []byte(v.(string))
				break
			}
		case []byte:
			{
				h.EntryData = v.([]byte)
				break
			}
		}
		delete(*message.GetData(), "RawData")
	}
	jMessage, jerr := json.Marshal(message)
	if jerr != nil {
		log.Fatalln("[Device] Error serialize history message:", jerr, "; Packet:", utils.InsertNth(utils.ByteToString(h.EntryData), 2, ' '))
		return 0, nil
	}
	h.ParsedEntryData = jMessage
	if v, f := message.GetValue("TimeStamp"); f {
		h.Time = v.(*utils.JSONTime).Time
	}

	h.DevID = message.GetStringValue("DevId", "")
	h.RecievedTime = time.Now().UTC()

	h.ReportType = message.EventCode()
	h.ReportClass = message.MessageType()

	if v, f := message.GetValue("Reason"); f {
		h.Reason = v.(int32)
	}

	if v, f := message.GetValue("Latitude"); f {
		h.Latitude = v.(float32)
	}
	if v, f := message.GetValue("Longitude"); f {
		h.Longitude = v.(float32)
	}
	if v, f := message.GetValue("Speed"); f {
		h.Speed = v.(float32)
	}
	if v, f := message.GetValue("GpsValidity"); f {
		h.ValidFix = v.(byte)
	}
	if v, f := message.GetValue("Altitude"); f {
		h.Altitude = v.(float32)
	}
	if v, f := message.GetValue("Heading"); f {
		h.Heading = v.(float32)
	}
	if v, f := message.GetValue("IgnitionState"); f {
		h.IgnitionState = v.(byte)
	}
	if v, f := message.GetValue("Odometer"); f {
		h.Odometer = v.(int32)
	}
	if v, f := message.GetValue("Satellites"); f {
		h.Satellites = v.(int32)
	}
	if v, f := message.GetValue("Supply"); f {
		h.Supply = v.(int32)
	}
	if v, f := message.GetValue("GPIO"); f {
		h.GPIO = v.(byte)
	}
	if v, f := message.GetValue("Relay"); f {
		h.Relay = v.(byte)
	}

	id, err := h.save()
	message.SetSourceID(id)
	return id, err
}

//DropMessageHistoryTable drop table from raw database
func DropMessageHistoryTable(tableName string) error {
	return rawdb.Exec("DROP TABLE " + tableName).Error
}

//GetMessageHistoryTableName returns last two digit of identity
func GetMessageHistoryTableName(identity string) string {
	return identity[len(identity)-2:]
}

//CreateMessageHistoryTable creates new table if table not exists detected
func CreateMessageHistoryTable(tableName string) error {
	return rawdb.Exec("CREATE TABLE IF NOT EXISTS  raw_data.`" + tableName + "` ( " +
		"`Id` bigint(20) NOT NULL AUTO_INCREMENT, " +
		"`DevId` varchar(100) NOT NULL, " +
		"`EntryData` blob, " +
		"`ParsedEntryData` blob, " +
		"`Time` datetime NOT NULL, " +
		"`RecievedTime` datetime NOT NULL, " +
		"`ReportClass` varchar(100) DEFAULT NULL, " +
		"`ReportType` int(11) DEFAULT NULL, " +
		"`Reason` varchar(5) DEFAULT NULL, " +
		"`Latitude` double DEFAULT NULL COMMENT 'degrees', " +
		"`Longitude` double DEFAULT NULL COMMENT 'degrees', " +
		"`Speed` double DEFAULT NULL, " +
		"`ValidFix` int(11) DEFAULT NULL, " +
		"`Altitude` double DEFAULT NULL, " +
		"`Heading` double DEFAULT NULL, " +
		"`IgnitionState` int(11) DEFAULT NULL, " +
		"`Odometer` int(10) DEFAULT NULL COMMENT 'm', " +
		"`Satellites` tinyint(3) unsigned DEFAULT NULL, " +
		"`Supply` int(10) DEFAULT NULL, " +
		"`GPIO` int(10) DEFAULT NULL COMMENT 'Input ports state', " +
		"`Relay` int(10) DEFAULT NULL COMMENT 'Output ports state', " +
		"`msg_id` binary(16) DEFAULT NULL, " +
		"`Extra` text, " +
		"`BatteryLow` double DEFAULT NULL, " +
		" PRIMARY KEY (`Id`,`Time`,`DevId`), " +
		"KEY `IX_RecievedTime` (`RecievedTime`,`DevId`) " +
		")" +
		"ENGINE = INNODB " +
		"AVG_ROW_LENGTH = 8192 " +
		"CHARACTER SET utf8 " +
		"COLLATE utf8_general_ci " +
		"PARTITION BY RANGE (to_days(Time)) " +
		"(" +
		"PARTITION p180201 VALUES LESS THAN (737091) ENGINE = InnoDB, " +
		"PARTITION p180301 VALUES LESS THAN (737119) ENGINE = InnoDB, " +
		"PARTITION p180401 VALUES LESS THAN (737150) ENGINE = InnoDB, " +
		"PARTITION p180501 VALUES LESS THAN (737180) ENGINE = InnoDB, " +
		"PARTITION p180601 VALUES LESS THAN (737211) ENGINE = InnoDB, " +
		"PARTITION p180701 VALUES LESS THAN (737241) ENGINE = InnoDB, " +
		"PARTITION p180801 VALUES LESS THAN (737272) ENGINE = InnoDB, " +
		"PARTITION p180901 VALUES LESS THAN (737303) ENGINE = InnoDB, " +
		"PARTITION p181001 VALUES LESS THAN (737333) ENGINE = InnoDB, " +
		"PARTITION p181101 VALUES LESS THAN (737364) ENGINE = InnoDB, " +
		"PARTITION p181201 VALUES LESS THAN (737394) ENGINE = InnoDB, " +
		"PARTITION p190101 VALUES LESS THAN (737425) ENGINE = InnoDB, " +
		"PARTITION p190201 VALUES LESS THAN (737456) ENGINE = InnoDB, " +
		"PARTITION p_cur VALUES LESS THAN MAXVALUE ENGINE = InnoDB " +
		");").Error
}
