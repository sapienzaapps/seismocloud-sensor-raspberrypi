package accelero

import (
	"math"
)

// Accelerometer represents a generic accelerometer
type Accelerometer interface {
	Start() error
	Calibration()
	ProbeValue() (AccelerometerData, error)
	ProbeValueRaw() (AccelerometerData, error)
	GetAccelerometerName() string
	Stop() error
	Shutdown() error
}

// AccelerometerData represents a 3D probe
type AccelerometerData struct {
	X float64
	Y float64
	Z float64
}

// GetTotalVector returns the modulus of the sum vector
func (a *AccelerometerData) GetTotalVector() float64 {
	return math.Sqrt(math.Pow(a.X, 2) + math.Pow(a.Y, 2) + math.Pow(a.Z, 2))
}

// Clone returns a clone of the accelerometer data
func (a *AccelerometerData) Clone() AccelerometerData {
	return AccelerometerData{a.X, a.Y, a.Z}
}

// Sub subtracts the given value from the target one
func (a *AccelerometerData) Sub(data AccelerometerData) {
	a.X = a.X - data.X
	a.Y = a.Y - data.Y
	a.Z = a.Z - data.Z
}
