package types

import (
	"github.com/jinzhu/copier"
)

// SystemStatus status contains information about network state, system time and other common attributes of the Shelly device.
// Presence of some keys is optional, depending on the underlying hardware components.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#Status
type SystemStatus struct {
	// MAC address of the device
	MAC *string `json:"mac" yaml:"mac"`
	// RestartRequired true if restart is required, false otherwise
	RestartRequired *bool `json:"restart_required" yaml:"restart_required"`
	// Time Current time in the format HH:MM (24-hour time format in the current timezone with leading zero).
	// null when time is not synced from NTP server.
	Time *string `json:"time" yaml:"time"`
	// Unixtime Unix timestamp (in UTC), null when time is not synced from NTP server.
	Unixtime *float64 `json:"unixtime" yaml:"unixtime"`
	// Uptime Time in seconds since last reboot
	Uptime *float64 `json:"uptime" yaml:"uptime"`
	// RAMSize Total size of the RAM in the system in Bytes
	RAMSize *float64 `json:"ram_size" yaml:"ram_size"`
	// RAMFree Size of the free RAM in the system in Bytes
	RAMFree *float64 `json:"ram_free" yaml:"ram_free"`
	// FsSize Total size of the file system in Bytes
	FsSize *float64 `json:"fs_size" yaml:"fs_size"`
	// FsFree Size of the free file system in Bytes
	FsFree *float64 `json:"fs_free" yaml:"fs_free"`
	// CfgRev Configuration revision number
	CfgRev *float64 `json:"cfg_rev" yaml:"cfg_rev"`
	// KvsRev KVS (Key-Value Store) revision number
	KvsRev *float64 `json:"kvs_rev" yaml:"kvs_rev"`
	// ScheduleRev Schedules revision number, present if schedules are enabled
	ScheduleRev *float64 `json:"schedule_rev" yaml:"schedule_rev"`
	// WebhookRev Webhooks revision number, present if webhooks are enabled
	WebhookRev *float64 `json:"webhook_rev" yaml:"webhook_rev"`
	// AvailableUpdates Information about available updates, similar to the one returned by Shelly.CheckForUpdate
	// (empty object: {}, if no updates available). This information is automatically updated every 24 hours.
	// Note that build_id and url for an update are not displayed here
	AvailableUpdates *SystemAvailableUpdates `json:"available_updates" yaml:"available_updates"`
	// WakeupReason Information about boot type and cause (only for battery-operated devices)
	WakeupReason *SystemWakeupReason `json:"wakeup_reason" yaml:"wakeup_reason"`
	// WakeupPeriod Period (in seconds) at which device wakes up and sends "keep-alive" packet to cloud, readonly.
	// Count starts from last full wakeup
	WakeupPeriod *int `json:"wakeup_period" yaml:"wakeup_period"`
}

// Clone return copy
func (t *SystemStatus) Clone() *SystemStatus {
	c := &SystemStatus{}
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

// SystemConfig System component config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemConfig struct {
	// Device information about the device
	Device *SystemDevice `json:"device,omitempty" yaml:"device,omitempty"`
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
func (t *SystemConfig) Clone() *SystemConfig {
	c := &SystemConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *SystemConfig) Markup() {

	if t == nil {
		return
	}

	t.Device.Markup()
	t.CfgRev = nil
	t.Debug = nil
	t.UIData = nil
	t.RPCUDP = nil
}

// Sanatize sanatize config
func (t *SystemConfig) Sanatize() {

	if t == nil {
		return
	}

	t.CfgRev = nil
	t.Device.Sanatize()
}

// SystemDevice information about the device
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemDevice struct {
	// Name of the device
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// EcoMode experimental Decreases power consumption when set to true, at the cost of reduced execution speed and increased network latency
	EcoMode *bool `json:"eco_mode" yaml:"eco_mode"`
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

// // Markup markup config
// func (t *SystemDevice) SetFromEnv() error {

// 	if t == nil {
// 		return nil
// 	}

// 	nameVar := os.Getenv(ShellyEnvVar + ".name")
// 	ecomodeVar := os.Getenv(ShellyEnvVar + ".eco_mode")
// 	profileVar := os.Getenv(ShellyEnvVar + ".profile")
// 	discoverableVar := os.Getenv(ShellyEnvVar + ".discoverable")
// 	addonTypeVar := os.Getenv(ShellyEnvVar + ".addon_type")

// 	if nameVar != "" {
// 		t.Name = &nameVar
// 	}

// 	if ecomodeVar != "" {
// 		x, err := strconv.ParseBool(ecomodeVar)
// 		if err != nil {
// 			return err
// 		}
// 		t.EcoMode = &x
// 	}

// 	if profileVar != "" {
// 		t.Profile = &profileVar
// 	}

// 	if discoverableVar != "" {
// 		x, err := strconv.ParseBool(discoverableVar)
// 		if err != nil {
// 			return err
// 		}
// 		t.Discoverable = &x
// 	}

// 	if addonTypeVar != "" {
// 		t.Profile = &profileVar
// 	}

// 	return nil
// }

// Clone return copy
func (t *SystemDevice) Clone() *SystemDevice {
	c := &SystemDevice{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *SystemDevice) Markup() {

	if t == nil {
		return
	}

	t.MAC = nil
	t.FwID = nil
}

// Sanatize sanatize config
func (t *SystemDevice) Sanatize() {

	if t == nil {
		return
	}

	t.MAC = nil
	t.FwID = nil
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
	UDP *SystemUDP `json:"udp,omitempty" yaml:"udp,omitempty"`
}

// Clone return copy
func (t *SystemDebug) Clone() *SystemDebug {
	c := &SystemDebug{}
	copier.Copy(&c, &t)
	return c
}

// SystemMqtt Configuration of logs streamed over MQTT
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemMqtt struct {
	Enable *bool `json:"enable" yaml:"enable"`
}

// Clone return copy
func (t *SystemMqtt) Clone() *SystemMqtt {
	c := &SystemMqtt{}
	copier.Copy(&c, &t)
	return c
}

// SystemWebsocket Configuration of logs streamed over websocket. Attention: Access to log streams
// over websocket is not restricted, even when authentication is enabled!
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemWebsocket struct {
	// True if enabled, false otherwise
	Enable *bool `json:"enable" yaml:"enable"`
}

// Clone return copy
func (t *SystemWebsocket) Clone() *SystemWebsocket {
	c := &SystemWebsocket{}
	copier.Copy(&c, &t)
	return c
}

// SystemUDP Configuration of logs streamed over UDP. Used by component System.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#configuration
type SystemUDP struct {
	Addr *string `json:"addr,omitempty" yaml:"addr,omitempty"`
}

// Clone return copy
func (t *SystemUDP) Clone() *SystemUDP {
	c := &SystemUDP{}
	copier.Copy(&c, &t)
	return c
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
