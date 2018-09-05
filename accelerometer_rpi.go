package main

import (
	"encoding/binary"
	"os"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// https://www.sparkfun.com/datasheets/Sensors/Accelerometer/ADXL345.pdf

const (
	I2C_ADDR_HIGH = 0x53
	I2C_ADDR_LOW  = 0x1D

	EARTH_GRAVITY_MS2 = 9.80665
	SCALE_MULTIPLIER  = 0.004

	POWER_CLT_MEASURE = 0x08
	POWER_CLT_STANDBY = 0x00

	// Data rate
	ADXL345_RATE_3200HZ = 0x0F // 3200 Hz
	ADXL345_RATE_1600HZ = 0x0E // 1600 Hz
	ADXL345_RATE_800HZ  = 0x0D // 800 Hz
	ADXL345_RATE_400HZ  = 0x0C // 400 Hz
	ADXL345_RATE_200HZ  = 0x0B // 200 Hz
	ADXL345_RATE_100HZ  = 0x0A // 100 Hz
	ADXL345_RATE_50HZ   = 0x09 // 50 Hz
	ADXL345_RATE_25HZ   = 0x08 // 25 Hz
	ADXL345_RATE_12_5HZ = 0x07 // 12.5 Hz
	ADXL345_RATE_6_25HZ = 0x06 // 6.25 Hz
	ADXL345_RATE_3_13HZ = 0x05 // 3.13 Hz
	ADXL345_RATE_1_56HZ = 0x04 // 1.56 Hz
	ADXL345_RATE_0_78HZ = 0x03 // 0.78 Hz
	ADXL345_RATE_0_39HZ = 0x02 // 0.39 Hz
	ADXL345_RATE_0_20HZ = 0x01 // 0.20 Hz
	ADXL345_RATE_0_10HZ = 0x00 // 0.10 Hz

	// Data range
	ADXL345_RANGE_2G  = 0x00 // +-2 g
	ADXL345_RANGE_4G  = 0x01 // +-4 g
	ADXL345_RANGE_8G  = 0x02 // +-8 g
	ADXL345_RANGE_16G = 0x03 // +-16 g)

	ADXL345_REG_DEVID          = 0x00 // R,     11100101,   Device ID
	ADXL345_REG_THRESH_TAP     = 0x1D // R/W,   00000000,   Tap threshold
	ADXL345_REG_OFSX           = 0x1E // R/W,   00000000,   X-axis offset
	ADXL345_REG_OFSY           = 0x1F // R/W,   00000000,   Y-axis offset
	ADXL345_REG_OFSZ           = 0x20 // R/W,   00000000,   Z-axis offset
	ADXL345_REG_DUR            = 0x21 // R/W,   00000000,   Tap duration
	ADXL345_REG_LATENT         = 0x22 // R/W,   00000000,   Tap latency
	ADXL345_REG_WINDOW         = 0x23 // R/W,   00000000,   Tap window
	ADXL345_REG_THRESH_ACT     = 0x24 // R/W,   00000000,   Activity threshold
	ADXL345_REG_THRESH_INACT   = 0x25 // R/W,   00000000,   Inactivity threshold
	ADXL345_REG_TIME_INACT     = 0x26 // R/W,   00000000,   Inactivity time
	ADXL345_REG_ACT_INACT_CTL  = 0x27 // R/W,   00000000,   Axis enable control for activity and inactiv ity detection
	ADXL345_REG_THRESH_FF      = 0x28 // R/W,   00000000,   Free-fall threshold
	ADXL345_REG_TIME_FF        = 0x29 // R/W,   00000000,   Free-fall time
	ADXL345_REG_TAP_AXES       = 0x2A // R/W,   00000000,   Axis control for single tap/double tap
	ADXL345_REG_ACT_TAP_STATUS = 0x2B // R,     00000000,   Source of single tap/double tap
	ADXL345_REG_BW_RATE        = 0x2C // R/W,   00001010,   Data rate and power mode control
	ADXL345_REG_POWER_CTL      = 0x2D // R/W,   00000000,   Power-saving features control
	ADXL345_REG_INT_ENABLE     = 0x2E // R/W,   00000000,   Interrupt enable control
	ADXL345_REG_INT_MAP        = 0x2F // R/W,   00000000,   Interrupt mapping control
	ADXL345_REG_INT_SOUCE      = 0x30 // R,     00000010,   Source of interrupts
	ADXL345_REG_DATA_FORMAT    = 0x31 // R/W,   00000000,   Data format control
	ADXL345_REG_DATAX0         = 0x32 // R,     00000000,   X-Axis Data 0
	ADXL345_REG_DATAX1         = 0x33 // R,     00000000,   X-Axis Data 1
	ADXL345_REG_DATAY0         = 0x34 // R,     00000000,   Y-Axis Data 0
	ADXL345_REG_DATAY1         = 0x35 // R,     00000000,   Y-Axis Data 1
	ADXL345_REG_DATAZ0         = 0x36 // R,     00000000,   Z-Axis Data 0
	ADXL345_REG_DATAZ1         = 0x37 // R,     00000000,   Z-Axis Data 1
	ADXL345_REG_FIFO_CTL       = 0x38 // R/W,   00000000,   FIFO control
	ADXL345_REG_FIFO_STATUS    = 0x39 // R,     00000000,   FIFO status
)

type ADXL345Accelerometer struct {
	fd  *os.File
	bus i2c.BusCloser
	dev i2c.Dev
}

func CreateNewADXL345Accelerometer() (Accelerometer, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	var err error
	ret := ADXL345Accelerometer{}

	ret.bus, err = i2creg.Open("1")
	if err != nil {
		return nil, err
	}

	ret.dev = i2c.Dev{Bus: ret.bus, Addr: I2C_ADDR_HIGH}

	err = ret.setBandwidthRate(ADXL345_RATE_100HZ)
	if err != nil {
		return nil, err
	}
	err = ret.setRange(ADXL345_RANGE_2G)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (a *ADXL345Accelerometer) setBandwidthRate(refreshRate uint8) error {
	return a.dev.Tx([]byte{ADXL345_REG_BW_RATE, refreshRate}, nil)
}

func (a *ADXL345Accelerometer) setRange(rangeFlags uint8) error {
	value := [1]byte{}
	err := a.dev.Tx([]byte{ADXL345_REG_DATA_FORMAT}, value[:])
	if err != nil {
		return err
	}

	value[0] &= ^uint8(0x0F)
	value[0] |= rangeFlags
	value[0] |= 0x08

	return a.dev.Tx([]byte{ADXL345_REG_DATA_FORMAT, value[0]}, nil)
}

func (a *ADXL345Accelerometer) Start() error {
	return a.dev.Tx([]byte{ADXL345_REG_POWER_CTL, POWER_CLT_MEASURE}, nil)
}

func (a *ADXL345Accelerometer) GetAccelerometerName() string {
	return "ADXL345"
}

func (a *ADXL345Accelerometer) ProbeValue() (AccelerometerData, error) {
	ret := AccelerometerData{0, 0, 0}

	buf := [6]byte{}

	err := a.dev.Tx([]byte{ADXL345_REG_DATAX0}, buf[:])
	if err == nil {
		X := int16(binary.LittleEndian.Uint16(buf[0:2]))
		Y := int16(binary.LittleEndian.Uint16(buf[2:4]))
		Z := int16(binary.LittleEndian.Uint16(buf[4:6]))

		//X := bytes[0] | (bytes[1] << 8)
		//Y := bytes[2] | (bytes[3] << 8)
		//Z := bytes[4] | (bytes[5] << 8)
		//
		ret.X = float64(X) * SCALE_MULTIPLIER
		ret.Y = float64(Y) * SCALE_MULTIPLIER
		ret.Z = float64(Z) * SCALE_MULTIPLIER

		//if (!gforce) {
		//	ret.X = ret.X * EARTH_GRAVITY_MS2
		//	ret.Y = ret.Y * EARTH_GRAVITY_MS2
		//	ret.Z = ret.Z * EARTH_GRAVITY_MS2
		//}
	}
	return ret, nil
}

func (a *ADXL345Accelerometer) Stop() error {
	return a.dev.Tx([]byte{ADXL345_REG_POWER_CTL, POWER_CLT_STANDBY}, nil)
}

func (a *ADXL345Accelerometer) Shutdown() error {
	return a.bus.Close()
}
