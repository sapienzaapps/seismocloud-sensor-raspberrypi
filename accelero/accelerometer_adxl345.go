// +build adxl345

package accelero

import (
	"encoding/binary"
	"os"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/utils"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// https://www.sparkfun.com/datasheets/Sensors/Accelerometer/ADXL345.pdf
// https://github.com/hybridgroup/gobot/blob/635adea96f788ab19bf77470b075f8db44d34fa0/drivers/i2c/adxl345_driver.go

const (
	i2cAddrHigh = 0x53
	i2cAddrLow  = 0x1D

	earthGravityMs2 = 9.80665
	scaleMultiplier = 0.004

	powerCltMeasure = 0x08
	powerCltStandBy = 0x00

	// Data rate
	adxl345Rate3200HZ = 0x0F // 3200 Hz
	adxl345Rate1600HZ = 0x0E // 1600 Hz
	adxl345Rate800HZ  = 0x0D // 800 Hz
	adxl345Rate400HZ  = 0x0C // 400 Hz
	adxl345Rate200HZ  = 0x0B // 200 Hz
	adxl345Rate100HZ  = 0x0A // 100 Hz
	adxl345Rate50HZ   = 0x09 // 50 Hz
	adxl345Rate25HZ   = 0x08 // 25 Hz
	adxl345Rate12_5HZ = 0x07 // 12.5 Hz
	adxl345Rate6_25HZ = 0x06 // 6.25 Hz
	adxl345Rate3_13HZ = 0x05 // 3.13 Hz
	adxl345Rate1_56HZ = 0x04 // 1.56 Hz
	adxl345Rate0_78HZ = 0x03 // 0.78 Hz
	adxl345Rate0_39HZ = 0x02 // 0.39 Hz
	adxl345Rate0_20HZ = 0x01 // 0.20 Hz
	adxl345Rate0_10HZ = 0x00 // 0.10 Hz

	// Data range
	adxl345Range2G  = 0x00 // +-2 g
	adxl345Range4G  = 0x01 // +-4 g
	adxl345Range8G  = 0x02 // +-8 g
	adxl345Range16G = 0x03 // +-16 g)

	adxl345RegDevID        = 0x00 // R,     11100101,   Device ID
	adxl345RegThreshTap    = 0x1D // R/W,   00000000,   Tap threshold
	adxl345RegOFSX         = 0x1E // R/W,   00000000,   X-axis offset
	adxl345RegOFSY         = 0x1F // R/W,   00000000,   Y-axis offset
	adxl345RegOFSZ         = 0x20 // R/W,   00000000,   Z-axis offset
	adxl345RegDUR          = 0x21 // R/W,   00000000,   Tap duration
	adxl345RegLATENT       = 0x22 // R/W,   00000000,   Tap latency
	adxl345RegWINDOW       = 0x23 // R/W,   00000000,   Tap window
	adxl345RegThreshACT    = 0x24 // R/W,   00000000,   Activity threshold
	adxl345RegThreshINACT  = 0x25 // R/W,   00000000,   Inactivity threshold
	adxl345RegTimeINACT    = 0x26 // R/W,   00000000,   Inactivity time
	adxl345RegAactInactCTL = 0x27 // R/W,   00000000,   Axis enable control for activity and inactiv ity detection
	adxl345RegThreshFF     = 0x28 // R/W,   00000000,   Free-fall threshold
	adxl345RegTimeFF       = 0x29 // R/W,   00000000,   Free-fall time
	adxl345RegTapAXES      = 0x2A // R/W,   00000000,   Axis control for single tap/double tap
	adxl345RegActTapSTATUS = 0x2B // R,     00000000,   Source of single tap/double tap
	adxl345RegBwRATE       = 0x2C // R/W,   00001010,   Data rate and power mode control
	adxl345RegPowerCTL     = 0x2D // R/W,   00000000,   Power-saving features control
	adxl345RegIntENABLE    = 0x2E // R/W,   00000000,   Interrupt enable control
	adxl345RegIntMAP       = 0x2F // R/W,   00000000,   Interrupt mapping control
	adxl345RegIntSOUCE     = 0x30 // R,     00000010,   Source of interrupts
	adxl345RegDataFORMAT   = 0x31 // R/W,   00000000,   Data format control
	adxl345RegDATAX0       = 0x32 // R,     00000000,   X-Axis Data 0
	adxl345RegDATAX1       = 0x33 // R,     00000000,   X-Axis Data 1
	adxl345RegDATAY0       = 0x34 // R,     00000000,   Y-Axis Data 0
	adxl345RegDATAY1       = 0x35 // R,     00000000,   Y-Axis Data 1
	adxl345RegDATAZ0       = 0x36 // R,     00000000,   Z-Axis Data 0
	adxl345RegDATAZ1       = 0x37 // R,     00000000,   Z-Axis Data 1
	adxl345RegFifoCTL      = 0x38 // R/W,   00000000,   FIFO control
	adxl345RegFifoSTATUS   = 0x39 // R,     00000000,   FIFO status
)

type adxl345Accelerometer struct {
	fd  *os.File
	bus i2c.BusCloser
	dev i2c.Dev

	normX float64
	normY float64
	normZ float64
}

// New creates a new instance of accelerometer for ADXL type sensors
func New() (Accelerometer, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	var err error
	ret := adxl345Accelerometer{}

	ret.bus, err = i2creg.Open("1")
	if err != nil {
		return nil, err
	}

	ret.dev = i2c.Dev{Bus: ret.bus, Addr: i2cAddrHigh}

	err = ret.setBandwidthRate(adxl345Rate100HZ)
	if err != nil {
		return nil, err
	}
	err = ret.setRange(adxl345Range2G)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (a *adxl345Accelerometer) setBandwidthRate(refreshRate uint8) error {
	return a.dev.Tx([]byte{adxl345RegBwRATE, refreshRate}, nil)
}

func (a *adxl345Accelerometer) setRange(rangeFlags uint8) error {
	value := [1]byte{}
	err := a.dev.Tx([]byte{adxl345RegDataFORMAT}, value[:])
	if err != nil {
		return err
	}

	value[0] &= ^uint8(0x0F)
	value[0] |= rangeFlags
	value[0] |= 0x08

	return a.dev.Tx([]byte{adxl345RegDataFORMAT, value[0]}, nil)
}

func (a *adxl345Accelerometer) Start() error {
	return a.dev.Tx([]byte{adxl345RegPowerCTL, powerCltMeasure}, nil)
}

func (a *adxl345Accelerometer) GetAccelerometerName() string {
	return "ADXL345"
}

func (a *adxl345Accelerometer) ProbeValueRaw() (AccelerometerData, error) {
	ret := AccelerometerData{0, 0, 0}

	buf := [6]byte{}

	err := a.dev.Tx([]byte{adxl345RegDATAX0}, buf[:])
	if err == nil {
		X := int16(binary.LittleEndian.Uint16(buf[0:2]))
		Y := int16(binary.LittleEndian.Uint16(buf[2:4]))
		Z := int16(binary.LittleEndian.Uint16(buf[4:6]))

		//X := bytes[0] | (bytes[1] << 8)
		//Y := bytes[2] | (bytes[3] << 8)
		//Z := bytes[4] | (bytes[5] << 8)
		//
		ret.X = float64(X) * scaleMultiplier
		ret.Y = float64(Y) * scaleMultiplier
		ret.Z = float64(Z) * scaleMultiplier

		//if (!gforce) {
		//	ret.X = ret.X * EARTH_GRAVITY_MS2
		//	ret.Y = ret.Y * EARTH_GRAVITY_MS2
		//	ret.Z = ret.Z * EARTH_GRAVITY_MS2
		//}
	}

	return ret, nil
}

func (a *adxl345Accelerometer) ProbeValue() (AccelerometerData, error) {
	ret, err := a.ProbeValueRaw()
	if err != nil {
		return ret, err
	}

	ret.X = ret.X - a.normX
	ret.Y = ret.Y - a.normY
	ret.Z = ret.Z - a.normZ
	return ret, nil
}

func (a *adxl345Accelerometer) Calibration() {
	t := time.NewTicker(50 * time.Millisecond)

	// Calibration
	avgX := utils.NewRunningAvgFloat64()
	avgY := utils.NewRunningAvgFloat64()
	avgZ := utils.NewRunningAvgFloat64()
	startms := time.Now()
	for time.Since(startms) < 10*time.Second {
		probe, _ := a.ProbeValue()
		avgX.AddValue(probe.X)
		avgY.AddValue(probe.Y)
		avgZ.AddValue(probe.Z)
		<-t.C
	}

	a.normX = avgX.GetAverage()
	a.normY = avgY.GetAverage()
	a.normZ = avgZ.GetAverage()
}

func (a *adxl345Accelerometer) Stop() error {
	return a.dev.Tx([]byte{adxl345RegPowerCTL, powerCltStandBy}, nil)
}

func (a *adxl345Accelerometer) Shutdown() error {
	return a.bus.Close()
}
