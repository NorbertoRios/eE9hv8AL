package gv75

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv75/evt75"
	"queclink-go/queclinkreport/gv75/resp75"
)

//GetReason returns reason for gv300 devices
func GetReason(message report.IMessage) int32 {
	switch message.MessageType() {
	case "+RSP", "+BSP":
		return getLocationReason(message)
	case "+EVT", "+BVT":
		return getEventCode(message)

	}
	return int32(6)
}

func getLocationReason(message report.IMessage) int32 {
	switch message.EventCode() {
	case resp75.GTDIS:
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
	case resp75.GTSPD:
		return translateSpeeding(message)
	case resp75.GTHBM:
		return translateHarsh(message)
	case resp75.GTTOW:
		return 112 //ongoing towing
	default:
		return 6 //periodical
	}
}

func getReportIDType(message report.IMessage) (byte, bool) {
	ir, found := message.GetValue(fields.ReportIDType)
	if !found {
		return 0, false
	}

	reportIDType, valid := ir.(byte)
	if !valid {
		return 0, false
	}
	return reportIDType, true
}

func decodeInputs(reportIDType byte) (int32, bool) {
	reportID := utils.GetHighestBits(reportIDType)
	reportType := utils.GetLowestBits(reportIDType)
	if reportID == 0 && reportType == 0 {
		return 2, true //IgnitionOff
	}
	if reportID == 0 && reportType == 1 {
		return 3, true //IgnitionOn
	}
	if reportID == 1 && reportType == 0 {
		return 24, true //Switch1Off
	}
	if reportID == 1 && reportType == 1 {
		return 20, true //Switch1On
	}
	return -1, false
}

func translateSpeeding(message report.IMessage) int32 {

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

func translateHarsh(message report.IMessage) int32 {
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
	default:
		return 6 //periodical
	}
}

func getEventCode(message report.IMessage) int32 {
	switch message.EventCode() {
	case evt75.GTPNA:
		return 0 //PowerUp
	case evt75.GTPFA:
		return 5 //PowerOff
	case evt75.GTMPN:
		return 0 //PowerUp
	case evt75.GTMPF:
		return 49 //MainPowerLost

	case evt75.GTBTC:
		return 6 //Periodical
	case evt75.GTSTC:
		return 6 //Periodical
	case evt75.GTSTT:
		return 6 //Periodical
	case evt75.GTPDP:
		return 6 //Periodical
	case evt75.GTIDN:
		return 16 //IdleTimer
	case evt75.GTJDR:
		return 108 //Jamming
	case evt75.GTSTR:
		return 29 //BeginMove
	case evt75.GTSTP:
		return 16 //IdleTimer
	case evt75.GTLSP:
		return 16 //IdleTimer
	case evt75.GTBPL:
		return 31 //PowerOffBatt
	case evt75.GTIGN:
		return 3 //IgnitionOn
	case evt75.GTIGF:
		return 2 //IgnitionOff
	case evt75.GTUPD:
		return 6 //Periodical
	case evt75.GTIDF:
		return 29 //BeginMove;
	case evt75.GTGSS:
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
	case evt75.GTDOS:
		return 6 //Periodical
	case evt75.GTGES:
		return 6 //Periodical
	case evt75.GTCRA:
		return 101 //CrashDetection
	default:
		return 6 //Periodical
	}
}
