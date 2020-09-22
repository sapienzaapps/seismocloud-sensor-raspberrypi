package main

import (
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
)

// TODO: close all (seismometer, MQTT connection, etc) gracefully
func onReboot(client scsclient.Client) {
	log.Info("Reboot")
	err := reboot()
	if err != nil {
		log.Error("error trying to reboot: ", err)
	}
}
