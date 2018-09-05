package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Config interface {
	GetDeviceId() string
	GetSigma() float64
	NewConfigReceived(sigma float32)
	RegisterSigmaCallback(SigmaCallback)
	RemoteReboot()
	UpdateCallback(hostname string, path string)
	Save() error
}

type SigmaCallback func(sigma float64)

func NewConfig() Config {
	cfgbuf, err := ioutil.ReadFile(CONFIG_PATH)
	ret := configImpl{
		sigmaCallback: nil,
	}
	if err == nil {
		err = json.Unmarshal(cfgbuf, &ret)
	}
	if err != nil || ret.DeviceId == "" {
		// Migration?
		migrfile, err := os.Open(CONFIG_PATH_OLD)
		if err == nil {
			s := bufio.NewScanner(migrfile)
			for s.Scan() {
				line := s.Text()
				if strings.HasPrefix(line, "deviceid:") {
					ret.DeviceId = strings.Replace(line, "deviceid:", "", -1)
					ret.Sigma = 6.0
					ret.Save()
				}
			}
			migrfile.Close()
			os.Remove(CONFIG_PATH_OLD)
		}
		if ret.DeviceId == "" {
			// No migration
			ret.Sigma = 6.0
			ret.DeviceId = strings.Replace(getMACAddress(), ":", "", -1)
			ret.Save()
		}
	}
	return &ret
}
