// +build !rpi

package leds

func (l *LEDImpl) Init() error {
	return nil
}

type LEDImpl struct {
}

func (l *LEDImpl) Green(s bool) error {
	return nil
}

func (l *LEDImpl) Yellow(s bool) error {
	return nil
}

func (l *LEDImpl) Red(s bool) error {
	return nil
}

func (l *LEDImpl) StartupBlink() error {
	return nil
}

func (l *LEDImpl) StartLoading() {
}

func (l *LEDImpl) StopLoading() {
}
