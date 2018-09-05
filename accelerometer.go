package main

import "math"

type Accelerometer interface {
	Start() error
	ProbeValue() (AccelerometerData, error)
	GetAccelerometerName() string
	Stop() error
	Shutdown() error
}

type AccelerometerData struct {
	X float64
	Y float64
	Z float64
}

func (a *AccelerometerData) GetTotalVector() float64 {
	return math.Sqrt(math.Pow(a.X, 2) + math.Pow(a.Y, 2) + math.Pow(a.Z, 2))
}
