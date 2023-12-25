package types

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/util"
)

type Request = msg_types.Request
type Response = msg_types.Response
type Error = msg_types.Error

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

// Params internal use only
type Params struct {
	Config *Config `json:"config,omitempty"`
}

// Status status of the Cloud component it can be checked whether the device is connected to the cloud.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cloud#status
type Status struct {
	// Connected true if the device is connected to the Shelly cloud, false otherwise
	Connected bool `json:"connected" yaml:"connected"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// Config configuration of the Cloud component shows information about the connection to the cloud
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cloud#configuration
type Config struct {
	// Enable true if cloud connection is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
	// Server name of the server to which the device is connected
	Server *string `json:"server,omitempty" yaml:"server,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *Config) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.Server = nil
	}

}

// Equals returns true if equal
func (t *Config) Equals(x *Config) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("Config receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("Config receiver is not nil but input is")
		return false
	}

	result := true

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("Config Enable not equal")
		result = false
	}

	if !util.CompareString(t.Server, x.Server) {
		zap.L().Info("Config Server not equal")
		result = false
	}

	return result
}

func (t *Config) Merge(x *Config) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

	if t.Server == nil {
		t.Server = x.Server
	}

}
