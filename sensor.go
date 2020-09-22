package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/config"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
)

func sensor() int {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	log.Info("Logging started")

	// Init LEDS
	err = ledset.Init()
	if err != nil {
		log.Error("error initializing LEDs: ", err)
		return -1
	}
	ledset.StartLoading()

	// Load config if exists
	log.Info("Loading config")
	cfg, err = config.New()
	if err != nil {
		log.Error("error loading configuration: ", err)
		return -3
	}
	log.Info("Device Id:", cfg.GetDeviceID(), "Sigma:", cfg.GetSigma())

	// TODO: check for updates

	// Init seismometer
	log.Info("Init Seismometer")
	seismometer, err := CreateNewSeismometer(cfg.GetSigma())
	if err != nil {
		log.Error("error initializing the seismometer: ", err)
		return -5
	}

	// Setup client
	log.Info("Connecting to MQTT Server", config.MqttServer)
	scs, err = scsclient.New(scsclient.ClientOptions{
		DeviceId:          cfg.GetDeviceID(),
		Model:             config.Model,
		Version:           config.Version,
		OnNewSigma:        onNewSigmaFunc(seismometer),
		OnReboot:          onReboot,
		OnStreamCommand:   nil,
		OnProbeSpeedSet:   nil,
		OnTimeReceived:    onTimeReceived,
		SeismoCloudBroker: config.MqttServer,
		Username:          "embedded",
		Password:          "embedded",
	})
	if err != nil {
		log.Error("error creating the seismocloud client: ", err)
		return -4
	}
	defer func() {
		err := scs.Close()
		if err != nil {
			log.Error("error closing the connection to seismocloud network: ", err)
		}
	}()

	// Try to connect indefinitely
	for {
		err = scs.Connect()
		if err != nil {
			log.Error("connection error: ", err)

			// Connect error - retry in 10s if not interrupted
			ledset.StopLoading()
			_ = ledset.Yellow(true)

			select {
			case <-sigs:
				return 0
			default:
				time.Sleep(10 * time.Second)
			}
		} else {
			break
		}
	}

	// Connect OK
	ledset.StartLoading()

	// NTP sync
	log.Info("Synchronizing time")
	checkTime()
	servertimeLastCheck := time.Now()

	// Start local broadcaster
	log.Info("Starting LAN interface")
	StartLANInterface(cfg.GetDeviceID())
	defer func() {
		log.Info("Stopping LAN interface")
		StopLANInterface()
	}()

	ledset.StopLoading()
	_ = ledset.StartupBlink()
	_ = ledset.Green(true)
	log.Info("Ready")

	seismometer.StartSeismometer()
	defer func() {
		seismometer.StopSeismometer()
	}()

	running := true
	for running {
		// Check internet every X minutes
		if !scs.IsConnected() {
			_ = ledset.Green(false)
			_ = ledset.Yellow(true)
			_ = ledset.Red(false)

			for {
				err = scs.Connect()
				if err != nil {
					log.Error("connection error: ", err)
					select {
					case <-sigs:
						return 0
					default:
						time.Sleep(10 * time.Second)
					}
				} else {
					break
				}
			}
			_ = ledset.Green(true)
			_ = ledset.Yellow(false)
			_ = ledset.Red(false)
		}

		// Check server time
		if time.Now().Sub(servertimeLastCheck).Hours() >= 24 {
			checkTime()
			servertimeLastCheck = time.Now()
		}

		// Check if termination is requested
		select {
		case <-sigs:
			running = false
		case <-time.After(1 * time.Second):
		}
	}
	return 0
}
