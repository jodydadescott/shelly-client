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

// Params internal use only
type Params struct {
	Config *Config `json:"config,omitempty" yaml:"config,omitempty"`
	ID     int     `json:"id" yaml:"id"`
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

// Status status of the Input component contains information about the state of the chosen input instance.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Input#status
type Status struct {
	// ID Id of the Input component instance
	ID *int `json:"id" yaml:"id"`
	// State (only for type switch, button) State of the input (null if the input instance is stateless, i.e. for type button)
	State *bool `json:"state,omitempty" yaml:"state,omitempty"`
	// Percent (only for type analog) Analog value in percent (null if valid value could not be obtained)
	Percent *int `json:"percent" yaml:"percent"`
	// Errors shown only if at least one error is present. May contain out_of_range, read
	Errors []string `json:"errors" yaml:"errors"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// Config configuration of the Input component contains information about the type, invert and factory reset
// settings of the chosen input instance. To Get/Set the configuration of the Input component its id must be specified.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Input#configuration
type Config struct {
	// ID of the Input component instance
	ID *int `json:"id" yaml:"id"`
	// Name of the input instance
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// Type of associated input. Range of values switch, button, analog (only if applicable).
	Type *string `json:"type,omitempty" yaml:"type,omitempty"`
	// Invert (only for type switch, button) True if the logical state of the associated input is inverted,
	// false otherwise. For the change to be applied, the physical switch has to be toggled once after invert is set.
	Invert *bool `json:"invert,omitempty" yaml:"invert,omitempty"`
	// FactoryReset (only for type switch, button) True if input-triggered factory reset option is enabled,
	// false otherwise (shown if applicable)
	FactoryReset *bool `json:"factory_reset,omitempty" yaml:"factory_reset,omitempty"`
	// ReportThreshold (only for type analog) Analog input report threshold in percent.
	// Accepted range is device-specific, default [1.0..50.0]% unless specified otherwise
	ReportThreshold *float64 `json:"report_thr,omitempty" yaml:"report_thr,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
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

	if !util.CompareInt(t.ID, x.ID) {
		zap.L().Info("Config ID not equal")
		return false
	}

	if !util.CompareString(t.Name, x.Name) {
		zap.L().Info("Config Name not equal")
		return false
	}

	if !util.CompareString(t.Type, x.Type) {
		zap.L().Info("Config Type not equal")
		return false
	}

	if !util.CompareBool(t.Invert, x.Invert) {
		zap.L().Info("Config Invert not equal")
		return false
	}

	if !util.CompareFloat64(t.ReportThreshold, x.ReportThreshold) {
		zap.L().Info("Config ReportThreshold not equal")
		return false
	}

	return true
}

func (t *Config) Merge(x *Config) {

	if x == nil {
		return
	}

	if t.ID == nil {
		t.ID = x.ID
	}

	if t.Name == nil {
		t.Name = x.Name
	}

	if t.Type == nil {
		t.Type = x.Type
	}

	if t.Invert == nil {
		t.Invert = x.Invert
	}

	if t.FactoryReset == nil {
		t.FactoryReset = x.FactoryReset
	}

	if t.ReportThreshold == nil {
		t.ReportThreshold = x.ReportThreshold
	}

}
