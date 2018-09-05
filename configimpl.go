package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type configImpl struct {
	DeviceId      string
	Sigma         float64
	sigmaCallback SigmaCallback `json:"-"`
}

func (cfg *configImpl) RegisterSigmaCallback(callback SigmaCallback) {
	cfg.sigmaCallback = callback
}

func (cfg *configImpl) GetSigma() float64 {
	return cfg.Sigma
}

func (cfg *configImpl) GetDeviceId() string {
	return cfg.DeviceId
}

func (cfg *configImpl) NewConfigReceived(sigma float32) {
	cfg.Sigma = float64(sigma)
	if cfg.sigmaCallback != nil {
		cfg.sigmaCallback(cfg.Sigma)
	}
}

func (cfg *configImpl) RemoteReboot() {
	log.Info("Reboot")
	exec.Command("reboot")
	os.Exit(0)
}

func (cfg *configImpl) UpdateCallback(hostname string, path string) {
	// TODO: implement self-update
}

func (cfg *configImpl) Save() error {
	// Save configuration
	newcfg, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CONFIG_PATH, newcfg, 0600)
}
