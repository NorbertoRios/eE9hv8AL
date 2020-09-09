package evt600

const (
	//GTPNA Power on report
	GTPNA = 1

	//GTPFA Power off report
	GTPFA = 2

	//GTMPN Connecting main power supply
	GTMPN = 3

	//GTMPF Disconnecting main power supply
	GTMPF = 4

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

	//GTUPD Update device software report
	GTUPD = 15

	//GTIDN Enter into idling status
	GTIDN = 16

	//GTIDF Leave idling status
	GTIDF = 17

	//GTDAT transparent data transmition
	GTDAT = 18

	//GTGSS GPS signal status
	GTGSS = 21

	//GTFLA Unusual fuel consumption alarm
	GTFLA = 22

	//GTSTR Vehicle enters into start status
	GTSTR = 23

	//GTSTP Vehicle enters into stop status
	GTSTP = 24

	//GTCRA --removed 2018-01-18
	GTCRA = 25

	//GTDOS Output status change with wave shape 1
	GTDOS = 27

	//GTGES Geo-fence event report
	GTGES = 28

	//GTLSP Vehicle enters into long stop status
	GTLSP = 29

	//GTTMP Temperature alarm
	GTTMP = 30

	//GTRMD roaming detection configuration
	GTRMD = 33

	//GTUPC update configuration OTA
	GTUPC = 39

	//GTGPJ GPS jamming status report
	GTGPJ = 40

	//GTVGN virtual ignition on
	GTVGN = 45

	//GTVGF virtual ignition off
	GTVGF = 46
)
