package main

import (
	"bytes"
	"net"
)

const (
	PKTTYPE_DISCOVERY       = 1
	PKTTYPE_DISCOVERY_REPLY = 2
	PKYTYPE_PING            = 3
	PKYTYPE_PONG            = 4
	PKTTYPE_START           = 5
	PKTTYPE_STOP            = 6
	PKTTYPE_SENDGPS         = 7
	PKTTYPE_OK              = 8
	PKTTYPE_SETSYSLOG       = 9
	PKTTYPE_REBOOT          = 10
	PKTTYPE_GETINFO         = 11
	PKTTYPE_GETINFO_REPLY   = 12
	PKTTYPE_RESET           = 13
	PKTTYPE_TRACE           = 14
)

var lanInterfaceStop chan interface{} = nil

func lanInterfaceWorker(deviceId string) {
	addr := net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Error(err.Error())
		lanInterfaceStop = nil
		return
	}

	go func() {
		<-lanInterfaceStop
		conn.Close()
	}()

	var buf [1024]byte

	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			break
		}

		if rlen > 5 && bytes.Compare(buf[:5], []byte("INGV\000")) == 0 {
			switch buf[5] {
			case PKTTYPE_DISCOVERY:
				// TODO: reply
				newbuf := make([]byte, 5+1+6+4+8)
				copy(newbuf, "INGV\000")
				newbuf[5] = PKTTYPE_DISCOVERY_REPLY
				copy(newbuf[6:], deviceId)
				copy(newbuf[6+6:], VERSION)
				copy(newbuf[6+6+4:], MODEL[:minint(8, len(MODEL))])

				conn.WriteToUDP(newbuf, remote)
			}
		}
	}
}

func StartLANInterface(deviceid string) {
	if lanInterfaceStop == nil {
		lanInterfaceStop = make(chan interface{}, 1)
		go lanInterfaceWorker(deviceid)
	}
}

func StopLANInterface() {
	lanInterfaceStop <- 1
}
