module git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi

go 1.14

replace github.com/jrcichra/gophidgets v0.0.0-20200812000151-517b9885f2a8 => github.com/sapienzaapps/gophidgets v0.0.0-20200924153930-a64521ae3d90

require (
	git.sapienzaapps.it/seismocloud/seismocloud-client-go v0.0.0-20200924142525-c922caba6d74
	github.com/jrcichra/gophidgets v0.0.0-20200812000151-517b9885f2a8
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/satori/go.uuid v1.2.0
	periph.io/x/periph v3.6.4+incompatible
)
