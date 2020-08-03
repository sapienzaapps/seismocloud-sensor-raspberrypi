package main

import (
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
)

func onNewSigma(client scsclient.Client, newSigma float64) {
	log.Info("New sigma received: ", newSigma)
	cfg.SetSigma(newSigma)
	err := cfg.Save()
	if err != nil {
		log.Error("can't save configuration: ", err)
	}
}
