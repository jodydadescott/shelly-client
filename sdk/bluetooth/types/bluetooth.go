package types

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/util"
)

// type Request = types.Request
type Response = msg_types.Response
type Error = msg_types.Error

// type MessageHandlerFactory = types.MessageHandlerFactory
// type MessageHandler = types.MessageHandler

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

// Status status of the BLE component contains information about the bluetooth on/off state and
// does not own any status properties.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#status
type Status struct {
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// Config configuration of the Bluetooth Low Energy component shows whether the bluetooth connection is enabled.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type Config struct {
	// Enable True if bluetooth is enabled, false otherwise
	Enable *bool `json:"enable" yaml:"enable"`
	// RPC configuration of the rpc service
	RPC *RPC `json:"rpc,omitempty" yaml:"rpc,omitempty"`
	// Observer configuration of the BT LE observer
	Observer *Observer `json:"observer,omitempty" yaml:"observer,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *Config) Sanatize() {

	if t.Enable == nil || !*t.Enable {
		t.RPC = nil
		t.Observer = nil
		return
	}

	t.RPC.Sanatize()
	t.Observer.Sanatize()
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

	if !t.RPC.Equals(x.RPC) {
		zap.L().Info("Config RPC not equal")
		result = false
	}

	if !t.Observer.Equals(x.Observer) {
		zap.L().Info("Config Observer not equal")
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

	if t.RPC == nil {
		if x.RPC != nil {
			t.RPC = x.RPC.Clone()
		} else {
			t.RPC.Merge(x.RPC)
		}
	}

	if t.Observer == nil {
		if x.Observer != nil {
			t.Observer = x.Observer.Clone()
		} else {
			t.Observer.Merge(x.Observer)
		}
	}
}

// RPC configuration of the rpc service
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type RPC struct {
	// Enable True if rpc service is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// Clone return copy
func (t *RPC) Clone() *RPC {
	c := &RPC{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *RPC) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}
}

// Equals returns true if equal
func (t *RPC) Equals(x *RPC) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("RPC receiver is nil but input is not")
		return x == nil
	}

	if x == nil {
		zap.L().Info("RPC receiver is not nil but input is")
		return false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("RPC enabled not equal")
		return false
	}

	return true
}

func (t *RPC) Merge(x *RPC) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}
}

// Observer configuration of the BT LE observer
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type Observer struct {
	// Enable true if BT LE observer is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// Clone return copy
func (t *Observer) Clone() *Observer {
	c := &Observer{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *Observer) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}
}

// Equals returns true if equal
func (t *Observer) Equals(x *Observer) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("Observer receiver is nil but input is not")
		return x == nil
	}

	if x == nil {
		zap.L().Info("Observer receiver is not nil but input is")
		return false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("Observer Enable not equal")
		return false
	}

	return true
}

func (t *Observer) Merge(x *Observer) {
	if x == nil {
		return
	}
	if t.Enable == nil {
		t.Enable = x.Enable
	}
}
