package main

type LEDs interface {
	StartupBlink()
	StartLoading()
	StopLoading()
	Green(bool)
	Yellow(bool)
	Red(bool)
}
