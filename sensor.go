package main

import (
	"github.com/sapienzaapps/seismocloud-client-go"
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
	leds, err = NewLEDs()
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	leds.StartLoading()

	// Load config if exists
	log.Info("Loading config")
	cfg = NewConfig()
	log.Info("Device Id:", cfg.GetDeviceId(), "Sigma:", cfg.GetSigma())

	// Connect to MQTT
	log.Info("Connecting to MQTT Server", MQTT_SERVER)
	scs, err = scsclient.NewClientV1(scsclient.ClientV1Options{
		ClientId:       cfg.GetDeviceId(),
		Server:         MQTT_SERVER,
		DeviceId:       cfg.GetDeviceId(),
		User:           "embedded",
		Pass:           "embedded",
		Logger:         log,
		ConfigCallback: cfg.NewConfigReceived,
		RebootCallback: cfg.RemoteReboot,
		UpdateCallback: cfg.UpdateCallback,
	})
	scs.SetSkipTLS(SKIP_TLS)
	for {
		err = scs.Connect()
		if err != nil {
			leds.StopLoading()
			leds.Yellow(true)

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
	leds.StartLoading()

	// Start local broadcaster
	log.Info("Starting LAN interface")
	StartLANInterface(cfg.GetDeviceId())

	// NTP sync
	servertime := scs.GetTime()
	if servertime.Unix()-time.Now().Unix() > 2 {
		log.Error("Local time is not synchronized")
		StopLANInterface()
		scs.Disconnect()
		leds.Green(false)
		leds.Yellow(true)
		leds.Red(true)
		os.Exit(-1)
	}
	servertimeLastCheck := *servertime

	// Init seismometer
	log.Info("Init Seismometer")
	seismometer, err := CreateNewSeismometer(cfg.GetSigma())
	if err != nil {
		panic(err)
	}
	cfg.RegisterSigmaCallback(seismometer.SetSigma)

	leds.StopLoading()
	leds.StartupBlink()
	leds.Green(true)
	log.Info("Ready")

	seismometer.StartSeismometer()

	running := true
	for running {
		// Check internet every X minutes
		if !scs.IsConnected() {
			leds.Green(false)
			leds.Yellow(true)
			leds.Red(false)

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
			leds.Green(true)
			leds.Yellow(false)
			leds.Red(false)
		}

		// Check server time
		if time.Now().Sub(servertimeLastCheck).Hours() >= 24 {
			servertime := scs.GetTime()
			if servertime.Unix()-time.Now().Unix() > 2 {
				log.Error("Local time is not synchronized")
				StopLANInterface()
				scs.Disconnect()
				leds.Green(false)
				leds.Yellow(true)
				leds.Red(true)
				os.Exit(-1)
			}
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
