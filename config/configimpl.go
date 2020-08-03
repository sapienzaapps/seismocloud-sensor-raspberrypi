package config

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
)

type configImpl struct {
	DeviceId uuid.UUID
	Sigma    float64
}

func (cfg *configImpl) GetSigma() float64 {
	return cfg.Sigma
}

func (cfg *configImpl) GetDeviceId() uuid.UUID {
	return cfg.DeviceId
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
	return ioutil.WriteFile(CONFIG_PATH, newcfg, 0600)
}

/*
func (cfg *configImpl) UpdateCallback(hostname string, path string) {
	err := updateStage1(fmt.Sprintf("http://%s/%s", hostname, path))
	if err != nil {
		log.Error(err.Error())
	}
}
*/
