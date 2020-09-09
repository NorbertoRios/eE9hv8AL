package resp300

const (

	//GTRTLPNL description
	GTRTLPNL = 0

	//GTTOW motion sensors report
	GTTOW = 1
	//    reserved = 2

	//GTLBC report location by call
	GTLBC = 3

	//GTEPS external power supply report
	GTEPS = 4

	//GTDIS report for customizable digital input
	GTDIS = 5

	//GTIOB report for input port when changed
	GTIOB = 6

	//GTFRI fixed report information
	GTFRI = 7

	//GTGEO enter/exit geo-fence
	GTGEO = 8

	//GTSPD speeding alarm
	GTSPD = 9

	//GTSOS report send to backend server when input port is activated
	GTSOS = 10

	//GTRTL after AT+GTRTO start GPS and send message with current position to backend server
	GTRTL = 11

	//GTDOG the protocol watchdog reboot message
	GTDOG = 12

	//reserved = 13,

	//GTAIS Alarm for analog input voltage enters the alarm range
	GTAIS = 14

	//GTHBM If harsh behavior is detected, this message will be sent
	GTHBM = 15

	//GTIGL location message for ignition on
	GTIGL = 16

	//GTIDA Protect unauthorized device use
	GTIDA = 17

	//GTERI extended report
	GTERI = 18

	//GTGIN enter Geo-Fence report
	GTGIN = 20

	//GTGOT leave Geo-Fence report
	GTGOT = 21
)
