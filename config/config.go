package config

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
)

type Config interface {
	GetDeviceId() uuid.UUID
	GetSigma() float64
	SetSigma(newSigma float64)
	Save() error
}

type SigmaCallback func(sigma float64)

func New() (Config, error) {
	ret := configImpl{}

	cfgbuf, err := ioutil.ReadFile(CONFIG_PATH)
	if err == nil {
		err = json.Unmarshal(cfgbuf, &ret)
	} else if err == os.ErrNotExist {
		ret.DeviceId = uuid.NewV4()
		ret.Sigma = 6.0
		err = ret.Save()
	}

	return &ret, err
}
