package evt55

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

	//reserved = 19

	//GTJDR Jamming indication
	GTJDR = 20

	//GTGSS GPS signal status
	GTGSS = 21

	//GTCRA Crash incident report
	GTCRA = 23

	//GTDOS Output status change with wave shape 1
	GTDOS = 25

	//GTGES Geo-fence event report
	GTGES = 26

	//GTSTR Vehicle enters into start status
	GTSTR = 28

	//GTSTP Vehicle enters into stop status
	GTSTP = 29

	//GTLSP Vehicle enters into long stop status
	GTLSP = 30

	//GTGPJ GPS Jamming status report
	GTGPJ = 31

	//GTRMD packet
	GTRMD = 32
)
