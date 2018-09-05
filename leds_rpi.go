// +build rpi

package main

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

func NewLEDs() (LEDs, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}
	l := LEDImpl{}

	l.green = gpioreg.ByName(GPIO0)
	l.yellow = gpioreg.ByName(GPIO1)
	if isRPiRev1() {
		l.red = gpioreg.ByName(GPIO2)
	} else {
		l.red = gpioreg.ByName(GPIO2_REV2)
	}
	l.loadingchan = nil

	return &l, nil
}

func (l *LEDImpl) Green(s bool) {
	l.green.Out(gpio.Level(s))
}

func (l *LEDImpl) Yellow(s bool) {
	l.yellow.Out(gpio.Level(s))
}

func (l *LEDImpl) Red(s bool) {
	l.red.Out(gpio.Level(s))
}

func (l *LEDImpl) StartupBlink() {
	for i := 0; i < 10; i++ {
		l.green.Out(gpio.High)
		l.yellow.Out(gpio.High)
		l.red.Out(gpio.High)
		time.Sleep(20 * time.Millisecond)
		l.green.Out(gpio.Low)
		l.yellow.Out(gpio.Low)
		l.red.Out(gpio.Low)
		time.Sleep(20 * time.Millisecond)
	}
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
				l.green.Out(gpio.High)
				l.yellow.Out(gpio.Low)
				l.red.Out(gpio.Low)
			case 1:
				l.green.Out(gpio.Low)
				l.yellow.Out(gpio.High)
				l.red.Out(gpio.Low)
			case 2:
				l.green.Out(gpio.Low)
				l.yellow.Out(gpio.Low)
				l.red.Out(gpio.High)
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
