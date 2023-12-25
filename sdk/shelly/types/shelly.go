package types

import (
	"fmt"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/sdk/util"
)

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

// SetConfigResponse internal use only
type SetConfigResponse struct {
	Response
	Result *Result `json:"result,omitempty"`
}

// GetConfigResponse internal use only
type GetConfigResponse struct {
	Response
	Result *RawConfig `json:"result,omitempty"`
}

// GetStatusResponse internal use only
type GetStatusResponse struct {
	Response
	Result *RawShellyStatus `json:"result,omitempty"`
}

// DeviceInfoResponse internal use only
type DeviceInfoResponse struct {
	Response
	Result *DeviceInfo `json:"result,omitempty"`
}

// CheckForUpdateResponse Shelly component object
type CheckForUpdateResponse struct {
	Response
	Result *SystemAvailableUpdates `json:"result,omitempty"`
}

// ListMethodsResponse internal use only
type ListMethodsResponse struct {
	Response
	Result *RPCMethods `json:"result,omitempty"`
}

type ConfigReport struct {
	RebootRequired bool `json:"rebootRequired,omitempty" yaml:"rebootRequired,omitempty"`
	NoChange       bool `json:"noChange,omitempty" yaml:"noChange,omitempty"`
}

// Clone return copy
func (t *ConfigReport) Clone() *ConfigReport {
	c := &ConfigReport{}
	copier.Copy(&c, &t)
	return c
}

// Status status of all the components of the device.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly
type Status struct {
	Bluetooth *BluetoothStatus      `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud     *CloudStatus          `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt      *MqttStatus           `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet  *EthernetStatus       `json:"eth,omitempty" yaml:"eth,omitempty"`
	System    *SystemStatus         `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi      *WifiStatus           `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Light     map[int]*LightStatus  `json:"light,omitempty" yaml:"light,omitempty"`
	Input     map[int]*InputStatus  `json:"input,omitempty" yaml:"input,omitempty"`
	Switch    map[int]*SwitchStatus `json:"switch,omitempty" yaml:"switch,omitempty"`
}

// RPCMethods lists of all available RPC methods. It takes into account both ACL and authentication
// restrictions and only lists the methods allowed for the particular user/channel that's making the request.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellylistmethods
type RPCMethods struct {
	// Methods names of the methods allowed
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty"`
}

// Clone return copy
func (t *RPCMethods) Clone() *RPCMethods {
	c := &RPCMethods{}
	copier.Copy(&c, &t)
	return c
}

// Config Shelly component config. The config is composed of each components config.
// Shelly devices can have zero or more 'Light', 'Input' and 'Switch' types. Because these
// are explicity named and not members of a JSON array we have statically created them.
// This seemed to be a cleaner solution then a customized JSON/YAML encoder/decoder. We have
// created 8 for each which is currently more then enough as the max for any Shelly product as
// of today is 4.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type Config struct {
	Auth          *AuthConfig           `json:"auth,omitempty" yaml:"auth,omitempty"`
	TLSClientCert *TLSConfig            `json:"tls_client_cert,omitempty" yaml:"tls_client_cert,omitempty"`
	TLSClientKey  *TLSConfig            `json:"tls_client_key,omitempty" yaml:"tls_client_key,omitempty"`
	UserCA        *TLSConfig            `json:"user_ca,omitempty" yaml:"user_ca,omitempty"`
	Bluetooth     *BluetoothConfig      `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud         *CloudConfig          `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt          *MqttConfig           `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet      *EthernetConfig       `json:"eth,omitempty" yaml:"eth,omitempty"`
	System        *SystemConfig         `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi          *WifiConfig           `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Websocket     *WebsocketConfig      `json:"ws,omitempty" yaml:"ws,omitempty"`
	Light         map[int]*LightConfig  `json:"light,omitempty" yaml:"light,omitempty"`
	Input         map[int]*InputConfig  `json:"input,omitempty" yaml:"input,omitempty"`
	Switch        map[int]*SwitchConfig `json:"switch,omitempty" yaml:"switch,omitempty"`
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

	if !t.Auth.Equals(x.Auth) {
		zap.L().Info("Config Auth not equal")
		result = false
	}

	if !t.TLSClientCert.Equals(x.TLSClientCert) {
		zap.L().Info("Config TLSClientCert not equal")
		result = false
	}

	if !t.TLSClientKey.Equals(x.TLSClientKey) {
		zap.L().Info("Config TLSClientKey not equal")
		result = false
	}

	if !t.UserCA.Equals(x.UserCA) {
		zap.L().Info("Config UserCA not equal")
		result = false
	}

	if !t.Bluetooth.Equals(x.Bluetooth) {
		zap.L().Info("Config Bluetooth not equal")
		result = false
	}

	if !t.Cloud.Equals(x.Cloud) {
		zap.L().Info("Config Cloud not equal")
		result = false
	}

	if !t.Mqtt.Equals(x.Mqtt) {
		zap.L().Info("Config Mqtt not equal")
		result = false
	}

	if !t.Ethernet.Equals(x.Ethernet) {
		zap.L().Info("Config Ethernet not equal")
		result = false
	}

	if !t.System.Equals(x.System) {
		zap.L().Info("Config System not equal")
		result = false
	}

	if !t.Wifi.Equals(x.Wifi) {
		zap.L().Info("Config Wifi not equal")
		result = false
	}

	if !t.Websocket.Equals(x.Websocket) {
		zap.L().Info("Config Websocket not equal")
		result = false
	}

	compareLight := func() bool {

		for i, a := range t.Light {
			b := x.Light[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Light %d not equal", i))
				return false
			}
		}

		for i, a := range x.Light {
			b := t.Light[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Light %d not equal", i))
				return false
			}
		}

		return true
	}

	if !compareLight() {
		zap.L().Info("Config Light")
		result = false
	}

	compareInput := func() bool {

		for i, a := range t.Input {
			b := x.Input[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Input %d not equal", i))
				return false
			}
		}

		for i, a := range x.Input {
			b := t.Input[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Input %d not equal", i))
				return false
			}
		}

		return true
	}

	if !compareInput() {
		zap.L().Info("Config Input")
		result = false
	}

	compareSwitch := func() bool {

		for i, a := range t.Switch {
			b := x.Switch[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Switch %d not equal", i))
				return false
			}
		}

		for i, a := range x.Switch {
			b := t.Switch[i]
			if !a.Equals(b) {
				zap.L().Info(fmt.Sprintf("Config Switch %d not equal", i))
				return false
			}
		}

		return true
	}

	if !compareSwitch() {
		zap.L().Info("Config Switch")
		result = false
	}

	return result
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

func (t *Config) Merge(x *Config) *Config {

	if x == nil {
		return t
	}

	if t.Auth == nil {
		if x.Auth != nil {
			t.Auth = x.Auth.Clone()
		} else {
			t.Auth.Merge(x.Auth)
		}
	}

	if t.TLSClientCert == nil {
		if x.TLSClientCert != nil {
			t.TLSClientCert = x.TLSClientCert.Clone()
		} else {
			t.TLSClientCert.Merge(x.TLSClientCert)
		}
	}

	if t.TLSClientKey == nil {
		if x.TLSClientKey != nil {
			t.TLSClientKey = x.TLSClientKey.Clone()
		} else {
			t.TLSClientKey.Merge(x.TLSClientKey)
		}
	}

	if t.UserCA == nil {
		if x.UserCA != nil {
			t.UserCA = x.UserCA.Clone()
		} else {
			t.UserCA.Merge(x.UserCA)
		}
	}

	if t.Bluetooth == nil {
		if x.Bluetooth != nil {
			t.Bluetooth = x.Bluetooth.Clone()
		} else {
			t.Bluetooth.Merge(x.Bluetooth)
		}
	}

	if t.Cloud == nil {
		if x.Cloud != nil {
			t.Cloud = x.Cloud.Clone()
		} else {
			t.Cloud.Merge(x.Cloud)
		}
	}

	if t.Mqtt == nil {
		if x.Mqtt != nil {
			t.Mqtt = x.Mqtt.Clone()
		} else {
			t.Mqtt.Merge(x.Mqtt)
		}
	}

	if t.Ethernet == nil {
		if x.Ethernet != nil {
			t.Ethernet = x.Ethernet.Clone()
		} else {
			t.Ethernet.Merge(x.Ethernet)
		}
	}

	if t.System == nil {
		if x.System != nil {
			t.System = x.System.Clone()
		} else {
			t.System.Merge(x.System)
		}
	}

	if t.Wifi == nil {
		if x.Wifi != nil {
			t.Wifi = x.Wifi.Clone()
		} else {
			t.Wifi.Merge(x.Wifi)
		}
	}

	if t.Websocket == nil {
		if x.Websocket != nil {
			t.Websocket = x.Websocket.Clone()
		} else {
			t.Websocket.Merge(x.Websocket)
		}
	}

	if t.Light != nil {
		if x.Light == nil {
			x.Light = make(map[int]*LightConfig)
		}
		for _, j := range t.Light {
			k := x.Light[*j.ID]
			if k == nil {
				x.Light[*j.ID] = j.Clone()
			} else {
				k.Merge(j)
			}
		}
	}

	if t.Input != nil {
		if x.Input == nil {
			x.Input = make(map[int]*InputConfig)
		}
		for _, j := range t.Input {
			k := x.Input[*j.ID]
			if k == nil {
				x.Input[*j.ID] = j.Clone()
			} else {
				k.Merge(j)
			}
		}
	}

	if t.Switch != nil {
		if x.Switch == nil {
			x.Switch = make(map[int]*SwitchConfig)
		}
		for _, j := range t.Switch {
			k := x.Switch[*j.ID]
			if k == nil {
				x.Switch[*j.ID] = j.Clone()
			} else {
				k.Merge(j)
			}
		}
	}

	return t
}

func (t *Config) Sanatize() *Config {

	t.Auth.Sanatize()
	t.TLSClientCert.Sanatize()
	t.TLSClientKey.Sanatize()

	t.UserCA.Sanatize()
	t.Bluetooth.Sanatize()
	t.Cloud.Sanatize()
	t.Mqtt.Sanatize()
	t.System.Sanatize()
	t.Wifi.Sanatize()
	t.Websocket.Sanatize()
	t.TLSClientCert.Sanatize()
	return t
}

// GetLight returns Light with specified ID, otherwise nil
func (t *Config) GetLight(id int) *LightConfig {
	for _, v := range t.Light {
		if *v.ID == id {
			return v
		}
	}
	return nil
}

// GetInput returns Input with specified ID, otherwise nil
func (t *Config) GetInput(id int) *InputConfig {
	for _, v := range t.Input {
		if *v.ID == id {
			return v
		}
	}
	return nil
}

// GetInput returns Input with specified ID, otherwise nil
func (t *Config) GetSwitch(id int) *SwitchConfig {
	for _, v := range t.Switch {
		if *v.ID == id {
			return v
		}
	}
	return nil
}

// DeviceInfo Shelly component top level device info
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellygetdeviceinfo
type DeviceInfo struct {
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	// ID Id of the device
	ID *string `json:"id" yaml:"id"`
	// MAC address of the device
	MAC *string `json:"mac,omitempty" yaml:"mac,omitempty"`
	// Model of the device
	Model *string `json:"model,omitempty" yaml:"model,omitempty"`
	// Generation of the device
	Generation *float32 `json:"gen,omitempty" yaml:"gen,omitempty"`
	// FirmwareID Id of the firmware of the device
	FirmwareID *string `json:"fw_id,omitempty" yaml:"fw_id,omitempty"`
	// Version of the firmware of the device
	Version *string `json:"ver,omitempty" yaml:"ver,omitempty"`
	// App name
	App *string `json:"app,omitempty" yaml:"app,omitempty"`
	// Profile name of the device profile (only applicable for multi-profile devices)
	Profile *string `json:"profile,omitempty" yaml:"profile,omitempty"`
	// AuthEnabled true if authentication is enabled, false otherwise
	AuthEnabled bool `json:"auth_en,omitempty" yaml:"auth_en,omitempty"`
	// AuthDomain name of the domain (null if authentication is not enabled)
	AuthDomain *string `json:"auth_domain,omitempty" yaml:"auth_domain,omitempty"`
	// Discoverable present only when false. If true, device is shown in 'Discovered devices'. If false, the device is hidden.
	Discoverable bool `json:"discoverable,omitempty" yaml:"discoverable,omitempty"`
	// Key cloud key of the device (see note below), present only when the ident parameter is set to true
	Key *string `json:"key,omitempty" yaml:"key,omitempty"`
	// Batch used to provision the device, present only when the ident parameter is set to true
	Batch *string `json:"batch,omitempty" yaml:"batch,omitempty"`
	// FwSbits Shelly internal flags, present only when the ident parameter is set to true
	FwSbits *string `json:"fw_sbits,omitempty" yaml:"fw_sbits,omitempty"`
}

// Clone return copy
func (t *DeviceInfo) Clone() *DeviceInfo {
	c := &DeviceInfo{}
	copier.Copy(&c, &t)
	return c
}

// UpdateConfig Shelly firmware update config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type UpdateConfig struct {
	// Stage is used by the following methods:
	// Update : The type of the new version - either stable or beta. By default updates to stable version. Optional
	Stage *string `json:"stage,omitempty" yaml:"stage,omitempty"`
	// Url is used by the following methods:
	// Update : Url address of the update. Optional
	Url *string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Clone return copy
func (t *UpdateConfig) Clone() *UpdateConfig {
	c := &UpdateConfig{}
	copier.Copy(&c, &t)
	return c
}

func (t *UpdateConfig) Merge(x *UpdateConfig) {

	if x == nil {
		return
	}

	if t.Stage == nil {
		t.Stage = x.Stage
	}

	if t.Url == nil {
		t.Url = x.Url
	}

}

// AuthConfig Shelly Auth Config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type AuthConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// Pass password
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
}

// Clone return copy
func (t *AuthConfig) Clone() *AuthConfig {
	c := &AuthConfig{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *AuthConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.Pass = nil
	}
}

// Equals returns true if equal
func (t *AuthConfig) Equals(x *AuthConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("AuthConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("AuthConfig receiver is not nil but input is")
		return false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("AuthConfig Enable not equal")
		return false
	}

	return true
}

func (t *AuthConfig) Merge(x *AuthConfig) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

	if t.Pass == nil {
		t.Pass = x.Pass
	}

}

// TLSConfig Shelly TLS Client Cert config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type TLSConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
	// Data is used by the following methods:
	// PutUserCA : Contents of the PEM file (null if you want to delete the existing data). Required
	// PutTLSClientCert : Contents of the client.crt file (null if you want to delete the existing data). Required
	// PutTLSClientKey : Contents of the client.key file (null if you want to delete the existing data). Required
	Data *string `json:"data,omitempty" yaml:"data,omitempty"`
}

// Clone return copy
func (t *TLSConfig) Clone() *TLSConfig {
	c := &TLSConfig{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatizes config
func (t *TLSConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.Data = nil
	}
}

// Equals returns true if equal
func (t *TLSConfig) Equals(x *TLSConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("TLSConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("TLSConfig receiver is not nil but input is")
		return false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("TLSConfig Enable not equal")
		return false
	}

	if !util.CompareString(t.Data, x.Data) {
		zap.L().Info("TLSConfig Data not equal")
		return false
	}

	return false
}

func (t *TLSConfig) Merge(x *TLSConfig) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

	if t.Data == nil {
		t.Data = x.Data
	}

}

// UpdatesReport checks for new firmware version for the device and returns information about it.
// If no update is available returns empty JSON object as result.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellycheckforupdate
type UpdatesReport struct {
	Src              *string                 `json:"src,omitempty" yaml:"src,omitempty"`
	AvailableUpdates *SystemAvailableUpdates `json:"available_updates,omitempty" yaml:"available_updates,omitempty"`
}

// Clone return copy
func (t *UpdatesReport) Clone() *UpdatesReport {
	c := &UpdatesReport{}
	copier.Copy(&c, &t)
	return c
}

type ComponentReport struct {
	RebootRequired *bool `json:"reboot_required,omitempty" yaml:"reboot_required,omitempty"`
	Error          error `json:"error,omitempty" yaml:"error,omitempty"`
	ID             *int  `json:"id,omitempty" yaml:"id,omitempty"`
}

// Clone return copy
func (t *ComponentReport) Clone() *ComponentReport {
	c := &ComponentReport{}
	copier.Copy(&c, &t)
	return c
}
