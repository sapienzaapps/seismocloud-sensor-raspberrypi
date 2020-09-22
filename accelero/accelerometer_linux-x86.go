// +build !rpi

package accelero

type dummyAccelerometer struct {
}

// New creates a new instance of accelerometer for ADXL type sensors
func New() (Accelerometer, error) {
	return &dummyAccelerometer{}, nil
}

func (a *dummyAccelerometer) Start() error {
	return nil
}

func (a *dummyAccelerometer) GetAccelerometerName() string {
	return "DUMMY"
}

func (a *dummyAccelerometer) ProbeValueRaw() (AccelerometerData, error) {
	ret := AccelerometerData{0, 0, 0}

	return ret, nil
}

func (a *dummyAccelerometer) ProbeValue() (AccelerometerData, error) {
	return a.ProbeValueRaw()
}

func (a *dummyAccelerometer) Calibration() {
}

func (a *dummyAccelerometer) Stop() error {
	return nil
}

func (a *dummyAccelerometer) Shutdown() error {
	return nil
}
