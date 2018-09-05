package main

import (
	"bytes"
	"net"
)

func minint(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func getMACAddress() string {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				return i.HardwareAddr.String()
			}
		}
	}
	return ""
}
