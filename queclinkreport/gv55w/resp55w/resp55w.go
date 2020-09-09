package resp55w

const (

	//GTRTLPNL After AT+GTRTO start GPS and send message with current position to backend server
	GTRTLPNL = 0

	//GTTOW Motion sensors report
	GTTOW = 1
	//    reserved = 2

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

	//GTGEO Enter/Exit geo-fence
	GTGEO = 8

	//GTSPD Speeding alarm
	GTSPD = 9

	//GTSOS Report send to backend server when input port is activated
	GTSOS = 10

	//GTRTL After AT+GTRTO start GPS and send message with current position to backend server
	GTRTL = 11

	//GTDOG The protocol watchdog reboot message
	GTDOG = 12

	//reserved = 13
	//reserved = 14

	//GTHBM If harsh behavior is detected this message will be sent
	GTHBM = 15

	//GTIGL Location message for ignition on
	GTIGL = 16

	//GTGIN packet
	GTGIN = 25

	//GTGOT packet
	GTGOT = 26

	//GTVGL virtual ignition location
	GTVGL = 27
)
