package types

import (
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/util"
)

type Response = msg_types.Response
type Error = msg_types.Error

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

// Params internal use only
type Params struct {
	Config *Config `json:"config,omitempty"`
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

// ScanResponse internal use only
type ScanResponse struct {
	Response
	Result *ScanResults `json:"result,omitempty"`
}

// ListAPClientsResponse internal use only
type ListAPClientsResponse struct {
	Response
	Result *APClients `json:"result,omitempty"`
}

// Status status of the WiFi component contains information about the state of the WiFi connection of the device.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#status
type Status struct {
	// StaIP Ip of the device in the network (null if disconnected)
	StaIP *string `json:"sta_ip,omitempty" yaml:"sta_ip,omitempty"`
	// Status of the connection. Range of values: disconnected, connecting, connected, got ip
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	// Ssid of the network (null if disconnected)
	SSID *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	// Rssi Strength of the signal in dBms
	RSSI *int `json:"rssi,omitempty" yaml:"rssi,omitempty"`
	// ApClientCount Number of clients connected to the access point. Present only when AP is
	// enabled and range extender functionality is present and enabled.
	ApClientCount *int `json:"ap_client_count,omitempty" yaml:"ap_client_count,omitempty"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

func (t *Status) GetStatus() StatusStatus {
	return StatusStatusFromString(t.Status)
}

type StatusStatus string

const (
	StatusStatusInvalid      StatusStatus = "invalid"
	StatusStatusDisconnected StatusStatus = "disconnected"
	StatusStatusConnecting   StatusStatus = "connecting"
	StatusStatusConnected    StatusStatus = "connected"
	StatusStatusGotIP        StatusStatus = "got ip"
)

func StatusStatusFromString(s string) StatusStatus {

	switch strings.ToLower(s) {

	case string(StatusStatusDisconnected):
		return StatusStatusDisconnected

	case string(StatusStatusConnecting):
		return StatusStatusConnecting

	case string(StatusStatusConnected):
		return StatusStatusConnected

	case string(StatusStatusGotIP):
		return StatusStatusGotIP

	}

	return StatusStatusInvalid
}

// Config configuration of the WiFi component contains information about the access point of the device,
// the network stations and the roaming settings.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type Config struct {
	// Ap Information about the access point
	Ap *APConfig `json:"ap,omitempty" yaml:"ap,omitempty"`
	// Sta information about the sta configuration
	Sta *STAConfig `json:"sta,omitempty" yaml:"sta,omitempty"`
	// Sta1 information about the sta configuration
	Sta1 *STAConfig `json:"sta1,omitempty" yaml:"sta1,omitempty"`
	// Roam WiFi roaming configuration
	Roam *WifiRoamConfig `json:"roam,omitempty" yaml:"roam,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatize config
func (t *Config) Sanatize() {

	if t == nil {
		return
	}

	t.Ap.Sanatize()
	t.Sta.Sanatize()
	t.Sta1.Sanatize()
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

	if !t.Ap.Equals(x.Ap) {
		zap.L().Info("Config Ap not equal")
		result = false
	}

	if !t.Sta.Equals(x.Sta) {
		zap.L().Info("Config Sta not equal")
		result = false
	}

	if !t.Sta1.Equals(x.Sta1) {
		zap.L().Info("Config Sta1 not equal")
		result = false
	}

	if !t.Roam.Equals(x.Roam) {
		zap.L().Info("Config Roam not equal")
		result = false
	}

	return result
}

func (t *Config) Merge(x *Config) {

	if x == nil {
		return
	}

	if t.Ap == nil {
		if x.Ap != nil {
			t.Ap = x.Ap.Clone()
		} else {
			t.Ap.Merge(x.Ap)
		}
	}

	if t.Sta == nil {
		if x.Sta != nil {
			t.Sta = x.Sta.Clone()
		} else {
			t.Sta.Merge(x.Sta)
		}
	}

	if t.Sta1 == nil {
		if x.Sta1 != nil {
			t.Sta1 = x.Sta1.Clone()
		} else {
			t.Sta1.Merge(x.Sta1)
		}
	}

	if t.Roam == nil {
		if x.Roam != nil {
			t.Roam = x.Roam.Clone()
		} else {
			t.Roam.Merge(x.Roam)
		}
	}

}

// APConfig WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type APConfig struct {
	// SSID readonly SSID of the access point
	SSID *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	// Pass password for the ssid, writeonly. Must be provided if you provide ssid
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
	// IsOpen True if the access point is open, false otherwise
	IsOpen *bool `json:"is_open,omitempty" yaml:"is_open,omitempty"`
	// Enable true if the access point is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
	// RangeExtender range extender configuration object, available only when range extender functionality is present.
	RangeExtender *RangeExtenderConfig `json:"range_extender,omitempty" yaml:"range_extender,omitempty"`
}

// Clone return copy
func (t *APConfig) Clone() *APConfig {
	c := &APConfig{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *APConfig) Equals(x *APConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("APConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("APConfig receiver is not nil but input is")
		return false
	}

	result := true

	if !util.CompareString(t.SSID, x.SSID) {
		zap.L().Info("APConfig SSID not equal")
		result = false
	}

	if !util.CompareBool(t.IsOpen, x.IsOpen) {
		zap.L().Info("APConfig IsOpen not equal")
		result = false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("APConfig Enable not equal")
		result = false
	}

	if !t.RangeExtender.Equals(x.RangeExtender) {
		zap.L().Info("APConfig RangeExtender not equal")
		result = false
	}

	return result
}

func (t *APConfig) Merge(x *APConfig) {

	if x == nil {
		return
	}

	if t.SSID == nil {
		t.SSID = x.SSID
	}

	if t.Pass == nil {
		t.Pass = x.Pass
	}

	if t.IsOpen == nil {
		t.IsOpen = x.IsOpen
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

	if t.RangeExtender == nil {
		if x.RangeExtender != nil {
			t.RangeExtender = x.RangeExtender.Clone()
		} else {
			t.RangeExtender.Merge(x.RangeExtender)
		}
	}

}

// Sanatize sanatize config
func (t *APConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.SSID = nil
		t.Pass = nil
		t.IsOpen = nil
		t.RangeExtender = nil
	}

	t.RangeExtender.Sanatize()
}

// RangeExtenderConfig Range extender configuration object, available only when range extender functionality is present.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type RangeExtenderConfig struct {
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// Clone return copy
func (t *RangeExtenderConfig) Clone() *RangeExtenderConfig {
	c := &RangeExtenderConfig{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *RangeExtenderConfig) Equals(x *RangeExtenderConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("RangeExtenderConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("RangeExtenderConfig receiver is not nil but input is")
		return false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("RangeExtenderConfig Enable not equal")
		return false
	}

	return true
}

func (t *RangeExtenderConfig) Merge(x *RangeExtenderConfig) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}
}

// Sanatize sanatize config
func (t *RangeExtenderConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}
}

// STAConfig WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type STAConfig struct {
	// SSID of the network
	SSID *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	// Password for the ssid, writeonly. Must be provided if you provide ssid
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
	// IsOpen true if the network is open, i.e. no password is set, false otherwise, readonly
	IsOpen *bool `json:"is_open,omitempty" yaml:"is_open,omitempty"`
	// Enable True if the configuration is enabled, false otherwise
	Enable *bool `json:"enable" yaml:"enable"`
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
func (t *STAConfig) Clone() *STAConfig {
	c := &STAConfig{}
	copier.Copy(&c, &t)
	return c
}

// Sanatize sanatize config
func (t *STAConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.SSID = nil
		t.Pass = nil
		t.IsOpen = nil
		t.Ipv4Mode = nil
		t.IP = nil
		t.Netmask = nil
		t.Gateway = nil
		t.Nameserver = nil
	}
}

// Equals returns true if equal
func (t *STAConfig) Equals(x *STAConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("STAConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("STAConfig receiver is not nil but input is")
		return false
	}

	result := true

	if !util.CompareString(t.SSID, x.SSID) {
		zap.L().Info("STAConfig SSID not equal")
		result = false
	}

	if !util.CompareBool(t.IsOpen, x.IsOpen) {
		zap.L().Info("STAConfig IsOpen not equal")
		result = false
	}

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("STAConfig Enable not equal")
		result = false
	}

	if !util.CompareString(t.Ipv4Mode, x.Ipv4Mode) {
		zap.L().Info("STAConfig Ipv4Mode not equal")
		result = false
	}

	if !util.CompareString(t.IP, x.IP) {
		zap.L().Info("STAConfig IP not equal")
		result = false
	}

	if !util.CompareString(t.Netmask, x.Netmask) {
		zap.L().Info("STAConfig Netmask not equal")
		result = false
	}

	if !util.CompareString(t.Gateway, x.Gateway) {
		zap.L().Info("STAConfig Gateway not equal")
		result = false
	}

	if !util.CompareString(t.Nameserver, x.Nameserver) {
		zap.L().Info("STAConfig Nameserver not equal")
		result = false
	}

	return result
}

func (t *STAConfig) Merge(x *STAConfig) {

	if x == nil {
		return
	}

	if t.SSID == nil {
		t.SSID = x.SSID
	}

	if t.Pass == nil {
		t.Pass = x.Pass
	}

	if t.IsOpen == nil {
		t.IsOpen = x.IsOpen
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

// WifiRoamConfig WiFi roaming configuration
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi#configuration
type WifiRoamConfig struct {
	// RSSIThreshold - when reached will trigger the access point roaming. Default value: -80
	RSSIThreshold *int `json:"rssi_thr,omitempty" yaml:"rssi_thr,omitempty"`
	// Interval at which to scan for better access points. Enabled if set to positive number,
	// disabled if set to 0. Default value: 60
	Interval *int `json:"interval,omitempty" yaml:"interval,omitempty"`
}

// Clone return copy
func (t *WifiRoamConfig) Clone() *WifiRoamConfig {
	c := &WifiRoamConfig{}
	copier.Copy(&c, &t)
	return c
}

// Equals returns true if equal
func (t *WifiRoamConfig) Equals(x *WifiRoamConfig) bool {

	if t == nil {
		if x == nil {
			return true
		}

		zap.L().Info("WifiRoamConfig receiver is nil but input is not")
		return false
	}

	if x == nil {
		zap.L().Info("WifiRoamConfig receiver is not nil but input is")
		return false
	}

	result := true

	if !util.CompareInt(t.RSSIThreshold, x.RSSIThreshold) {
		zap.L().Info("WifiRoamConfig RSSIThreshold not equal")
		result = false
	}

	if !util.CompareInt(t.Interval, x.Interval) {
		zap.L().Info("WifiRoamConfig Interval not equal")
		result = false
	}

	return result
}

func (t *WifiRoamConfig) Merge(x *WifiRoamConfig) {

	if x == nil {
		return
	}

	if t.RSSIThreshold == nil {
		t.RSSIThreshold = x.RSSIThreshold
	}

	if t.Interval == nil {
		t.Interval = x.Interval
	}
}

// Net WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi
type Net struct {
	SSID    *string `json:"ssid,omitempty" yaml:"ssid,omitempty"`
	BSSID   *string `json:"bssid,omitempty" yaml:"bssid,omitempty"`
	Auth    *int    `json:"auth,omitempty" yaml:"auth,omitempty"`
	Channel *int    `json:"channel,omitempty" yaml:"channel,omitempty"`
	RSSI    *int    `json:"rssi,omitempty" yaml:"rssi,omitempty"`
}

// Clone return copy
func (t *Net) Clone() *Net {
	c := &Net{}
	copier.Copy(&c, &t)
	return c
}

// APClient WiFi component object
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi
type APClient struct {
	MAC      *string `json:"mac,omitempty" yaml:"mac,omitempty"`
	IP       *string `json:"ip,omitempty" yaml:"ip,omitempty"`
	IPStatic *bool   `json:"ip_static,omitempty" yaml:"ip_static,omitempty"`
	Mport    *int    `json:"mport,omitempty" yaml:"mport,omitempty"`
	Since    *int    `json:"since,omitempty" yaml:"since,omitempty"`
}

// Clone return copy
func (t *APClient) Clone() *APClient {
	c := &APClient{}
	copier.Copy(&c, &t)
	return c
}

// WifiScanResults Wifi Scan Results
type ScanResults struct {
	Results []Net `json:"results,omitempty" yaml:"results,omitempty"`
}

// Clone return copy
func (t *ScanResults) Clone() *ScanResults {
	c := &ScanResults{}
	copier.Copy(&c, &t)
	return c
}

// WifiAPClients Wifi AP Clients
type APClients struct {
	Ts      *int       `json:"ts,omitempty" yaml:"ts,omitempty"`
	Clients []APClient `json:"ap_clients,omitempty" yaml:"ap_clients,omitempty"`
}

// Clone return copy
func (t *APClients) Clone() *APClients {
	c := &APClients{}
	copier.Copy(&c, &t)
	return c
}
