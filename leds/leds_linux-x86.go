// +build !rpi

package leds

import "fmt"

func (l *ledImpl) Init() error {
	return nil
}

type ledImpl struct {
	red    bool
	yellow bool
	green  bool
}

func (l *ledImpl) Green(s bool) error {
	if l.green != s {
		l.green = s
		fmt.Printf("[G] %v\n", s)
	}
	return nil
}

func (l *ledImpl) Yellow(s bool) error {
	if l.yellow != s {
		l.yellow = s
		fmt.Printf("[Y] %v\n", s)
	}
	return nil
}

func (l *ledImpl) Red(s bool) error {
	if l.red != s {
		l.red = s
		fmt.Printf("[R] %v\n", s)
	}
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
