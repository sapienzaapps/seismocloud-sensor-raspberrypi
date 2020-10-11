// +build phidget

package accelero

import (
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/utils"
	"github.com/jrcichra/gophidgets/phidgets"
)

type phidgetAccelerometer struct {
	accelerometer phidgets.PhidgetAccelerometer

	norm AccelerometerData
	last AccelerometerData
}

// New creates a new instance of a dummy accelerometer
func New() (Accelerometer, error) {

	// TODO: check if this permission is always needed, or there is a way to grant permissions over the sensor
	// (it should exists)
	if os.Geteuid() != 0 {
		return nil, errors.New("Can't access to sensor as non-root user")
	}

	accelero := phidgets.PhidgetAccelerometer{}
	accelero.Create()

	// TODO
	accelero.SetDeviceSerialNumber(429840)
	//accelero.SetHubPort(0)
	//accelero.SetIsRemote(true)

	return &phidgetAccelerometer{
		accelerometer: accelero,
	}, nil
}

func (a *phidgetAccelerometer) Start() error {
	err := a.accelerometer.OpenWaitForAttachment(2000)
	// TODO: parametrizzare dall'esterno
	a.accelerometer.SetDataInterval(50)
	return err
}

func (a *phidgetAccelerometer) GetAccelerometerName() string {
	return "Phidget"
}

func (a *phidgetAccelerometer) ProbeValueRaw() (AccelerometerData, error) {
	accelData, err := a.accelerometer.GetAcceleration()
	if err != nil {
		return AccelerometerData{}, err
	}

	ret := AccelerometerData{
		float64(accelData[0]),
		float64(accelData[1]),
		float64(accelData[2]),
	}

	// Calculate the difference between values - the phidget returns the absolute acceleration
	ret.Sub(a.last)
	a.last = AccelerometerData{
		float64(accelData[0]),
		float64(accelData[1]),
		float64(accelData[2]),
	}

	return ret, nil
}

func (a *phidgetAccelerometer) ProbeValue() (AccelerometerData, error) {
	ret, err := a.ProbeValueRaw()
	if err != nil {
		return ret, err
	}

	// Normalize with calibration data
	ret.Sub(a.norm)

	return ret, nil
}

func (a *phidgetAccelerometer) Calibration() {
	t := time.NewTicker(50 * time.Millisecond)

	// Calibration
	avgX := utils.NewRunningAvgFloat64()
	avgY := utils.NewRunningAvgFloat64()
	avgZ := utils.NewRunningAvgFloat64()
	startms := time.Now()
	for time.Since(startms) < 10*time.Second {
		probe, _ := a.ProbeValue()
		avgX.AddValue(probe.X)
		avgY.AddValue(probe.Y)
		avgZ.AddValue(probe.Z)
		<-t.C
	}

	a.norm = AccelerometerData{
		avgX.GetAverage(),
		avgY.GetAverage(),
		avgZ.GetAverage(),
	}
}

func (a *phidgetAccelerometer) Stop() error {
	return nil
}

func (a *phidgetAccelerometer) Shutdown() error {
	a.accelerometer.Close()
	return nil
}
