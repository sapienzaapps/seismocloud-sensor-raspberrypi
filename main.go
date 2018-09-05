package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"github.com/sapienzaapps/seismocloud-client-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MQTT_SERVER = "tcp://mqtt.seismocloud.com:1883"

var log = logging.MustGetLogger("sensor")
var cfg Config
var leds LEDs
var scs scsclient.SCSClientOldProtocol

func main() {
	showDeviceId := flag.Bool("showdeviceid", false, "Show the device ID and exit")
	testLanDiscovery := flag.Bool("testlandiscovery", false, "Boot only the LAN discovery handler")
	rawLog := flag.Bool("rawlog", false, "Dump raw Accelerometer data")
	rawLogAbsolute := flag.Bool("absolute", false, "Use absolute values on raw logs")
	ledTest := flag.Bool("testleds", false, "Test LEDs")

	flag.Parse()

	if *showDeviceId {
		cfg := NewConfig()
		fmt.Println(cfg.GetDeviceId())
		return
	} else if *testLanDiscovery {
		fmt.Println("Starting LAN discovery")
		StartLANInterface("test")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		<-sigs

		StopLANInterface()
	} else if *rawLog {
		rawLogMain(*rawLogAbsolute)
	} else if *ledTest {
		leds, err := NewLEDs()
		if err != nil {
			panic(err)
		}
		leds.StartLoading()
		time.Sleep(10 * time.Second)
		leds.StopLoading()
	} else {
		sensor()
	}
}

func sensor() {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	log.Info("Logging started")

	// Init LEDS
	leds, err = NewLEDs()
	if err != nil {
		// TODO: handle
		log.Error(err.Error())
		os.Exit(-1)
	}
	leds.StartLoading()

	// Load config if exists
	log.Info("Loading config")
	cfg = NewConfig()
	log.Info("Device Id: ", cfg.GetDeviceId(), " Sigma: ", cfg.GetSigma())

	// Connect to MQTT
	log.Info("Connecting to MQTT Server ", MQTT_SERVER)
	scs = scsclient.NewSCSClientOldProtocol(log, cfg.NewConfigReceived, cfg.RemoteReboot, cfg.UpdateCallback)
	for {
		err = scs.Connect(cfg.GetDeviceId(), MQTT_SERVER, cfg.GetDeviceId(), "embedded", "embedded")
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
	scs.Alive()

	// Start local broadcaster
	log.Info("Starting LAN interface")
	StartLANInterface(cfg.GetDeviceId())

	// NTP sync
	scs.GetTime()
	// TODO

	// Init seismometer
	log.Info("Init Seismometer")
	seismometer, err := CreateNewSeismometer(cfg.GetSigma())
	cfg.RegisterSigmaCallback(seismometer.SetSigma)
	if err != nil {
		panic(err)
	}

	leds.StopLoading()
	leds.StartupBlink()
	leds.Green(true)
	log.Info("Ready")

	seismometer.StartSeismometer()

	running := true
	for running {
		// Check internet every X minutes
		// TODO

		// Update NTP every X minutes
		// TODO

		// Check if termination is needed
		select {
		case <-sigs:
			running = false
		default:
			time.Sleep(1 * time.Second)
		}
	}

	seismometer.StopSeismometer()

	log.Info("Stopping LAN interface")
	StopLANInterface()

	log.Info("End")
}
