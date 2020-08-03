package main

import (
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/config"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func sensor() {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	log.Info("Logging started")

	// Init LEDS
	err = ledset.Init()
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	ledset.StartLoading()

	// Load config if exists
	log.Info("Loading config")
	cfg, err = config.New()
	if err != nil {
		log.Error("error loading configuration: ", err)
		os.Exit(-3)
	}
	log.Info("Device Id:", cfg.GetDeviceId(), "Sigma:", cfg.GetSigma())

	// TODO: check for updates

	// Setup client
	log.Info("Connecting to MQTT Server", config.MQTT_SERVER)
	scs, err = scsclient.New(scsclient.ClientOptions{
		DeviceId:          cfg.GetDeviceId(),
		Model:             config.MODEL,
		Version:           config.VERSION,
		OnNewSigma:        onNewSigma,
		OnReboot:          onReboot,
		OnStreamCommand:   nil,
		OnProbeSpeedSet:   nil,
		OnTimeReceived:    onTimeReceived,
		SeismoCloudBroker: config.MQTT_SERVER,
		Username:          "embedded",
		Password:          "embedded",
	})
	if err != nil {
		panic(err)
	}

	// Try to connect indefinitely
	for {
		err = scs.Connect()
		if err != nil {
			// Connect error - retry in 10s if not interrupted
			ledset.StopLoading()
			_ = ledset.Yellow(true)

			log.Error(err.Error())
			select {
			case <-sigs:
				os.Exit(0)
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
	checkTime()
	servertimeLastCheck := time.Now()

	// Start local broadcaster
	log.Info("Starting LAN interface")
	StartLANInterface(cfg.GetDeviceId())

	// Init seismometer
	log.Info("Init Seismometer")
	seismometer, err := CreateNewSeismometer(cfg.GetSigma())
	if err != nil {
		panic(err)
	}

	ledset.StopLoading()
	_ = ledset.StartupBlink()
	_ = ledset.Green(true)
	log.Info("Ready")

	seismometer.StartSeismometer()

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
					log.Error(err.Error())
					select {
					case <-sigs:
						os.Exit(-2)
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

	seismometer.StopSeismometer()

	log.Info("Stopping LAN interface")
	StopLANInterface()

	log.Info("End")
}
