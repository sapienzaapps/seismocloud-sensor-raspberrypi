package main

import (
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/accelero"
	"time"
)

func rawLogMain(absolute bool) {
	a, err := accelero.NewADXL345Accelerometer()
	if err != nil {
		panic(err)
	}

	err = a.Start()
	if err != nil {
		panic(err)
	}
	t := time.NewTicker(50 * time.Millisecond)

	avgX := NewRunningAvgFloat64()
	avgY := NewRunningAvgFloat64()
	avgZ := NewRunningAvgFloat64()
	if absolute {
		for i := 0; i < 100; i++ {
			probe, _ := a.ProbeValue()
			avgX.AddValue(probe.X)
			avgY.AddValue(probe.Y)
			avgZ.AddValue(probe.Z)
			<-t.C
		}
	}

	for {
		probe, err := a.ProbeValue()
		if err != nil {
			panic(err)
		}

		if absolute {
			probe.X = probe.X - avgX.GetAverage()
			probe.Y = probe.Y - avgY.GetAverage()
			probe.Z = probe.Z - avgZ.GetAverage()
		}

		fmt.Printf("%f\t%f\t%f\n", probe.X, probe.Y, probe.Z)
		<-t.C
	}
	err = a.Stop()
	if err != nil {
		panic(err)
	}
}
