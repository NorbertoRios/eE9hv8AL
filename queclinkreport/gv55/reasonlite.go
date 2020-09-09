package gv55

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv55/evt55"
	"queclink-go/queclinkreport/gv55/resp55"
)

//GetReasonLite returns reason for gv300 devices
func GetReasonLite(message report.IMessage) int32 {
	switch message.MessageType() {
	case "+RSP", "+BSP":
		return getLocationReasonLite(message)
	case "+EVT", "+BVT":
		return getEventCodeLite(message)

	}
	return int32(6)
}

func getLocationReasonLite(message report.IMessage) int32 {
	switch message.EventCode() {
	case resp55.GTDIS:
		{
			reportIDType, valid := getReportIDType(message)
			if !valid {
				return 6 //periodical
			}

			event, found := decodeInputs(reportIDType)
			if !found {
				return 6 //periodical
			}
			return event
		}
	case resp55.GTSPD:
		return translateSpeedingLite(message)
	case resp55.GTHBM:
		return translateHarshLite(message)
	case resp55.GTTOW:
		return 112 //ongoing towing
	default:
		return 6 //periodical
	}
}

func translateSpeedingLite(message report.IMessage) int32 {

	reportIDType, valid := getReportIDType(message)
	if !valid {
		return 6 //periodical
	}
	reportType := utils.GetLowestBits(reportIDType)
	switch reportType {
	case 1:
		return 18 //SpeedingStart
	case 0:
		return 19 //SpeedingStop
	default:
		return 6 //Periodical
	}
}

func translateHarshLite(message report.IMessage) int32 {
	reportIDType, valid := getReportIDType(message)
	if !valid {
		return 6 //periodical
	}
	reportType := utils.GetLowestBits(reportIDType)
	switch reportType {
	case 0:
		return 62 //Deceleration
	case 1:
		return 61 //Acceleration
	case 2, 3, 4: //Turn, Brake_Turn, Accelerate_Turn
		return 100 //HarshTurn
	default:
		return 6 //periodical
	}
}

func getEventCodeLite(message report.IMessage) int32 {
	switch message.EventCode() {
	case evt55.GTPNA:
		return 0 //PowerUp
	case evt55.GTPFA:
		return 5 //PowerOff
	case evt55.GTMPN:
		return 0 //PowerUp
	case evt55.GTMPF:
		return 49 //MainPowerLost
	case evt55.GTBTC:
		return 6 //Periodical
	case evt55.GTSTC:
		return 6 //Periodical
	case evt55.GTSTT:
		return 6 //Periodical
	case evt55.GTPDP:
		return 6 //Periodical
	case evt55.GTIDN:
		return 16 //IdleTimer
	case evt55.GTJDR:
		return 108 //Jamming
	case evt55.GTSTR:
		return 29 //BeginMove
	case evt55.GTSTP:
		return 16 //IdleTimer
	case evt55.GTLSP:
		return 16 //IdleTimer
	case evt55.GTBPL:
		return 31 //PowerOffBatt
	case evt55.GTIGN:
		return 3 //IgnitionOn
	case evt55.GTIGF:
		return 2 //IgnitionOff
	case evt55.GTUPD:
		return 6 //Periodical
	case evt55.GTIDF:
		return 29 //BeginMove;
	case evt55.GTGSS:
		{
			iv, f := message.GetValue(fields.GPSSignalStatus)
			if !f {
				return 6 //Periodical
			}
			gps, valid := iv.(byte)
			if !valid {
				return 6 //Periodical
			}
			switch gps {
			case 0:
				return 10 //GpsLost
			case 1:
				return 11 //GpsFound
			}
			return 6 //Periodical
		}
	case evt55.GTDOS:
		return 6 //Periodical
	case evt55.GTGES:
		return 6 //Periodical
	case evt55.GTGPJ:
		{
			iv, f := message.GetValue(fields.GPSJammingValue)
			if !f {
				return 6 //Periodical
			}
			gps, valid := iv.(byte)
			if !valid {
				return 6 //Periodical
			}
			switch gps {
			case 0, 1, 2:
				return 6 //periodical
			case 3:
				return 108 //Jamming
			}
			return 6 //periodical
		}
	default:
		return 6 //Periodical
	}
}
