package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gofrs/uuid"
)

// Interface is the interface for the configuration system
type Interface interface {
	GetDeviceID() uuid.UUID
	GetSigma() float64
	SetSigma(newSigma float64)
	Save() error
}

// New creates a new configuration instance
func New() (Interface, error) {
	ret := configImpl{}

	cfgbuf, err := ioutil.ReadFile(ConfigPath)
	if err == nil {
		err = json.Unmarshal(cfgbuf, &ret)
	} else if os.IsNotExist(err) {
		deviceId, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		ret.DeviceID = deviceId
		ret.Sigma = 6.0
		err = ret.Save()
	}

	return &ret, err
}

type configImpl struct {
	DeviceID uuid.UUID
	Sigma    float64
}

func (cfg *configImpl) GetSigma() float64 {
	return cfg.Sigma
}

func (cfg *configImpl) GetDeviceID() uuid.UUID {
	return cfg.DeviceID
}

func (cfg *configImpl) SetSigma(newSigma float64) {
	cfg.Sigma = newSigma
}

func (cfg *configImpl) Save() error {
	// Save configuration
	newcfg, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigPath, newcfg, 0600)
}

/*
func (cfg *configImpl) UpdateCallback(hostname string, path string) {
	err := updateStage1(fmt.Sprintf("http://%s/%s", hostname, path))
	if err != nil {
		log.Error(err.Error())
	}
}
*/
