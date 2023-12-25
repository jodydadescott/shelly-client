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

// Status status of the Light component contains information about the brightness level and output state of the light instance.
// To obtain the status of the Light component its id must be specified.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Light#status
type Status struct {
	// ID Id of the Switch component instance
	ID *int `json:"id" yaml:"id"`
	// Source of the last command, for example: init, WS_in, http, ...
	Source *string `json:"source,omitempty" yaml:"source,omitempty"`
	// Output true if the output channel is currently on, false otherwise
	Output *bool `json:"output,omitempty" yaml:"output,omitempty"`
	// Brightness current brightness level (in percent)
	Brightness *float64 `json:"brightness" yaml:"brightness"`
	// TimerStartedAt Unix timestamp, start time of the timer (in UTC) (shown if the timer is triggered)
	TimerStartedAt *float64 `json:"timer_started_at,omitempty" yaml:"timer_started_at,omitempty"`
	// TimerDuration duration of the timer in seconds (shown if the timer is triggered)
	TimerDuration *float64 `json:"timer_duration,omitempty" yaml:"timer_duration,omitempty"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

type Config struct {
	// ID Id of the Switch component instance
	ID *int `json:"id" yaml:"id"`
	// Name of the switch instance
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
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
	// DefaultBrightness brightness level (in percent) after power on
	DefaultBrightness *float64 `json:"default.brightness,omitempty" yaml:"default.brightness,omitempty"`
	// NightModeEnable Enable or disable night mode
	NightModeEnable *bool `json:"night_mode.enable,omitempty" yaml:"night_mode.enable,omitempty"`
	// NightModeBrightness brightness level limit when night mode is active
	NightModeBrightness *float64 `json:"night_mode.brightness,omitempty" yaml:"night_mode.brightness,omitempty"`
	// NightModeActiveBetween containing 2 elements of type string, the first element indicates the start of
	// the period during which the night mode will be active, the second indicates the end of that period.
	// Both start and end are strings in the format HH:MM, where HH and MM are hours and minutes with optinal
	// leading zeros
	NightModeActiveBetween []string `json:"night_mode.active_between,omitempty" yaml:"night_mode.active_between,omitempty"`
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

	if !util.CompareFloat64(t.DefaultBrightness, x.DefaultBrightness) {
		zap.L().Info("Config DefaultBrightness not equal")
		return false
	}

	if !util.CompareBool(t.NightModeEnable, x.NightModeEnable) {
		zap.L().Info("Config NightModeEnable not equal")
		return false
	}

	if !util.CompareFloat64(t.NightModeBrightness, x.NightModeBrightness) {
		zap.L().Info("Config NightModeBrightness not equal")
		return false
	}

	if !util.CompareStringSlice(t.NightModeActiveBetween, x.NightModeActiveBetween) {
		zap.L().Info("Config NightModeActiveBetween not equal")
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

	if t.AutoOff == nil {
		t.AutoOff = x.AutoOff
	}

	if t.DefaultBrightness == nil {
		t.DefaultBrightness = x.DefaultBrightness
	}

	if t.NightModeEnable == nil {
		t.NightModeEnable = x.NightModeEnable
	}

	if t.NightModeBrightness == nil {
		t.NightModeBrightness = x.NightModeBrightness
	}

	t.NightModeActiveBetween = append(t.NightModeActiveBetween, x.NightModeActiveBetween...)
}
