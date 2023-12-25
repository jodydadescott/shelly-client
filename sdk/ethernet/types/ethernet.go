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

// BasicResult internal use only
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

// EthernetStatus Ethernet component top level status
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Eth#status
type Status struct {
	// IP of the device in the network
	IP *string `json:"ip" yaml:"ip"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// Config Ethernet component top level config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Eth#configuration
type Config struct {
	// Enable True if the configuration is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
	// Ipv4Mode IPv4 mode. Range of values: dhcp, static
	Ipv4Mode *string `json:"ipv4mode,omitempty" yaml:"ipv4mode,omitempty"`
	// IP Ip to use when ipv4mode is static
	IP *string `json:"ip,omitempty" yaml:"ip,omitempty"`
	// Netmask to use when ipv4mode is static
	Netmask *string `json:"netmask,omitempty" yaml:"netmask,omitempty"`
	// Gateway to use when ipv4mode is static
	Gateway *string `json:"gw,omitempty" yaml:"gw,omitempty"`
	// Nameserver to use when ipv4mode is static
	Nameserver *string `json:"nameserver,omitempty" yaml:"nameserver,omitempty"`
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
		t.Ipv4Mode = nil
		t.IP = nil
		t.Netmask = nil
		t.Gateway = nil
		t.Nameserver = nil
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

	if !util.CompareString(t.Ipv4Mode, x.Ipv4Mode) {
		zap.L().Info("Config Ipv4Mode not equal")
		result = false
	}

	if !util.CompareString(t.IP, x.IP) {
		zap.L().Info("Config IP not equal")
		result = false
	}

	if !util.CompareString(t.Netmask, x.Netmask) {
		zap.L().Info("Config Netmask not equal")
		result = false
	}

	if !util.CompareString(t.Gateway, x.Gateway) {
		zap.L().Info("Config Gateway not equal")
		result = false
	}

	if !util.CompareString(t.Nameserver, x.Nameserver) {
		zap.L().Info("Config Nameserver not equal")
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

	if t.Ipv4Mode == nil {
		t.Ipv4Mode = x.Ipv4Mode
	}

	if t.IP == nil {
		t.IP = x.IP
	}

	if t.Netmask == nil {
		t.Netmask = x.Netmask
	}

	if t.Gateway == nil {
		t.Gateway = x.Gateway
	}

	if t.Nameserver == nil {
		t.Nameserver = x.Nameserver
	}

}
