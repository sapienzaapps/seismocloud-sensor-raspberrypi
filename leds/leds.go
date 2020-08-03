package leds

type LEDSet interface {
	Init() error

	StartupBlink() error
	StartLoading()
	StopLoading()
	Green(bool) error
	Yellow(bool) error
	Red(bool) error
}

func New() LEDSet {
	return &LEDImpl{}
}
