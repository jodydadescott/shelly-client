package types

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/copier"
)

// ShellyStatus status of all the components of the device.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly
type ShellyStatus struct {
	Bluetooth *BluetoothStatus `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud     *CloudStatus     `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt      *MqttStatus      `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet  *EthernetStatus  `json:"eth,omitempty" yaml:"eth,omitempty"`
	System    *SystemStatus    `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi      *WifiStatus      `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Light     []*LightStatus   `json:"light,omitempty" yaml:"light,omitempty"`
	Input     []*InputStatus   `json:"input,omitempty" yaml:"input,omitempty"`
	Switch    []*SwitchStatus  `json:"switch,omitempty" yaml:"switch,omitempty"`
}

// ShellyRPCMethods lists of all available RPC methods. It takes into account both ACL and authentication
// restrictions and only lists the methods allowed for the particular user/channel that's making the request.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#shellylistmethods
type ShellyRPCMethods struct {
	// Methods names of the methods allowed
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty"`
}

// Clone return copy
func (t *ShellyRPCMethods) Clone() *ShellyRPCMethods {
	c := &ShellyRPCMethods{}
	copier.Copy(&c, &t)
	return c
}

// ShellyConfig Shelly component config. The config is composed of each components config.
// Shelly devices can have zero or more 'Light', 'Input' and 'Switch' types. Because these
// are explicity named and not members of a JSON array we have statically created them.
// This seemed to be a cleaner solution then a customized JSON/YAML encoder/decoder. We have
// created 8 for each which is currently more then enough as the max for any Shelly product as
// of today is 4.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyConfig struct {
	Auth          *ShellyAuthConfig          `json:"auth,omitempty" yaml:"auth,omitempty"`
	TLSClientCert *ShellyTLSClientCertConfig `json:"tls_client_cert,omitempty" yaml:"tls_client_cert,omitempty"`
	TLSClientKey  *ShellyTLSClientKeyConfig  `json:"tls_client_key,omitempty" yaml:"tls_client_key,omitempty"`
	UserCA        *ShellyUserCAConfig        `json:"user_ca,omitempty" yaml:"user_ca,omitempty"`
	Bluetooth     *BluetoothConfig           `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud         *CloudConfig               `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt          *MqttConfig                `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet      *EthernetConfig            `json:"eth,omitempty" yaml:"eth,omitempty"`
	System        *SystemConfig              `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi          *WifiConfig                `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Websocket     *WebsocketConfig           `json:"ws,omitempty" yaml:"ws,omitempty"`
	Light         []*LightConfig             `json:"light,omitempty" yaml:"light,omitempty"`
	Input         []*InputConfig             `json:"input,omitempty" yaml:"input,omitempty"`
	Switch        []*SwitchConfig            `json:"switch,omitempty" yaml:"switch,omitempty"`
}

// Clone return copy
func (t *ShellyConfig) Clone() *ShellyConfig {
	c := &ShellyConfig{}
	copier.Copy(&c, &t)
	return c
}

// GetLight returns Light with specified ID, otherwise nil
func (t *ShellyConfig) GetLight(id int) *LightConfig {
	for _, v := range t.Light {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// GetInput returns Input with specified ID, otherwise nil
func (t *ShellyConfig) GetInput(id int) *InputConfig {
	for _, v := range t.Input {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// GetInput returns Input with specified ID, otherwise nil
func (t *ShellyConfig) GetSwitch(id int) *SwitchConfig {
	for _, v := range t.Switch {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// Markup markup config
func (t *ShellyConfig) Markup() {

	if t == nil {
		return
	}

	if t.Auth == nil {
		t.Auth = &ShellyAuthConfig{}
	}

	t.Auth.Markup()

	if t.UserCA == nil {
		t.UserCA = &ShellyUserCAConfig{}
	}

	t.UserCA.Markup()

	if t.TLSClientCert == nil {
		t.TLSClientCert = &ShellyTLSClientCertConfig{}
	}

	t.TLSClientCert.Markup()

	if t.TLSClientKey == nil {
		t.TLSClientKey = &ShellyTLSClientKeyConfig{}
	}

	t.TLSClientKey.Markup()

	t.Bluetooth.Markup()
	t.Cloud.Markup()
	t.Mqtt.Markup()
	t.Ethernet.Markup()
	t.System.Markup()
	t.Wifi.Markup()
	t.Websocket.Markup()

	for _, v := range t.Light {
		v.Markup()
	}

	for _, v := range t.Input {
		v.Markup()
	}

	for _, v := range t.Switch {
		v.Markup()
	}

}

// Sanatize sanatize config
func (t *ShellyConfig) Sanatize() {

	if t == nil {
		return
	}

	t.Auth.Sanatize()
	t.TLSClientCert.Sanatize()
	t.TLSClientKey.Sanatize()
	t.UserCA.Sanatize()
	t.Bluetooth.Sanatize()
	t.Cloud.Sanatize()
	t.Mqtt.Sanatize()
	t.Ethernet.Sanatize()
	t.System.Sanatize()
	t.Wifi.Sanatize()
	t.Websocket.Sanatize()

	for _, v := range t.Light {
		v.Sanatize()
	}

	for _, v := range t.Input {
		v.Sanatize()
	}

	for _, v := range t.Switch {
		v.Sanatize()
	}
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

// ShellyUpdateConfig Shelly firmware update config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyUpdateConfig struct {
	// Stage is used by the following methods:
	// Update : The type of the new version - either stable or beta. By default updates to stable version. Optional
	Stage *string `json:"stage,omitempty" yaml:"stage,omitempty"`
	// Url is used by the following methods:
	// Update : Url address of the update. Optional
	Url *string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Clone return copy
func (t *ShellyUpdateConfig) Clone() *ShellyUpdateConfig {
	c := &ShellyUpdateConfig{}
	copier.Copy(&c, &t)
	return c
}

// ShellyAuthConfig Shelly Auth Config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyAuthConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// Pass password
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
}

// Clone return copy
func (t *ShellyAuthConfig) Clone() *ShellyAuthConfig {
	c := &ShellyAuthConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *ShellyAuthConfig) Markup() {

	if t == nil {
		return
	}

	if t.Pass == nil {

		if t.Enable {
			tmp := passwordIfEnabled
			t.Pass = &tmp
		} else {
			tmp := passwordIfNotEnabled
			t.Pass = &tmp
		}
	}
}

// Sanatize sanatize config
func (t *ShellyAuthConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Pass != nil {
		if *t.Pass == passwordIfNotEnabled {
			t.Pass = nil
		}
	}
}

// ShellyUserCAConfig Shelly UserCA config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyUserCAConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// Data is used by the following methods:
	// PutUserCA : Contents of the PEM file (null if you want to delete the existing data). Required
	// PutTLSClientCert : Contents of the client.crt file (null if you want to delete the existing data). Required
	// PutTLSClientKey : Contents of the client.key file (null if you want to delete the existing data). Required
	Data *string `json:"data,omitempty" yaml:"data,omitempty"`
}

// Clone return copy
func (t *ShellyUserCAConfig) Clone() *ShellyUserCAConfig {
	c := &ShellyUserCAConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *ShellyUserCAConfig) Markup() {

	if t == nil {
		return
	}

	if t.Data == nil {
		tmp := shellyUserCAExample
		t.Data = &tmp
	}

}

// Sanatize sanatize config
func (t *ShellyUserCAConfig) Sanatize() {

	if t == nil {
		return
	}

	if !t.Enable {
		t.Data = nil
	}
}

// ShellyTLSClientCertConfig Shelly TLS Client Cert config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyTLSClientCertConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// Data is used by the following methods:
	// PutUserCA : Contents of the PEM file (null if you want to delete the existing data). Required
	// PutTLSClientCert : Contents of the client.crt file (null if you want to delete the existing data). Required
	// PutTLSClientKey : Contents of the client.key file (null if you want to delete the existing data). Required
	Data *string `json:"data,omitempty" yaml:"data,omitempty"`
}

// Clone return copy
func (t *ShellyTLSClientCertConfig) Clone() *ShellyTLSClientCertConfig {
	c := &ShellyTLSClientCertConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *ShellyTLSClientCertConfig) Markup() {

	if t == nil {
		return
	}

	if t.Data == nil {
		tmp := shellyTLSClientCertExample
		t.Data = &tmp
	}

}

// Sanatize sanatize config
func (t *ShellyTLSClientCertConfig) Sanatize() {

	if t == nil {
		return
	}

	if !t.Enable {
		t.Data = nil
	}
}

// ShellyTLSClientKeyConfig Shelly TLS Client Key config
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type ShellyTLSClientKeyConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// Data is used by the following methods:
	// PutUserCA : Contents of the PEM file (null if you want to delete the existing data). Required
	// PutTLSClientCert : Contents of the client.crt file (null if you want to delete the existing data). Required
	// PutTLSClientKey : Contents of the client.key file (null if you want to delete the existing data). Required
	Data *string `json:"data,omitempty" yaml:"data,omitempty"`
}

// Clone return copy
func (t *ShellyTLSClientKeyConfig) Clone() *ShellyTLSClientKeyConfig {
	c := &ShellyTLSClientKeyConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *ShellyTLSClientKeyConfig) Markup() {

	if t == nil {
		return
	}

	if t.Data == nil {
		tmp := shellyTLSClientKey
		t.Data = &tmp
	}

}

// Sanatize sanatize config
func (t *ShellyTLSClientKeyConfig) Sanatize() {

	if t == nil {
		return
	}

	if !t.Enable {
		t.Data = nil
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

type ShellyReport struct {
	Auth          *ComponentReport   `json:"auth,omitempty" yaml:"auth,omitempty"`
	TLSClientCert *ComponentReport   `json:"tls_client_cert,omitempty" yaml:"tls_client_cert,omitempty"`
	TLSClientKey  *ComponentReport   `json:"tls_client_key,omitempty" yaml:"tls_client_key,omitempty"`
	UserCA        *ComponentReport   `json:"user_ca,omitempty" yaml:"user_ca,omitempty"`
	Bluetooth     *ComponentReport   `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud         *ComponentReport   `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt          *ComponentReport   `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet      *ComponentReport   `json:"eth,omitempty" yaml:"eth,omitempty"`
	System        *ComponentReport   `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi          *ComponentReport   `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Websocket     *ComponentReport   `json:"ws,omitempty" yaml:"ws,omitempty"`
	Light         []*ComponentReport `json:"light,omitempty" yaml:"light,omitempty"`
	Input         []*ComponentReport `json:"input,omitempty" yaml:"input,omitempty"`
	Switch        []*ComponentReport `json:"switch,omitempty" yaml:"switch,omitempty"`
}

// Clone return copy
func (t *ShellyReport) Clone() *ShellyReport {
	c := &ShellyReport{}
	copier.Copy(&c, &t)
	return c
}

func (t *ShellyReport) RebootRequired() bool {

	if t.Bluetooth != nil {
		if t.Bluetooth.RebootRequired != nil {
			if *t.Bluetooth.RebootRequired {
				return true
			}
		}
	}

	if t.Cloud != nil {
		if t.Cloud.RebootRequired != nil {
			if *t.Cloud.RebootRequired {
				return true
			}
		}
	}

	if t.Mqtt != nil {
		if t.Mqtt.RebootRequired != nil {
			if *t.Mqtt.RebootRequired {
				return true
			}
		}
	}

	if t.Ethernet != nil {
		if t.Ethernet.RebootRequired != nil {
			if *t.Ethernet.RebootRequired {
				return true
			}
		}
	}

	if t.System != nil {
		if t.System.RebootRequired != nil {
			if *t.System.RebootRequired {
				return true
			}
		}
	}

	if t.Wifi != nil {
		if t.Wifi.RebootRequired != nil {
			if *t.Wifi.RebootRequired {
				return true
			}
		}
	}

	if t.Websocket != nil {
		if t.Websocket.RebootRequired != nil {
			if *t.Websocket.RebootRequired {
				return true
			}
		}
	}

	return false
}

func (t *ShellyReport) Error() error {

	var errors *multierror.Error

	if t.Auth != nil {
		if t.Auth.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("auth :: %v", t.Auth.Error))
		}
	}

	if t.TLSClientCert != nil {
		if t.TLSClientCert.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("tlsClientCert :: %v", t.TLSClientCert.Error))
		}
	}

	if t.TLSClientKey != nil {
		if t.TLSClientKey.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("tlsClientKey :: %v", t.TLSClientKey.Error))
		}
	}

	if t.UserCA != nil {
		if t.UserCA.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("userCA :: %v", t.UserCA.Error))
		}
	}

	if t.Bluetooth != nil {
		if t.Bluetooth.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("bluetooth :: %v", t.Bluetooth.Error))
		}
	}

	if t.Cloud != nil {
		if t.Cloud.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("cloud :: %v", t.Cloud.Error))
		}
	}

	if t.Mqtt != nil {
		if t.Mqtt.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("mqtt :: %v", t.Mqtt.Error))
		}
	}

	if t.Ethernet != nil {
		if t.Ethernet.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("ethernet :: %v", t.Ethernet.Error))
		}
	}

	if t.System != nil {
		if t.System.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("system :: %v", t.System.Error))
		}
	}

	if t.Wifi != nil {
		if t.Wifi.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("wifi :: %v", t.Wifi.Error))
		}
	}

	if t.Websocket != nil {
		if t.Websocket.Error != nil {
			errors = multierror.Append(errors, fmt.Errorf("websocket :: %v", t.Websocket.Error))
		}
	}

	if t.Light != nil {
		for _, v := range t.Light {
			if v.Error != nil {
				errors = multierror.Append(errors, fmt.Errorf("light %d :: %v", v.ID, v.Error))
			}
		}
	}

	if t.Input != nil {
		for _, v := range t.Input {
			if v.Error != nil {
				errors = multierror.Append(errors, fmt.Errorf("input %d :: %v", v.ID, v.Error))
			}
		}
	}

	if t.Switch != nil {
		for _, v := range t.Switch {
			if v.Error != nil {
				errors = multierror.Append(errors, fmt.Errorf("switch %d :: %v", v.ID, v.Error))
			}
		}
	}

	return errors.ErrorOrNil()
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
