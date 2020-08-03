package main

import (
	"flag"
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/config"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/leds"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logging.MustGetLogger("sensor")
var cfg config.Config
var ledset leds.LEDSet
var scs scsclient.Client

func main() {
	// TODO: force this only in platforms where LED access is root-only
	// TODO: document a way to grant permissions to runtime user
	if os.Geteuid() != 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Run as root")
		os.Exit(2)
	}

	ledset = leds.New()

	showDeviceId := flag.Bool("showdeviceid", false, "Show the device ID and exit")
	testLanDiscovery := flag.Bool("testlandiscovery", false, "Boot only the LAN discovery handler")
	rawLog := flag.Bool("rawlog", false, "Dump raw Accelerometer data")
	rawLogAbsolute := flag.Bool("absolute", false, "Use absolute values on raw logs")
	ledTest := flag.Bool("testleds", false, "Test LEDs")
	stage2update := flag.String("stage2update", "", "Internal flag")

	flag.Parse()

	if *stage2update != "" {
		err := updateStage2(*stage2update)
		if err != nil {
			log.Error(err.Error())
		}
	}

	if *showDeviceId {
		cfg, err := config.New()
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Println(cfg.GetDeviceId())
		}
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
		err := ledset.Init()
		if err != nil {
			panic(err)
		}
		ledset.StartLoading()
		time.Sleep(10 * time.Second)
		ledset.StopLoading()
	} else {
		sensor()
	}
}
