package types

import (
	"strings"

	"github.com/jinzhu/copier"
)

// WifiStatus status of the WiFi component contains information about the state of the WiFi connection of the device.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#status
type WifiStatus struct {
	// StaIP Ip of the device in the network (null if disconnected)
	StaIP *string `json:"sta_ip" yaml:"sta_ip"`
	// Status of the connection. Range of values: disconnected, connecting, connected, got ip
	Status string `json:"status" yaml:"status"`
	// Ssid of the network (null if disconnected)
	SSID *string `json:"ssid" yaml:"ssid"`
	// Rssi Strength of the signal in dBms
	RSSI *int `json:"rssi" yaml:"rssi"`
	// ApClientCount Number of clients connected to the access point. Present only when AP is
	// enabled and range extender functionality is present and enabled.
	ApClientCount *int `json:"ap_client_count" yaml:"ap_client_count"`
}

// Clone return copy
func (t *WifiStatus) Clone() *WifiStatus {
	c := &WifiStatus{}
	copier.Copy(&c, &t)
	return c
}

func (t *WifiStatus) GetStatus() WifiStatusStatus {
	return WifiStatusStatusFromString(t.Status)
}

type WifiStatusStatus string

const (
	WifiStatusStatusInvalid      WifiStatusStatus = "invalid"
	WifiStatusStatusDisconnected WifiStatusStatus = "disconnected"
	WifiStatusStatusConnecting   WifiStatusStatus = "connecting"
	WifiStatusStatusConnected    WifiStatusStatus = "connected"
	WifiStatusStatusGotIP        WifiStatusStatus = "got ip"
)

func WifiStatusStatusFromString(s string) WifiStatusStatus {

	switch strings.ToLower(s) {

	case string(WifiStatusStatusDisconnected):
		return WifiStatusStatusDisconnected

	case string(WifiStatusStatusConnecting):
		return WifiStatusStatusConnecting

	case string(WifiStatusStatusConnected):
		return WifiStatusStatusConnected

	case string(WifiStatusStatusGotIP):
		return WifiStatusStatusGotIP

	}

	return WifiStatusStatusInvalid
}

// WifiConfig configuration of the WiFi component contains information about the access point of the device,
// the network stations and the roaming settings.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiConfig struct {
	// Ap Information about the access point
	Ap *WifiAPConfig `json:"ap,omitempty" yaml:"ap,omitempty"`
	// Sta information about the sta configuration
	Sta *WifiSTAConfig `json:"sta,omitempty" yaml:"sta,omitempty"`
	// Sta1 information about the sta configuration
	Sta1 *WifiSTAConfig `json:"sta1,omitempty" yaml:"sta1,omitempty"`
	// Roam WiFi roaming configuration
	Roam *WifiRoamConfig `json:"roam,omitempty" yaml:"roam,omitempty"`
}

// Clone return copy
func (t *WifiConfig) Clone() *WifiConfig {
	c := &WifiConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *WifiConfig) Markup() {

	if t == nil {
		return
	}

	if t.Ap != nil {
		t.Ap.SSID = nil
		t.Ap.Markup()
	}

	t.Sta.Markup()
	t.Sta1.Markup()
	t.Roam.Markup()
}

// Sanatize sanatize config
func (t *WifiConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Ap != nil {
		// SSID is read only
		t.Ap.SSID = nil
	}

	t.Ap.Sanatize()
	t.Sta.Sanatize()
	t.Sta1.Sanatize()
	t.Roam.Sanatize()
}

// WifiAPConfig WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiAPConfig struct {
	// SSID readonly SSID of the access point
	SSID *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	// Pass password for the ssid, writeonly. Must be provided if you provide ssid
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
	// IsOpen True if the access point is open, false otherwise
	IsOpen *bool `json:"is_open,omitempty" yaml:"is_open,omitempty"`
	// Enable true if the access point is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// RangeExtender range extender configuration object, available only when range extender functionality is present.
	RangeExtender *WifiRangeExtenderConfig `json:"range_extender,omitempty" yaml:"range_extender,omitempty"`
}

// Clone return copy
func (t *WifiAPConfig) Clone() *WifiAPConfig {
	c := &WifiAPConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *WifiAPConfig) Markup() {

	if t == nil {
		return
	}

	t.SSID = nil

	if t.Enable {
		if t.Pass == nil {
			tmp := genericPassword
			t.Pass = &tmp
		}

		t.RangeExtender.Markup()
		return
	}

	t.Pass = nil
	t.IsOpen = nil
	t.RangeExtender = nil
}

// Sanatize sanatize config
func (t *WifiAPConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable {
		t.RangeExtender.Sanatize()
		return
	}

	t.Pass = nil
	t.RangeExtender = nil
}

// WifiRangeExtenderConfig Range extender configuration object, available only when range extender functionality is present.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiRangeExtenderConfig struct {
	Enable bool `json:"enable" yaml:"enable"`
}

// Clone return copy
func (t *WifiRangeExtenderConfig) Clone() *WifiRangeExtenderConfig {
	c := &WifiRangeExtenderConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *WifiRangeExtenderConfig) Markup() {

	if t == nil {
		return
	}

}

// Sanatize sanatize config
func (t *WifiRangeExtenderConfig) Sanatize() {

	if t == nil {
		return
	}

}

// WifiSTAConfig WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiSTAConfig struct {
	// SSID of the network
	SSID *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	// Password for the ssid, writeonly. Must be provided if you provide ssid
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
	// IsOpen true if the network is open, i.e. no password is set, false otherwise, readonly
	IsOpen *bool `json:"is_open,omitempty" yaml:"is_open,omitempty"`
	// Enable True if the configuration is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
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
func (t *WifiSTAConfig) Clone() *WifiSTAConfig {
	c := &WifiSTAConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *WifiSTAConfig) Markup() {

	if t == nil {
		return
	}

	if t.Enable {
		if t.Pass == nil {
			tmp := genericPassword
			t.Pass = &tmp
		}
		return
	}

	t.SSID = nil
	t.Pass = nil
	t.IsOpen = nil
	t.Ipv4Mode = nil
	t.IP = nil
	t.Netmask = nil
	t.Gateway = nil
	t.Nameserver = nil
}

// Sanatize sanatize config
func (t *WifiSTAConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Pass != nil {
		if *t.Pass == genericPassword {
			t.Pass = nil
		}
	}

	if t.Ipv4Mode != nil {
		if *t.Ipv4Mode == "dhcp" {
			t.IP = nil
			t.Netmask = nil
			t.Gateway = nil
			t.Nameserver = nil
		}
	}

}

// WifiRoamConfig WiFi roaming configuration
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiRoamConfig struct {
	// RSSIThreshold - when reached will trigger the access point roaming. Default value: -80
	RSSIThreshold *int `json:"rssi_thr" yaml:"rssi_thr"`
	// Interval at which to scan for better access points. Enabled if set to positive number,
	// disabled if set to 0. Default value: 60
	Interval *int `json:"interval" yaml:"interval"`
}

// Clone return copy
func (t *WifiRoamConfig) Clone() *WifiRoamConfig {
	c := &WifiRoamConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *WifiRoamConfig) Markup() {

	if t == nil {
		return
	}

}

// Sanatize sanatize config
func (t *WifiRoamConfig) Sanatize() {

	if t == nil {
		return
	}

}

// Scan WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi
type WifiNet struct {
	SSID    *string `json:"ssid" yaml:"ssid"`
	BSSID   *string `json:"bssid" yaml:"bssid"`
	Auth    *int    `json:"auth" yaml:"auth"`
	Channel *int    `json:"channel" yaml:"channel"`
	RSSI    *int    `json:"rssi" yaml:"rssi"`
}

// Clone return copy
func (t *WifiNet) Clone() *WifiNet {
	c := &WifiNet{}
	copier.Copy(&c, &t)
	return c
}

// WifiAPClient WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi
type WifiAPClient struct {
	MAC      *string `json:"mac" yaml:"mac"`
	IP       *string `json:"ip" yaml:"ip"`
	IPStatic *bool   `json:"ip_static" yaml:"ip_static"`
	Mport    *int    `json:"mport" yaml:"mport"`
	Since    *int    `json:"since" yaml:"since"`
}

// Clone return copy
func (t *WifiAPClient) Clone() *WifiAPClient {
	c := &WifiAPClient{}
	copier.Copy(&c, &t)
	return c
}

// WifiScanResults Wifi Scan Results
type WifiScanResults struct {
	Results []WifiNet `json:"results"`
}

// Clone return copy
func (t *WifiScanResults) Clone() *WifiScanResults {
	c := &WifiScanResults{}
	copier.Copy(&c, &t)
	return c
}

// WifiAPClients Wifi AP Clients
type WifiAPClients struct {
	Ts      *int           `json:"ts"`
	Clients []WifiAPClient `json:"ap_clients"`
}

// Clone return copy
func (t *WifiAPClients) Clone() *WifiAPClients {
	c := &WifiAPClients{}
	copier.Copy(&c, &t)
	return c
}
