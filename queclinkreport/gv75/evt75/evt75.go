package evt75

const (

	//GTPNA Power on report
	GTPNA = 1

	//GTPFA Power off report
	GTPFA = 2

	//GTMPN Connecting main power supply
	GTMPN = 3

	//GTMPF Disconnecting main power supply
	GTMPF = 4

	//reserved = 5

	//GTBPL Backup battery low. 4 times report before power off
	GTBPL = 6

	//GTBTC Backup battery start charging report
	GTBTC = 7

	//GTSTC Backup battery stop charging report
	GTSTC = 8

	//GTSTT Device motion state indication when the motion state is changed
	GTSTT = 9

	//reserved = 10
	//reserved = 11

	//GTPDP GPRS connection establishment report
	GTPDP = 12

	//GTIGN Ignition on report
	GTIGN = 13

	//GTIGF Ignition off report
	GTIGF = 14

	//GTUPD device update packet
	GTUPD = 15

	//GTIDN Enter into idling status
	GTIDN = 16

	//GTIDF Leave idling status
	GTIDF = 17

	//GTDAT raw data packet
	GTDAT = 18
	//reserved = 19

	//GTJDR Jamming indication
	GTJDR = 20

	//GTGSS GPS signal status
	GTGSS = 21

	//GTFLA fuel level alarm
	GTFLA = 22

	//GTSTR Vehicle enters into start status
	GTSTR = 23

	//GTSTP Vehicle enters into stop status
	GTSTP = 24

	//GTCRA Crash report
	GTCRA = 25
	//reserved = 26

	//GTDOS Output status change with wave shape 1
	GTDOS = 27

	//GTGES Geo-fence event report
	GTGES = 28

	//GTLSP Vehicle enters into long stop status
	GTLSP = 29

	//GTTMP Temperature alarm
	GTTMP = 30

	//GTDTT packet
	GTDTT = 31

	//GTJDS packet
	GTJDS = 32

	//GTRMD packet
	GTRMD = 33

	//GTPHL packet
	GTPHL = 34

	//GTEXP packet
	GTEXP = 35

	//reserved = 36

	//GTUFS packet
	GTUFS = 37
	//GTFTP packet
	GTFTP = 38
)
