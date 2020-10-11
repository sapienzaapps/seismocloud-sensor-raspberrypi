package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/accelero"
)

func rawLogMain(absolute bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	a, err := accelero.New()
	if err != nil {
		panic(err)
	}

	err = a.Start()
	if err != nil {
		panic(err)
	}
	a.Calibration()
	t := time.NewTicker(50 * time.Millisecond)

	fmt.Println("UNIX Time\tX\tY\tZ\tTotal vector")

	running := true
	for running {
		var probe accelero.AccelerometerData
		var err error
		if absolute {
			probe, err = a.ProbeValueRaw()
		} else {
			probe, err = a.ProbeValue()
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("%0.3f\t%f\t%f\t%f\t%f\n", float64(time.Now().UnixNano())/float64(time.Second), probe.X, probe.Y, probe.Z, probe.GetTotalVector())

		select {
		case <-sigs:
			// External signal received, exiting
			running = false
			t.Stop()
		case <-t.C:
			// Continue
		}
	}
	err = a.Stop()
	if err != nil {
		panic(err)
	}
}
