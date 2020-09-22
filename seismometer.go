package main

import (
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/accelero"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/utils"
)

// Seismometer describe a simple seismometer interface
type Seismometer interface {
	StartSeismometer()
	SetSigma(float64)
	StopSeismometer()
}

type seismometerImpl struct {
	accelerometer  accelero.Accelerometer
	quakeThreshold float64
	sigma          float64
	lastCMA        utils.RunningAvgFloat64
	partialCMA     utils.RunningAvgFloat64

	active bool
}

// CreateNewSeismometer returns a new seismometer instance
func CreateNewSeismometer(sigma float64) (Seismometer, error) {
	s, e := accelero.New()
	if e != nil {
		return nil, e
	}

	return &seismometerImpl{
		quakeThreshold: 0,
		active:         false,
		accelerometer:  s,
		sigma:          sigma,
		lastCMA:        utils.NewRunningAvgFloat64(),
	}, nil
}

func (s *seismometerImpl) StartSeismometer() {
	s.active = true
	go s.seismometer()
}

func (s *seismometerImpl) SetSigma(sigma float64) {
	s.sigma = sigma
	s.rotateThreshold()
}

func (s *seismometerImpl) rotateThreshold() {
	s.lastCMA = s.partialCMA
	s.partialCMA = utils.NewRunningAvgFloat64()
	s.quakeThreshold = s.lastCMA.GetAverage() + (s.lastCMA.GetStandardDeviation() * s.sigma)
}

func (s *seismometerImpl) preFill() {
	oneSecondMs := time.Now()
	for time.Since(oneSecondMs) < time.Second {
		probe, err := s.accelerometer.ProbeValue()
		if err != nil {
			// TODO: log and terminate gracefully
			panic(err)
		}
		s.partialCMA.AddValue(probe.GetTotalVector())
	}
	s.rotateThreshold()
}

func (s *seismometerImpl) seismometer() {
	// TODO: adjust with the speed limit from MQTT
	t := time.NewTicker(50 * time.Millisecond)

	s.accelerometer.Start()
	s.accelerometer.Calibration()

	// Pre-fill the threshold
	s.preFill()

	lastThresholdUpdate := time.Now()
	lastQuake := time.Unix(0, 0)
	for s.active {
		if time.Since(lastThresholdUpdate) > 15*time.Minute {
			s.rotateThreshold()
			lastThresholdUpdate = time.Now()
		}

		probe, err := s.accelerometer.ProbeValue()
		if err != nil {
			// TODO: log and terminate gracefully
			panic(err)
		}

		inQuake := time.Since(lastQuake) < 5*time.Second
		probeValue := probe.GetTotalVector()
		if !inQuake && probeValue > s.quakeThreshold {
			log.Debugf("New Event: v:%f - thr:%f - iter:%f - avg:%f - stddev:%f",
				probeValue, s.quakeThreshold, s.sigma, s.lastCMA.GetAverage(), s.lastCMA.GetStandardDeviation())

			_ = ledset.Red(true)

			scs.Quake(time.Now(), probe.X, probe.Y, probe.Z)
			// QUAKE
		} else if !inQuake {
			_ = ledset.Red(false)
			// End quake period
		}

		// TODO: stream data

		s.partialCMA.AddValue(probeValue)

		<-t.C
	}
	t.Stop()
	s.accelerometer.Stop()
}

func (s *seismometerImpl) StopSeismometer() {
	s.active = false
}
