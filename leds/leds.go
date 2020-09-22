package leds

// LEDSet represent an interface for leds
type LEDSet interface {
	Init() error

	StartupBlink() error
	StartLoading()
	StopLoading()
	Green(bool) error
	Yellow(bool) error
	Red(bool) error
}

// New creates a new instance of LEDSet
func New() LEDSet {
	return &ledImpl{}
}
