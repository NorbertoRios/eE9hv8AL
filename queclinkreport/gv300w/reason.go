package gv300w

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv300w/evt300w"
	"queclink-go/queclinkreport/gv300w/resp300w"
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
	case resp300w.GTDIS:
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
	case resp300w.GTSPD:
		return translateSpeeding(message)
	case resp300w.GTHBM:
		return translateHarsh(message)
	case resp300w.GTTOW:
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
	case 0:
		return 18 //SpeedingStart
	case 1:
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
	case evt300w.GTPNA:
		return 0 //PowerUp
	case evt300w.GTPFA:
		return 5 //PowerOff
	case evt300w.GTMPN:
		return 0 //PowerUp
	case evt300w.GTMPF:
		return 49 //MainPowerLost
	case evt300w.GTBTC:
		return 6 //Periodical
	case evt300w.GTSTC:
		return 6 //Periodical
	case evt300w.GTSTT:
		return 6 //Periodical
	case evt300w.GTPDP:
		return 6 //Periodical
	case evt300w.GTIDN:
		return 16 //IdleTimer
	case evt300w.GTJDR:
		return 108 //Jamming
	case evt300w.GTSTR:
		return 29 //BeginMove
	case evt300w.GTSTP:
		return 16 //IdleTimer
	case evt300w.GTLSP:
		return 16 //IdleTimer
	case evt300w.GTBPL:
		return 31 //PowerOffBatt
	case evt300w.GTIGN:
		return 3 //IgnitionOn
	case evt300w.GTIGF:
		return 2 //IgnitionOff
	case evt300w.GTVGN:
		return 3 //IgnitionOn
	case evt300w.GTVGF:
		return 2 //IgnitionOff
	case evt300w.GTUPD:
		return 6 //Periodical
	case evt300w.GTIDF:
		return 29 //BeginMove;
	case evt300w.GTGSS:
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
	case evt300w.GTDOS:
		return 6 //Periodical
	case evt300w.GTGES:
		return 6 //Periodical
	default:
		return 6 //Periodical
	}
}
