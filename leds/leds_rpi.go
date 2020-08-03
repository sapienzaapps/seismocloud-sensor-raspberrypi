// +build rpi

package leds

import (
	"io/ioutil"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"strings"
	"time"
)

const (
	GPIO0      = "17" // Green
	GPIO1      = "18" // Yellow
	GPIO2      = "21" // Red (revision 1)
	GPIO2_REV2 = "27" // Red (revision 2+)
)

type LEDImpl struct {
	green  gpio.PinIO
	yellow gpio.PinIO
	red    gpio.PinIO

	loadingchan chan interface{}
}

func isRPiRev1() bool {
	buf, err := ioutil.ReadFile("/proc/device-tree/model")
	if err != nil {
		return false
	}
	model := string(buf)
	return strings.HasSuffix(model, "Rev 1")
}

func (l *LEDImpl) Init() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	l.green = gpioreg.ByName(GPIO0)
	l.yellow = gpioreg.ByName(GPIO1)
	if isRPiRev1() {
		l.red = gpioreg.ByName(GPIO2)
	} else {
		l.red = gpioreg.ByName(GPIO2_REV2)
	}
	l.loadingchan = nil

	return nil
}

func (l *LEDImpl) Green(s bool) error {
	return l.green.Out(gpio.Level(s))
}

func (l *LEDImpl) Yellow(s bool) error {
	return l.yellow.Out(gpio.Level(s))
}

func (l *LEDImpl) Red(s bool) error {
	return l.red.Out(gpio.Level(s))
}

func (l *LEDImpl) StartupBlink() error {
	var err error
	for i := 0; i < 10; i++ {
		err = l.green.Out(gpio.High)
		if err != nil {
			return err
		}
		err = l.yellow.Out(gpio.High)
		if err != nil {
			return err
		}
		err = l.red.Out(gpio.High)
		if err != nil {
			return err
		}
		time.Sleep(20 * time.Millisecond)

		// Assume that the error was already fired in previous calls
		_ = l.green.Out(gpio.Low)
		_ = l.yellow.Out(gpio.Low)
		_ = l.red.Out(gpio.Low)
		time.Sleep(20 * time.Millisecond)
	}
	return nil
}

func (l *LEDImpl) StartLoading() {
	if l.loadingchan != nil {
		return
	}
	l.loadingchan = make(chan interface{}, 1)
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		i := 0
		for {
			i = i % 3
			switch i {
			default:
				fallthrough
			case 0:
				_ = l.green.Out(gpio.High)
				_ = l.yellow.Out(gpio.Low)
				_ = l.red.Out(gpio.Low)
			case 1:
				_ = l.green.Out(gpio.Low)
				_ = l.yellow.Out(gpio.High)
				_ = l.red.Out(gpio.Low)
			case 2:
				_ = l.green.Out(gpio.Low)
				_ = l.yellow.Out(gpio.Low)
				_ = l.red.Out(gpio.High)
			}

			select {
			case <-l.loadingchan:
				l.loadingchan = nil
				return
			case <-t.C:
				i++
			}
		}
	}()
}

func (l *LEDImpl) StopLoading() {
	l.loadingchan <- 1
}
