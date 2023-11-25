package system

import (
	"github.com/jodydadescott/shelly-client/sdk/types"
)

type Request = types.Request
type Response = types.Response
type Error = types.Error
type MessageHandlerFactory = types.MessageHandlerFactory
type MessageHandler = types.MessageHandler

type Config = types.SystemConfig
type Status = types.SystemStatus
type DeviceConfig = types.SystemDevice
type LocationConfig = types.SystemLocation
type DebugConfig = types.SystemDebug
type UIDataConfig = types.SystemUIData
type RPCUDPConfig = types.SystemRPCUDP
type SntpConfig = types.SystemSntp
type MqttDebug = types.SystemMqtt
type WebsocketDebug = types.SystemWebsocket
type UDP = types.SystemUDP

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

// GetConfigResponse internal use only
type GetConfigResponse struct {
	Response
	Result *Config `json:"result,omitempty"`
}

// SetConfigResponse internal use only
type SetConfigResponse struct {
	Response
	Result *Result `json:"result,omitempty"`
}

// GetStatusResponse internal use only
type GetStatusResponse struct {
	Response
	Result *Status `json:"result,omitempty"`
}

// Params internal use only
type Params struct {
	Config *Config `json:"config,omitempty"`
}
