package main

import (
	"os"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/utils"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsutils"
)

var timeSync = make(chan time.Time, 1)

func onTimeReceived(_ scsclient.Client, t0 int64, t1 int64, t2 int64, t3 int64) {
	// TODO: drain timeSync
	timeSync <- scsutils.SyncSNTPTime(t0, t1, t2, t3)
}

func checkTime() {
	for {
		err := scs.RequestTime()
		if err != nil {
			// TODO: better logging and error handling
			log.Error(err.Error())
			os.Exit(-1)
		} else {
			select {
			case serverTime := <-timeSync:
				if utils.AbsI64(serverTime.Unix()-time.Now().Unix()) > 2 {
					log.Error("Local time is not synchronized")
					_ = scs.Close()
					_ = ledset.Green(false)
					_ = ledset.Yellow(true)
					_ = ledset.Red(true)
					time.Sleep(10 * time.Second)
					// TODO: better logging and error handling
					os.Exit(-1)
				}
				return
			case <-time.After(10 * time.Second):
				// TODO: do not block for signals
				// Timeout, retrying
				log.Warning("timeout syncing time, retrying in 10s")
			}
		}
	}
}
