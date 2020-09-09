package evt55w

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

	//GTJDR Jamming indication
	GTJDR = 20

	//GTGSS GPS signal status
	GTGSS = 21

	//GTCRA Crash report
	GTCRA = 23

	//GTDOS Output status change with wave shape 1
	GTDOS = 25

	//GTGES Geo-fence event report
	GTGES = 26

	//GTSTR start packet
	GTSTR = 28

	//GTSTP stop packet
	GTSTP = 29

	//GTLSP packet
	GTLSP = 30

	//GTRMD roaming packet
	GTRMD = 32

	//GTJDS jamming packet
	GTJDS = 33

	//GTUPC packet
	GTUPC = 36

	//GTVGN virtual ignition on
	GTVGN = 39

	//GTVGF  virtual ignition off
	GTVGF = 40

	//GTPNR packet
	GTPNR = 41

	//GTPFR packet
	GTPFR = 42
)
