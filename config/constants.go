//+build prod

package config

const (
	// Version is the software/firmware version
	Version = "3.00"

	// MqttServer is the server for signaling
	MqttServer = "tls://mqtts.seismocloud.com:443"
)
