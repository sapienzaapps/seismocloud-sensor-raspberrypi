package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/config"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/leds"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
)

var log = logging.MustGetLogger("sensor")
var cfg config.Interface
var ledset leds.LEDSet = leds.New()
var scs scsclient.Client

func main() {
	// TODO: force this only in platforms where LED/sensors access is root-only
	// TODO: document a way to grant permissions to runtime user
	/*if os.Geteuid() != 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Run as root")
		os.Exit(2)
	}*/

	showDeviceID := flag.Bool("showdeviceid", false, "Show the device ID and exit")
	testLanDiscovery := flag.Bool("testlandiscovery", false, "Boot only the LAN discovery handler")
	rawLog := flag.Bool("rawlog", false, "Dump raw Accelerometer data")
	rawLogAbsolute := flag.Bool("absolute", false, "Use absolute values on raw logs")
	ledTest := flag.Bool("testleds", false, "Test LEDs")

	flag.Parse()

	if *showDeviceID {
		// Dumps only the device ID
		cfg, err := config.New()
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Println(cfg.GetDeviceID())
		}
	} else if *testLanDiscovery {
		// Test LAN discovery feature (no sensor activity)
		fmt.Println("Starting LAN discovery")
		StartLANInterface(uuid.Nil)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		<-sigs

		StopLANInterface()
	} else if *rawLog {
		// Dumps the raw sensor log (no sensor activity)
		rawLogMain(*rawLogAbsolute)
	} else if *ledTest {
		// Test LEDs (no sensor activity)
		err := ledset.Init()
		if err != nil {
			panic(err)
		}
		ledset.StartLoading()
		time.Sleep(10 * time.Second)
		ledset.StopLoading()
	} else {
		// Normal program flow
		sensor()
	}
}
