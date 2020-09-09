package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/devicetypes"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv300"
	"queclink-go/queclinkreport/gv300w"
	"queclink-go/queclinkreport/gv55"
	"queclink-go/queclinkreport/gv55w"
	"queclink-go/queclinkreport/gv600"
	"queclink-go/queclinkreport/gv75"
)

//GetReason returns reason for message
func GetReason(message report.IMessage) int32 {
	idt, found := message.GetValue(fields.DeviceType)
	if !found {
		return 6 //periodical
	}

	deviceType, valid := idt.(byte)
	if !valid {
		return 6 //periodical
	}

	switch deviceType {
	case devicetypes.GV320:
		return gv300.GetReason(message)

	case devicetypes.GV55, devicetypes.GV55N:
		return gv55.GetReason(message)

	case devicetypes.GV55Lite, devicetypes.GV55NLite:
		return gv55.GetReasonLite(message)

	case devicetypes.GV75, devicetypes.GV75W:
		return gv75.GetReason(message)

	case devicetypes.GV55W:
		return gv55w.GetReason(message)

	case devicetypes.GV600W:
		return gv600.GetReason(message)
	case devicetypes.GV300W:
		return gv300w.GetReason(message)
	default:
		return gv55.GetReason(message)
	}
}
