package resp350

const (
	//GTPNL Power on location
	GTPNL = 0

	//GTTOW Motion sensors report
	GTTOW = 1
	//   reserved = 2

	//GTLBC Report location by call
	GTLBC = 3

	//GTEPS External power supply report
	GTEPS = 4

	//GTDIS Report for customizable digital input
	GTDIS = 5

	//GTIOB Report for input port when changed
	GTIOB = 6

	//GTFRI Fixed report information
	GTFRI = 7

	//GTSPD Speeding alarm
	GTSPD = 9

	//GTSOS Report send to backend server when input port is activated
	GTSOS = 10

	//GTRTL After AT+GTRTO start GPS and send message with current position to backend server
	GTRTL = 11

	//GTDOG The protocol watchdog reboot message
	GTDOG = 12

	//reserved = 13

	//GTAIS Analog input port setting
	GTAIS = 14

	//GTHBM If harsh behavior is detected  this message will be sent
	GTHBM = 15

	//GTIGL Location message for ignition on
	GTIGL = 16

	//GTIDA ID authentication message
	GTIDA = 17

	//GTERI Extended report information
	GTERI = 18

	//GTGIN geofence in report
	GTGIN = 20

	//GTGOT geofence out report
	GTGOT = 21
)
