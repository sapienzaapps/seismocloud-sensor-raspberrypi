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
		partialCMA:     utils.NewRunningAvgFloat64(),
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

func (s *seismometerImpl) preFill() error {
	oneSecondMs := time.Now()
	for time.Since(oneSecondMs) < time.Second {
		probe, err := s.accelerometer.ProbeValue()
		if err != nil {
			return err
		}
		s.partialCMA.AddValue(probe.GetTotalVector())
	}
	s.rotateThreshold()
	return nil
}

func (s *seismometerImpl) seismometer() {
	log.Debugf("Starting accelerometer %s", s.accelerometer.GetAccelerometerName())
	err := s.accelerometer.Start()
	if err != nil {
		log.Error("error starting the accelerometer: ", err)
		return
	}

	log.Debug("Calibrating accelerometer")
	s.accelerometer.Calibration()
	defer func() {
		_ = s.accelerometer.Stop()
		_ = s.accelerometer.Shutdown()
	}()

	log.Debug("Pre-fill threshold")
	// Pre-fill the threshold
	err = s.preFill()
	if err != nil {
		log.Error("error reading accelerometer values: ", err)
		return
	}

	log.Debug("Starting probe")
	lastThresholdUpdate := time.Now()
	lastQuake := time.Unix(0, 0)
	// TODO: adjust with the speed limit from MQTT
	t := time.NewTicker(50 * time.Millisecond)
	for s.active {
		if time.Since(lastThresholdUpdate) > 15*time.Minute {
			s.rotateThreshold()
			lastThresholdUpdate = time.Now()
		}

		probe, err := s.accelerometer.ProbeValue()
		if err != nil {
			log.Error("error reading accelerometer values: ", err)
			break
		}

		inQuake := time.Since(lastQuake) < 5*time.Second
		probeValue := probe.GetTotalVector()
		if !inQuake && probeValue > s.quakeThreshold {
			log.Debugf("New Event: v:%f - thr:%f - iter:%f - avg:%f - stddev:%f",
				probeValue, s.quakeThreshold, s.sigma, s.lastCMA.GetAverage(), s.lastCMA.GetStandardDeviation())

			_ = ledset.Red(true)
			lastQuake = time.Now()

			err = scs.Quake(time.Now(), probe.X, probe.Y, probe.Z)
			if err != nil {
				log.Warning("error sending quake signal: ", err)
			}
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
}

func (s *seismometerImpl) StopSeismometer() {
	s.active = false
}
