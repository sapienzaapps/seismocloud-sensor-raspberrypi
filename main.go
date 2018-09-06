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

var log = logging.MustGetLogger("sensor")
var cfg Config
var leds LEDs
var scs scsclient.ClientV1

func main() {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Run as root")
		os.Exit(2)
	}
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
