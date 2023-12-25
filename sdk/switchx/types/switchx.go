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
	ID     int     `json:"id" yaml:"id"`
	Config *Config `json:"config,omitempty" yaml:"config,omitempty"`
	On     *bool   `json:"on" yaml:"on"`
}

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
	WasOn           *bool  `json:"was_on,omitempty"`
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

// Status status of the Switch component contains information about the temperature, voltage, energy level and
// other physical characteristics of the switch instance. To obtain the status of the Switch component its id must be specified.
// For switches with power metering capabilities the status payload contains an additional set of properties with information
// about instantaneous power, supply voltage parameters and energy counters.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Switch#status
type Status struct {
	// ID Id of the Switch component instance
	ID *int `json:"id" yaml:"id"`
	// Source of the last command, for example: init, WS_in, http, ...
	Source *string `json:"source,omitempty" yaml:"source,omitempty"`
	// Output true if the output channel is currently on, false otherwise
	Output bool `json:"output,omitempty" yaml:"output,omitempty"`
	// TimerStartedAt Unix timestamp, start time of the timer (in UTC) (shown if the timer is triggered)
	TimerStartedAt *float64 `json:"timer_started_at,omitempty" yaml:"timer_started_at,omitempty"`
	// TimerDuration duration of the timer in seconds (shown if the timer is triggered)
	TimerDuration *float64 `json:"timer_duration,omitempty" yaml:"timer_duration,omitempty"`
	// Apower last measured instantaneous active power (in Watts) delivered to the attached load (shown if applicable)
	Apower *float64 `json:"apower,omitempty" yaml:"apower,omitempty"`
	// Voltage last measured voltage in Volts (shown if applicable)
	Voltage *float64 `json:"voltage,omitempty" yaml:"voltage,omitempty"`
	// Current last measured current in Amperes (shown if applicable)
	Current *float64 `json:"current,omitempty" yaml:"current,omitempty"`
	// PowerFactor last measured power factor (shown if applicable)
	PowerFactor *float64 `json:"pf,omitempty" yaml:"pf,omitempty"`
	// Aenergy information about the active energy counter (shown if applicable)
	Aenergy *SwitchAenergy `json:"aenergy,omitempty" yaml:"aenergy,omitempty"`
	// Temperature information about the temperature
	Temperature *SwitchTemperature `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	// Error conditions occurred. May contain overtemp, overpower, overvoltage, undervoltage, (shown if at least one error is present)
	Errors []string `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// SwitchAenergy information about the active energy counter (shown if applicable)
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Switch#status
type SwitchAenergy struct {
	// Total energy consumed in Watt-hours
	Total *float64 `json:"total,omitempty" yaml:"total,omitempty"`
	// ByMinute energy consumption by minute (in Milliwatt-hours) for the last three minutes
	// (the lower the index of the element in the array, the closer to the current moment the minute)
	ByMinute []float64 `json:"by_minute,omitempty" yaml:"by_minute,omitempty"`
	// MinuteTs Unix timestamp of the first second of the last minute (in UTC)
	MinuteTs *int `json:"minute_ts,omitempty" yaml:"minute_ts,omitempty"`
}

// Clone return copy
func (t *SwitchAenergy) Clone() *SwitchAenergy {
	c := &SwitchAenergy{}
	copier.Copy(&c, &t)
	return c
}

// SwitchTemperature System component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#status
type SwitchTemperature struct {
	// TC temperature in Celsius (null if temperature is out of the measurement range)
	TC *float64 `json:"tC,omitempty" yaml:"tC,omitempty"`
	// TF temperature in Fahrenheit (null if temperature is out of the measurement
	TF *float64 `json:"tF,omitempty" yaml:"tF,omitempty"`
}

// Clone return copy
func (t *SwitchTemperature) Clone() *SwitchTemperature {
	c := &SwitchTemperature{}
	copier.Copy(&c, &t)
	return c
}

// Config configuration of the Switch component contains information about the input mode, the timers and the protection
// settings of the chosen switch instance. To Get/Set the configuration of the Switch component its id must be specified.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Switch#configuration
type Config struct {
	// ID Id of the Switch component instance
	ID *int `json:"id" yaml:"id"`
	// Name of the switch instance
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// InMode range of values: momentary, follow, flip, detached
	InMode *string `json:"in_mode,omitempty" yaml:"in_mode,omitempty"`
	// InitialState range of values: off, on, restore_last, match_input
	InitialState *string `json:"initial_state,omitempty" yaml:"initial_state,omitempty"`
	// AutoOn True if the "Automatic ON" function is enabled, false otherwise
	AutoOn *bool `json:"auto_on,omitempty" yaml:"auto_on,omitempty"`
	// AutoOnDelay Seconds to pass until the component is switched back on
	AutoOnDelay *float64 `json:"auto_on_delay,omitempty" yaml:"auto_on_delay,omitempty"`
	// AutoOff True if the "Automatic OFF" function is enabled, false otherwise
	AutoOff *bool `json:"auto_off,omitempty" yaml:"auto_off,omitempty"`
	// AutoOffDelay Seconds to pass until the component is switched back off
	AutoOffDelay *float64 `json:"auto_off_delay,omitempty" yaml:"auto_off_delay,omitempty"`
	// AutorecoverVoltageErrors True if switch output state should be restored after over/undervoltage error is cleared, false otherwise (shown if applicable)
	AutorecoverVoltageErrors *bool `json:"autorecover_voltage_errors,omitempty" yaml:"autorecover_voltage_errors,omitempty"`
	// InputID Id of the Input component which controls the Switch. Applicable only to Pro1 and Pro1PM devices. Valid values: 0, 1
	InputID *int `json:"input_id,omitempty" yaml:"input_id,omitempty"`
	// PowerLimit Limit (in Watts) over which overpower condition occurs (shown if applicable)
	PowerLimit *float64 `json:"power_limit,omitempty" yaml:"power_limit,omitempty"`
	// VoltageLimit Limit (in Volts) over which overvoltage condition occurs (shown if applicable)
	VoltageLimit *float64 `json:"voltage_limit,omitempty" yaml:"voltage_limit,omitempty"`
	// UndervoltageLimit Limit (in Volts) under which undervoltage condition occurs (shown if applicable)
	UndervoltageLimit *float64 `json:"undervoltage_limit,omitempty" yaml:"undervoltage_limit,omitempty"`
	// CurrentLimit Number, limit (in Amperes) over which overcurrent condition occurs (shown if applicable)
	CurrentLimit *float64 `json:"current_limit,omitempty" yaml:"current_limit,omitempty"`
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

	if !util.CompareString(t.InMode, x.InMode) {
		zap.L().Info("Config InMode not equal")
		return false
	}

	if !util.CompareString(t.InitialState, x.InitialState) {
		zap.L().Info("Config InitialState not equal")
		return false
	}

	if !util.CompareBool(t.AutoOn, x.AutoOn) {
		zap.L().Info("Config AutoOn not equal")
		return false
	}

	if !util.CompareFloat64(t.AutoOnDelay, x.AutoOnDelay) {
		zap.L().Info("Config AutoOnDelay not equal")
		return false
	}

	if !util.CompareBool(t.AutoOff, x.AutoOff) {
		zap.L().Info("Config AutoOff not equal")
		return false
	}

	if !util.CompareFloat64(t.AutoOffDelay, x.AutoOffDelay) {
		zap.L().Info("Config AutoOffDelay not equal")
		return false
	}

	if !util.CompareBool(t.AutorecoverVoltageErrors, x.AutorecoverVoltageErrors) {
		zap.L().Info("Config AutorecoverVoltageErrors not equal")
		return false
	}

	if !util.CompareInt(t.InputID, x.InputID) {
		zap.L().Info("Config InputID not equal")
		return false
	}

	if !util.CompareFloat64(t.AutoOffDelay, x.AutoOffDelay) {
		zap.L().Info("Config AutoOffDelay not equal")
		return false
	}

	if !util.CompareFloat64(t.PowerLimit, x.PowerLimit) {
		zap.L().Info("Config PowerLimit not equal")
		return false
	}

	if !util.CompareFloat64(t.VoltageLimit, x.VoltageLimit) {
		zap.L().Info("Config VoltageLimit not equal")
		return false
	}

	if !util.CompareFloat64(t.UndervoltageLimit, x.UndervoltageLimit) {
		zap.L().Info("Config UndervoltageLimit not equal")
		return false
	}

	if !util.CompareFloat64(t.CurrentLimit, x.CurrentLimit) {
		zap.L().Info("Config CurrentLimit not equal")
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

	if t.InMode == nil {
		t.InMode = x.InMode
	}

	if t.InitialState == nil {
		t.InitialState = x.InitialState
	}

	if t.AutoOn == nil {
		t.AutoOn = x.AutoOn
	}

	if t.AutoOnDelay == nil {
		t.AutoOnDelay = x.AutoOnDelay
	}

	if t.AutoOff == nil {
		t.AutoOff = x.AutoOff
	}

	if t.AutoOffDelay == nil {
		t.AutoOffDelay = x.AutoOffDelay
	}

	if t.AutorecoverVoltageErrors == nil {
		t.AutorecoverVoltageErrors = x.AutorecoverVoltageErrors
	}

	if t.InputID == nil {
		t.InputID = x.InputID
	}

	if t.PowerLimit == nil {
		t.PowerLimit = x.PowerLimit
	}

	if t.VoltageLimit == nil {
		t.VoltageLimit = x.VoltageLimit
	}

	if t.UndervoltageLimit == nil {
		t.UndervoltageLimit = x.UndervoltageLimit
	}

	if t.CurrentLimit == nil {
		t.CurrentLimit = x.CurrentLimit
	}
}
