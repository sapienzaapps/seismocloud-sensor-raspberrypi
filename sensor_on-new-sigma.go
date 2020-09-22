package main

import (
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
)

func onNewSigmaFunc(seismometer Seismometer) func(scsclient.Client, float64) {
	return func(client scsclient.Client, newSigma float64) {
		log.Info("New sigma received: ", newSigma)
		seismometer.SetSigma(newSigma)
		cfg.SetSigma(newSigma)
		err := cfg.Save()
		if err != nil {
			log.Error("can't save configuration: ", err)
		}
	}
}
