// +build !rpi

package leds

import "fmt"

func (l *ledImpl) Init() error {
	return nil
}

type ledImpl struct {
}

func (l *ledImpl) Green(s bool) error {
	fmt.Printf("[G] %v\n", s)
	return nil
}

func (l *ledImpl) Yellow(s bool) error {
	fmt.Printf("[Y] %v\n", s)
	return nil
}

func (l *ledImpl) Red(s bool) error {
	fmt.Printf("[R] %v\n", s)
	return nil
}

func (l *ledImpl) StartupBlink() error {
	fmt.Println("LED Startup blink")
	return nil
}

func (l *ledImpl) StartLoading() {
	fmt.Println("LED start loading")
}

func (l *ledImpl) StopLoading() {
	fmt.Println("LED stop loading")
}
