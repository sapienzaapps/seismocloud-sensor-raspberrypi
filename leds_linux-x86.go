// +build linuxx86

package main

import (
	"periph.io/x/periph/host"
)

func NewLEDs() (LEDs, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}
	l := LEDImpl{}
	return &l, nil
}

type LEDImpl struct {
}

func (l *LEDImpl) Green(s bool) {
}

func (l *LEDImpl) Yellow(s bool) {
}

func (l *LEDImpl) Red(s bool) {
}

func (l *LEDImpl) StartupBlink() {
}

func (l *LEDImpl) StartLoading() {
}

func (l *LEDImpl) StopLoading() {
}
