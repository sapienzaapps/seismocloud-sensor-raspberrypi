package main

import (
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"os"
	"os/exec"
)

func onReboot(client scsclient.Client) {
	log.Info("Reboot")
	err := exec.Command("reboot").Start()
	if err != nil {
		log.Error("error trying to reboot: ", err)
	} else {
		os.Exit(0)
	}
}
