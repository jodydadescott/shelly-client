package light

import (
	"github.com/jodydadescott/shelly-client/sdk/types"
)

type Request = types.Request
type Response = types.Response
type Error = types.Error
type MessageHandlerFactory = types.MessageHandlerFactory
type MessageHandler = types.MessageHandler

type ShelllyDeviceInfo = types.ShelllyDeviceInfo

type Status = types.LightStatus
type Config = types.LightConfig

// Params internal use only
type Params struct {
	ID         int      `json:"id" yaml:"id"`
	Config     *Config  `json:"config,omitempty" yaml:"on,omitempty"`
	On         *bool    `json:"on,omitempty" yaml:"on,omitempty"`
	Brightness *float64 `json:"brightness,omitempty" yaml:"brightness,omitempty"`
}

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

// GetConfigResponse internal use only
type GetConfigResponse struct {
	Response
	Result *Config `json:"result,omitempty"`
	Params *Params `json:"params,omitempty"`
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
