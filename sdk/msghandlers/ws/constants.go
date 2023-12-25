package ws

import (
	"time"
)

const (
	WsScheme           = "ws"
	wsPath             = "/rpc"
	defaultRetryWait   = time.Duration(3) * time.Second
	defaultSendTimeout = time.Duration(time.Second * 10)
	defaultSendTrys    = 3
	defaultShellyUser  = "admin"
)
