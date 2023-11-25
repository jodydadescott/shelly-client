package types

import (
	"github.com/jinzhu/copier"
)

// InputStatus status of the Input component contains information about the state of the chosen input instance.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Input#status
type InputStatus struct {
	// ID Id of the Input component instance
	ID *int `json:"id" yaml:"id"`
	// State (only for type switch, button) State of the input (null if the input instance is stateless, i.e. for type button)
	State *bool `json:"state" yaml:"state"`
	// Percent (only for type analog) Analog value in percent (null if valid value could not be obtained)
	Percent *int `json:"percent" yaml:"percent"`
	// Errors shown only if at least one error is present. May contain out_of_range, read
	Errors []string `json:"errors" yaml:"errors"`
}

// Clone return copy
func (t *InputStatus) Clone() *InputStatus {
	c := &InputStatus{}
	copier.Copy(&c, &t)
	return c
}

// InputConfig configuration of the Input component contains information about the type, invert and factory reset
// settings of the chosen input instance. To Get/Set the configuration of the Input component its id must be specified.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Input#configuration
type InputConfig struct {
	// ID of the Input component instance
	ID int `json:"id" yaml:"id"`
	// Name of the input instance
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// Type of associated input. Range of values switch, button, analog (only if applicable).
	Type *string `json:"type,omitempty" yaml:"type,omitempty"`
	// Invert (only for type switch, button) True if the logical state of the associated input is inverted,
	// false otherwise. For the change to be applied, the physical switch has to be toggled once after invert is set.
	Invert *bool `json:"invert" yaml:"invert"`
	// FactoryReset (only for type switch, button) True if input-triggered factory reset option is enabled,
	// false otherwise (shown if applicable)
	FactoryReset *bool `json:"factory_reset,omitempty" yaml:"factory_reset,omitempty"`
	// ReportThreshold (only for type analog) Analog input report threshold in percent.
	// Accepted range is device-specific, default [1.0..50.0]% unless specified otherwise
	ReportThreshold *float64 `json:"report_thr,omitempty" yaml:"report_thr,omitempty"`
}

// Clone return copy
func (t *InputConfig) Clone() *InputConfig {
	c := &InputConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *InputConfig) Markup() {

	if t == nil {
		return
	}

}

// Sanatize sanatize config
func (t *InputConfig) Sanatize() {

	if t == nil {
		return
	}

}
