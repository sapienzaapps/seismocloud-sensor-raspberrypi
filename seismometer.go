package main

import (
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/accelero"
	"time"
)

type Seismometer interface {
	StartSeismometer()
	SetSigma(float64)
	StopSeismometer()
}

type Detection struct {
	Timestamp     time.Time
	Acceleration  float64
	OverThreshold bool
}

func CreateNewSeismometer(sigma float64) (Seismometer, error) {
	s, e := accelero.NewADXL345Accelerometer()
	if e != nil {
		return nil, e
	}

	e = s.Start()
	if e != nil {
		return nil, e
	}

	t := time.NewTicker(50 * time.Millisecond)

	avgX := NewRunningAvgFloat64()
	avgY := NewRunningAvgFloat64()
	avgZ := NewRunningAvgFloat64()
	for i := 0; i < 100; i++ {
		probe, _ := s.ProbeValue()
		avgX.AddValue(probe.X)
		avgY.AddValue(probe.Y)
		avgZ.AddValue(probe.Z)
		<-t.C
	}

	s.Stop()

	return &SeismometerImpl{
		quakeThreshold: 0,
		active:         false,
		accelerometer:  s,
		sigma:          sigma,
		runningAVG:     NewRunningAvgFloat64(),
		normX:          avgX.GetAverage(),
		normY:          avgY.GetAverage(),
		normZ:          avgZ.GetAverage(),
	}, nil
}

type SeismometerImpl struct {
	accelerometer  accelero.Accelerometer
	runningAVG     RunningAvgFloat64
	quakeThreshold float64
	sigma          float64
	active         bool
	inEvent        bool
	lastEventWas   int64

	normX float64
	normY float64
	normZ float64
}

func (s *SeismometerImpl) StartSeismometer() {
	s.active = true
	go s.seismometer()
}

func (s *SeismometerImpl) SetSigma(sigma float64) {
	s.sigma = sigma
}

func (s *SeismometerImpl) seismometer() {
	t := time.NewTicker(50 * time.Millisecond)

	fmt.Println("Seismometer started")
	s.accelerometer.Start()
	for s.active {
		detectionAVG := s.runningAVG.GetAverage()
		detectionStdDev := s.runningAVG.GetStandardDeviation()

		probe, err := s.accelerometer.ProbeValue()
		if err != nil {
			panic(err)
		}
		probe.X = probe.X - s.normX
		probe.Y = probe.Y - s.normY
		probe.Z = probe.Z - s.normZ

		accel := probe.GetTotalVector()
		record := Detection{
			Timestamp:     time.Now(),
			Acceleration:  accel,
			OverThreshold: accel > s.quakeThreshold,
		}

		s.runningAVG.AddValue(record.Acceleration)
		s.quakeThreshold = s.runningAVG.GetAverage() + (s.runningAVG.GetStandardDeviation() * s.sigma)

		if s.inEvent && time.Now().Unix()-s.lastEventWas >= 5 {
			_ = ledset.Red(false)
			s.inEvent = false
		} else if s.inEvent && time.Now().Unix()-s.lastEventWas < 5 {
			continue
		}

		if record.OverThreshold && !s.inEvent && s.runningAVG.Elements() > 2 {
			log.Debugf("New Event: v:%f - thr:%f - iter:%f - avg:%f - stddev:%f", record.Acceleration, s.quakeThreshold,
				s.sigma, detectionAVG, detectionStdDev)

			_ = ledset.Red(true)

			s.inEvent = true
			s.lastEventWas = time.Now().Unix()

			scs.Quake(time.Now(), probe.X, probe.Y, probe.Z)
		}

		<-t.C
	}
	t.Stop()
	s.accelerometer.Stop()
}

func (s *SeismometerImpl) StopSeismometer() {
	s.active = false
}
