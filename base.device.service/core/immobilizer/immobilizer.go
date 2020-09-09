package immobilizer

import (
	"strings"

	"queclink-go/base.device.service/report"
)

//IImmobilizer interface for device immobilizer
type IImmobilizer interface {
	SendAPIImmobilizerCommand(callbackID, port, state string, ttl int, trigger string, safetyOption bool) string
	Update()
	CheckStatusOfImmobilizerCommands()
	UpdateMessageSent(message report.IMessage)
}

//GetPortValue returns port value based on state and trigger
func GetPortValue(state, trigger string) (byte, bool) {
	state = strings.ToUpper(state)
	trigger = strings.ToUpper(trigger)
	switch {
	case state == "MOBILE" && trigger == "HIGH":
		return 0, true
	case state == "MOBILE" && trigger == "LOW":
		return 1, true
	case state == "ARMED" && trigger == "HIGH":
		return 1, true
	case state == "ARMED" && trigger == "LOW":
		return 0, true
	}
	return 0, false
}

//StrPortToInt converts string port to int
func StrPortToInt(port string) (byte, bool) {
	port = strings.ToUpper(port)
	switch port {
	case "OUT0":
		return 0, true
	case "OUT1":
		return 1, true
	case "OUT2":
		return 2, true
	case "OUT3":
		return 3, true
	}
	return 0, false
}
