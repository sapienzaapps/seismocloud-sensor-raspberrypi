package main

import (
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsclient"
	"git.sapienzaapps.it/seismocloud/seismocloud-client-go/scsutils"
	"os"
	"time"
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
			log.Error(err.Error())
			os.Exit(-1)
		} else {
			select {
			case serverTime := <-timeSync:
				if AbsI64(serverTime.Unix()-time.Now().Unix()) > 2 {
					log.Error("Local time is not synchronized")
					_ = scs.Close()
					_ = ledset.Green(false)
					_ = ledset.Yellow(true)
					_ = ledset.Red(true)
					time.Sleep(10 * time.Second)
					os.Exit(-1)
				}
				return
			default:
				// Timeout, retrying
				log.Warning("timeout syncing time, retrying in 10s")
				time.Sleep(10 * time.Second)
			}
		}
	}
}
