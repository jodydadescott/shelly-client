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

// Status status contains information about network state, system time and other common attributes of the Shelly device.
// Presence of some keys is optional, depending on the underlying hardware components.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#Status
type Status struct {
	// MAC address of the device
	MAC *string `json:"enable,omitempty" yaml:"enable,omitempty"`
	// RestartRequired true if restart is required, false otherwise
	RestartRequired *bool `json:"restart_required,omitempty" yaml:"restart_required,omitempty"`
	// Time Current time in the format HH:MM (24-hour time format in the current timezone with leading zero).
	// null when time is not synced from NTP server.
	Time *string `json:"time,omitempty" yaml:"time,omitempty"`
	// Unixtime Unix timestamp (in UTC), null when time is not synced from NTP server.
	Unixtime *float64 `json:"unixtime,omitempty" yaml:"unixtime,omitempty"`
	// Uptime Time in seconds since last reboot
	Uptime *float64 `json:"uptime,omitempty" yaml:"uptime,omitempty"`
	// RAMSize Total size of the RAM in the system in Bytes
	RAMSize *float64 `json:"ram_size,omitempty" yaml:"ram_size,omitempty"`
	// RAMFree Size of the free RAM in the system in Bytes
	RAMFree *float64 `json:"ram_free,omitempty" yaml:"ram_free,omitempty"`
	// FsSize Total size of the file system in Bytes
	FsSize *float64 `json:"fs_size,omitempty" yaml:"fs_size,omitempty"`
	// FsFree Size of the free file system in Bytes
	FsFree *float64 `json:"fs_free,omitempty" yaml:"fs_free,omitempty"`
	// CfgRev Configuration revision number
	CfgRev *float64 `json:"cfg_rev,omitempty" yaml:"cfg_rev,omitempty"`
	// KvsRev KVS (Key-Value Store) revision number
	KvsRev *float64 `json:"kvs_rev,omitempty" yaml:"kvs_rev,omitempty"`
	// ScheduleRev Schedules revision number, present if schedules are enabled
	ScheduleRev *float64 `json:"schedule_rev,omitempty" yaml:"schedule_rev,omitempty"`
	// WebhookRev Webhooks revision number, present if webhooks are enabled
	WebhookRev *float64 `json:"webhook_rev,omitempty" yaml:"webhook_rev,omitempty"`
	// AvailableUpdates Information about available updates, similar to the one returned by Shelly.CheckForUpdate
	// (empty object: {}, if no updates available). This information is automatically updated every 24 hours.
	// Note that build_id and url for an update are not displayed here
	AvailableUpdates *SystemAvailableUpdates `json:"available_updates,omitempty" yaml:"available_updates,omitempty"`
	// WakeupReason Information about boot type and cause (only for battery-operated devices)
	WakeupReason *SystemWakeupReason `json:"wakeup_reason,omitempty" yaml:"wakeup_reason,omitempty"`
	// WakeupPeriod Period (in seconds) at which device wakes up and sends "keep-alive" packet to cloud, readonly.
	// Count starts from last full wakeup
	WakeupPeriod *int `json:"wakeup_period,omitempty" yaml:"wakeup_period,omitempty"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// SystemAvailableUpdates Information about available updates, similar to the one returned by Shelly.CheckForUpdate
// (empty object: {}, if no updates available). This information is automatically updated every 24 hours.
// Note that build_id and url for an update are not displayed here
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys/#status
type SystemAvailableUpdates struct {
	// Beta shown only if beta update is available
	Beta *FirmwareStatus `json:"beta,omitempty" yaml:"beta,omitempty"`
	// Stable version of the new firmware. Shown only if stable update is available
	Stable *FirmwareStatus `json:"stable,omitempty" yaml:"stable,omitempty"`
}

// Clone return copy
func (t *SystemAvailableUpdates) Clone() *SystemAvailableUpdates {
	c := &SystemAvailableUpdates{}
	copier.Copy(&c, &t)
	return c
}

// SystemWakeupReason information about boot type and cause (only for battery-operated devices)
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys
type SystemWakeupReason struct {
	// Boot type, one of: poweron, software_restart, deepsleep_wake, internal (e.g. brownout detection, watchdog timeout, etc.), unknown
	Boot *string `json:"boot,omitempty" yaml:"boot,omitempty"`
	// Cause one of: button, usb, periodic, status_update, alarm, alarm_test, undefined (in case of deep sleep, reset was not caused by exit from deep sleep)
	Cause *string `json:"cause,omitempty" yaml:"cause,omitempty"`
}

// Clone return copy
func (t *SystemWakeupReason) Clone() *SystemWakeupReason {
	c := &SystemWakeupReason{}
	copier.Copy(&c, &t)
	return c
}

// Config System component config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type Config struct {
	// Device information about the device
	Device *Device `json:"device,omitempty" yaml:"device,omitempty"`
	// Location information about the current location of the device
	Location *SystemLocation `json:"location,omitempty" yaml:"location,omitempty"`
	// Debug configuration of the device's debug logs.
	Debug *SystemDebug `json:"debug,omitempty" yaml:"debug,omitempty"`
	// UIData user interface data
	UIData *SystemUIData `json:"ui_data,omitempty" yaml:"ui_data,omitempty"`
	// RPCUDP configuration for the RPC over UDP
	RPCUDP *SystemRPCUDP `json:"rpc_udp,omitempty" yaml:"rpc_udp,omitempty"`
	// Sntp configuration for the sntp server
	Sntp *SystemSntp `json:"sntp,omitempty" yaml:"sntp,omitempty"`
	// CfgRev Configuration revision. This number will be incremented for every configuration change of a device component.
	// If the new config value is the same as the old one there will be no change of this property. Can not be modified
	// explicitly by a call to Sys.SetConfig
	CfgRev *int `json:"cfg_rev,omitempty" yaml:"cfg_rev,omitempty"`
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

	t.Device.Sanatize()
	t.CfgRev = nil
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

	if !t.Device.Equals(x.Device) {
		zap.L().Info("Config Device not equal")
		result = false
	}

	if !t.Location.Equals(x.Location) {
		zap.L().Info("Config Location not equal")
		result = false
	}

	if !t.Debug.Equals(x.Debug) {
		zap.L().Info("Config Debug not equal")
		result = false
	}

	if !t.UIData.Equals(x.UIData) {
		zap.L().Info("Config UIData not equal")
		result = false
	}

	if !t.RPCUDP.Equals(x.RPCUDP) {
		zap.L().Info("Config RPCUDP not equal")
		result = false
	}

	if !t.Sntp.Equals(x.Sntp) {
		zap.L().Info("Config Sntp not equal")
		result = false
	}

	return result
}

func (t *Config) Merge(x *Config) {

	if x == nil {
		return
	}

	if t.Device == nil {
		if x.Device != nil {
			t.Device = x.Device.Clone()
		} else {
			t.Device.Merge(x.Device)
		}
	}

	if t.Location == nil {
		if x.Location != nil {
			t.Location = x.Location.Clone()
		} else {
			t.Location.Merge(x.Location)
		}
	}

	if t.Debug == nil {
		if x.Debug != nil {
			t.Debug = x.Debug.Clone()
		} else {
			t.Debug.Merge(x.Debug)
		}
	}

	if t.UIData == nil {
		if x.UIData != nil {
			t.UIData = x.UIData.Clone()
		} else {
			t.UIData.Merge(x.UIData)
		}
	}

	if t.RPCUDP == nil {
		if x.RPCUDP != nil {
			t.RPCUDP = x.RPCUDP.Clone()
		} else {
			t.RPCUDP.Merge(x.RPCUDP)
		}
	}

	if t.Sntp == nil {
		if x.Sntp != nil {
			t.Sntp = x.Sntp.Clone()
		} else {
			t.Sntp.Merge(x.Sntp)
		}
	}

	if t.CfgRev == nil {
		t.CfgRev = x.CfgRev
	}

}

// Device information about the device
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type Device struct {
	// Name of the device
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// EcoMode experimental Decreases power consumption when set to true, at the cost of reduced execution speed and increased network latency
	EcoMode *bool `json:"eco_mode,omitempty" yaml:"eco_mode,omitempty"`
	// MAC read-only base MAC address of the device
	MAC *string `json:"mac,omitempty" yaml:"mac,omitempty"`
	// FwID read-only build identifier of the current firmware image
	FwID *string `json:"fw_id,omitempty" yaml:"fw_id,omitempty"`
	// Profile name of the device profile (only applicable for multi-profile devices)
	Profile *string `json:"profile,omitempty" yaml:"profile,omitempty"`
	// Discoverable if true, device is shown in 'Discovered devices'. If false, the device is hidden.
	Discoverable *bool `json:"discoverable,omitempty" yaml:"discoverable,omitempty"`
	// AddonType enable/disable addon board (if supported). Range of values: sensor; null to disable.
	AddonType *string `json:"addon_type,omitempty" yaml:"addon_type,omitempty"`
}

// Clone return copy
func (t *Device) Clone() *Device {
	c := &Device{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *Device) Sanatize() {

	if t == nil {
		return
	}

	t.MAC = nil
	t.FwID = nil
}

// Equals returns true if equal
func (t *Device) Equals(x *Device) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("Config receiver is nil but input is not")
		return x == nil
	}

	if x == nil {
		zap.L().Info("Config receiver is not nil but input is")
		return false
	}

	result := true

	if !util.CompareString(t.Name, x.Name) {
		zap.L().Info("Device Name not equal")
		result = false
	}

	if !util.CompareBool(t.EcoMode, x.EcoMode) {
		zap.L().Info("Device EcoMode not equal")
		result = false
	}

	if !util.CompareString(t.Profile, x.Profile) {
		zap.L().Info("Device Profile not equal")
		result = false
	}

	if !util.CompareBool(t.Discoverable, x.Discoverable) {
		zap.L().Info("Device Discoverable not equal")
		result = false
	}

	if !util.CompareString(t.AddonType, x.AddonType) {
		zap.L().Info("Device AddonType not equal")
		result = false
	}

	return result
}

func (t *Device) Merge(x *Device) {

	if x == nil {
		return
	}

	if t.Name == nil {
		t.Name = x.Name
	}

	if t.EcoMode == nil {
		t.EcoMode = x.EcoMode
	}

	if t.MAC == nil {
		t.MAC = x.MAC
	}

	if t.FwID == nil {
		t.FwID = x.FwID
	}

	if t.Profile == nil {
		t.Profile = x.Profile
	}

	if t.Discoverable == nil {
		t.Discoverable = x.Discoverable
	}

	if t.AddonType == nil {
		t.AddonType = x.AddonType
	}

}

// SystemLocationConfig Information about the current location of the device
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemLocation struct {
	// Timezone (null if unavailable)
	Tz *string `json:"tz,omitempty" yaml:"tz,omitempty"`
	// Lat latitude in degrees (null if unavailable)
	Lat *float64 `json:"lat,omitempty" yaml:"lat,omitempty"`
	// Lon longitude in degrees (null if unavailable)
	Lon *float64 `json:"lon,omitempty" yaml:"lon,omitempty"`
}

// Clone return copy
func (t *SystemLocation) Clone() *SystemLocation {
	c := &SystemLocation{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemLocation) Equals(x *SystemLocation) bool {

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

	if !util.CompareString(t.Tz, x.Tz) {
		zap.L().Info("SystemLocation Tz not equal")
		result = false
	}

	if !util.CompareFloat64(t.Lat, x.Lat) {
		zap.L().Info("SystemLocation Lat not equal")
		result = false
	}

	if !util.CompareFloat64(t.Lon, x.Lon) {
		zap.L().Info("SystemLocation Lon not equal")
		result = false
	}

	return result
}

func (t *SystemLocation) Merge(x *SystemLocation) {

	if x == nil {
		return
	}

	if t.Tz == nil {
		t.Tz = x.Tz
	}

	if t.Lat == nil {
		t.Lat = x.Lat
	}

	if t.Lon == nil {
		t.Lon = x.Lon
	}
}

// DebugConfig Configuration of the device's debug logs
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
// https://shelly-api-docs.shelly.cloud/gen2/General/DebugLogs
type SystemDebug struct {
	// Mqtt configuration of logs streamed over MQTT
	Mqtt *SystemMqtt `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	// Websocket configuration of logs streamed over websocket. Attention: Access to log streams over
	// websocket is not restricted, even when authentication is enabled!
	Websocket *SystemWebsocket `json:"websocket,omitempty" yaml:"websocket,omitempty"`
	// UDP Configuration of logs streamed over UDP
	UDP *UDP `json:"udp,omitempty" yaml:"udp,omitempty"`
}

func (t *SystemDebug) Merge(x *SystemDebug) {

	if x == nil {
		return
	}

	if t.Mqtt == nil {
		if x.Mqtt != nil {
			t.Mqtt = x.Mqtt.Clone()
		} else {
			t.Mqtt.Merge(x.Mqtt)
		}
	}

	if t.Websocket == nil {
		if x.Websocket != nil {
			t.Websocket = x.Websocket.Clone()
		} else {
			t.Websocket.Merge(x.Websocket)
		}
	}

	if t.UDP == nil {
		if x.UDP != nil {
			t.UDP = x.UDP.Clone()
		} else {
			t.UDP.Merge(x.UDP)
		}
	}

}

// Clone return copy
func (t *SystemDebug) Clone() *SystemDebug {
	c := &SystemDebug{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemDebug) Equals(x *SystemDebug) bool {

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

	if !t.Mqtt.Equals(x.Mqtt) {
		zap.L().Info("SystemDebug Mqtt not equal")
		result = false
	}

	if !t.Websocket.Equals(x.Websocket) {
		zap.L().Info("SystemDebug Websocket not equal")
		result = false
	}

	if !t.UDP.Equals(x.UDP) {
		zap.L().Info("SystemDebug UDP not equal")
		result = false
	}

	return result
}

// SystemMqtt Configuration of logs streamed over MQTT
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemMqtt struct {
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// Clone return copy
func (t *SystemMqtt) Clone() *SystemMqtt {
	c := &SystemMqtt{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemMqtt) Equals(x *SystemMqtt) bool {

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

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("SystemMqtt Enable not equal")
		return false
	}

	return true
}

func (t *SystemMqtt) Merge(x *SystemMqtt) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

}

// SystemWebsocket Configuration of logs streamed over websocket. Attention: Access to log streams
// over websocket is not restricted, even when authentication is enabled!
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemWebsocket struct {
	// True if enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// Clone return copy
func (t *SystemWebsocket) Clone() *SystemWebsocket {
	c := &SystemWebsocket{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemWebsocket) Equals(x *SystemWebsocket) bool {

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

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("SystemWebsocket Enable not equal")
		return false
	}

	return true
}

func (t *SystemWebsocket) Merge(x *SystemWebsocket) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

}

// UDP Configuration of logs streamed over UDP. Used by component System.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type UDP struct {
	Addr *string `json:"addr,omitempty" yaml:"addr,omitempty"`
}

// Clone return copy
func (t *UDP) Clone() *UDP {
	c := &UDP{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *UDP) Equals(x *UDP) bool {

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

	if !util.CompareString(t.Addr, x.Addr) {
		zap.L().Info("UDP Addr not equal")
		return false
	}

	return true
}

func (t *UDP) Merge(x *UDP) {

	if x == nil {
		return
	}

	if t.Addr == nil {
		t.Addr = x.Addr
	}

}

// SystemUIData user interface data. Used by component System.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemUIData struct {
}

// Clone return copy
func (t *SystemUIData) Clone() *SystemUIData {
	c := &SystemUIData{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemUIData) Equals(x *SystemUIData) bool {

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

	return true
}

func (t *SystemUIData) Merge(x *SystemUIData) {

	if x == nil {
		return
	}
}

// SystemRPCUDP configuration for the RPC over UDP
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemRPCUDP struct {
	// DstAddr destination IP address
	DstAddr *string `json:"dst_addr,omitempty" yaml:"dst_addr,omitempty"`
	// ListenPort port number for inbound UDP RPC channel, null disables. Restart is required for changes to apply
	ListenPort *string `json:"listen_port,omitempty" yaml:"listen_port,omitempty"`
}

// Clone return copy
func (t *SystemRPCUDP) Clone() *SystemRPCUDP {
	c := &SystemRPCUDP{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemRPCUDP) Equals(x *SystemRPCUDP) bool {

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

	if !util.CompareString(t.DstAddr, x.DstAddr) {
		zap.L().Info("SystemRPCUDP DstAddr not equal")
		result = false
	}

	if !util.CompareString(t.ListenPort, x.ListenPort) {
		zap.L().Info("SystemRPCUDP ListenPort not equal")
		result = false
	}

	return result
}

func (t *SystemRPCUDP) Merge(x *SystemRPCUDP) {

	if x == nil {
		return
	}

	if t.DstAddr == nil {
		t.DstAddr = x.DstAddr
	}

	if t.ListenPort == nil {
		t.ListenPort = x.ListenPort
	}
}

// SntpConfig configuration for the sntp server
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemSntp struct {
	// Server name of the sntp server
	Server *string `json:"server,omitempty" yaml:"server,omitempty"`
}

// Clone return copy
func (t *SystemSntp) Clone() *SystemSntp {
	c := &SystemSntp{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *SystemSntp) Equals(x *SystemSntp) bool {

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

	if !util.CompareString(t.Server, x.Server) {
		zap.L().Info("SystemSntp Server not equal")
		return false
	}

	return true
}

func (t *SystemSntp) Merge(x *SystemSntp) {

	if x == nil {
		return
	}

	if t.Server == nil {
		t.Server = x.Server
	}

}

// FirmwareStatus is common for components Sys and Shelly
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#status &
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#status
type FirmwareStatus struct {
	// Version of the new firmware
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`
	// BuildID Id of the new build
	BuildID *string `json:"build_id,omitempty" yaml:"build_id,omitempty"`
}
