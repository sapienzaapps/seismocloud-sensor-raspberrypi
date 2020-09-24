package leds

const (
	// StatusCalibration indicates that the calibration is in progress
	StatusCalibration int = iota

	// StatusWaitNetwork indicates that the connection to the network is in progress
	StatusWaitNetwork

	// StatusUpdate indicates that an update is in progress
	StatusUpdate

	// StatusConnecting indicates that the device is connecting
	StatusConnecting

	// StatusReady indicates that the device is ready
	StatusReady
)

// LEDSet represent an interface for leds
type LEDSet interface {
	Init() error

	StartupBlink() error
	StartLoading()
	StopLoading()
	Green(bool) error
	Yellow(bool) error
	Red(bool) error

	SetStatus(int) error
}

// New creates a new instance of LEDSet
func New() LEDSet {
	return &ledImpl{}
}

func (l *ledImpl) SetStatus(newStatus int) error {
	var err error
	switch newStatus {
	case StatusCalibration:
		err = l.Green(false)
		if err != nil {
			return err
		}
		err = l.Yellow(false)
		if err != nil {
			return err
		}
		err = l.Red(true)
		if err != nil {
			return err
		}
	case StatusWaitNetwork:
		err = l.Green(false)
		if err != nil {
			return err
		}
		err = l.Yellow(true)
		if err != nil {
			return err
		}
		err = l.Red(false)
		if err != nil {
			return err
		}
	case StatusUpdate:
		err = l.Green(false)
		if err != nil {
			return err
		}
		err = l.Yellow(true)
		if err != nil {
			return err
		}
		err = l.Red(true)
		if err != nil {
			return err
		}
	case StatusConnecting:
		err = l.Green(true)
		if err != nil {
			return err
		}
		err = l.Yellow(true)
		if err != nil {
			return err
		}
		err = l.Red(true)
		if err != nil {
			return err
		}
	case StatusReady:
		err = l.Green(true)
		if err != nil {
			return err
		}
		err = l.Yellow(false)
		if err != nil {
			return err
		}
		err = l.Red(false)
		if err != nil {
			return err
		}
	}
	return nil
}
