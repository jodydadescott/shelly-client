package mqtt

import "time"

const (
	defaultKeepAlive      = time.Second * 2
	defaultPingTimeout    = time.Second * 2
	defaultPublishTimeout = time.Second * 10
	defaultDaemonInterval = time.Minute * 10

	defaultClientID    = "shelly-client-mqtt"
	defaultLightSource = defaultClientID

	shellyPlusWallDimmer = "PlusWallDimmer"
)
