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
	Config    *Config `json:"config,omitempty"`
	Connected *bool   `json:"connected,omitempty"`
}

// Status MQTT component top level status
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Mqtt
type Status struct {
	Connected *bool `json:"connected,omitempty" yaml:"connected,omitempty"`
}

// Clone return copy
func (t *Status) Clone() *Status {
	c := &Status{}
	copier.Copy(&c, &t)
	return c
}

// Config configuration of the MQTT component contains information about the credentials and prefix used and the
// protection and notifications settings of the MQTT connection.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Mqtt#configuration
type Config struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`
	// Server host name of the MQTT server. Can be followed by port number - host:port
	Server *string `json:"server,omitempty" yaml:"server,omitempty"`
	// ClientID identifies each MQTT client that connects to an MQTT brokers
	ClientID *string `json:"client_id,omitempty" yaml:"client_id,omitempty"`
	// User username for the MQTT Server
	User *string `json:"user,omitempty" yaml:"user,omitempty"`
	// Pass password for the MQTT Server
	Pass *string `json:"pass,omitempty" yaml:"pass,omitempty"`
	// SslCa type of the TCP sockets:
	// null : Plain TCP connection
	// user_ca.pem : TLS connection verified by the user-provided CA
	// ca.pem : TLS connection verified by the built-in CA bundle
	SslCa *string `json:"ssl_ca,omitempty" yaml:"ssl_ca,omitempty"`
	// TopicPrefix prefix of the topics on which device publish/subscribe. Limited to 300 characters.
	// Could not start with $ and #, +, %, ? are not allowed.
	// Values
	// null : Device id is used as topic prefix
	TopicPrefix *string `json:"topic_prefix,omitempty" yaml:"topic_prefix,omitempty"`
	// RPCNtf enables RPC notifications (NotifyStatus and NotifyEvent) to be published on <device_id|topic_prefix>/events/rpc
	// (<topic_prefix> when a custom prefix is set, <device_id> otherwise). Default value: true.
	RPCNtf *bool `json:"rpc_ntf,omitempty" yaml:"rpc_ntf,omitempty"`
	// StatusNtf enables publishing the complete component status on <device_id|topic_prefix>/status/<component>:<id> (<topic_prefix>
	// when a custom prefix is set, <device_id> otherwise). The complete status will be published if a signifficant change occurred.
	// Default value: false
	StatusNtf *bool `json:"status_ntf,omitempty" yaml:"status_ntf,omitempty"`
	// UseClientCert enable or diable usage of client certifactes to use MQTT with encription, default: false
	UseClientCert *bool `json:"use_client_cert,omitempty" yaml:"use_client_cert,omitempty"`
	// EnableRPC enable RPC
	EnableRPC *bool `json:"enable_rpc,omitempty" yaml:"enable_rpc,omitempty"`
	// EnableControl enable the MQTT control feature. Defalut value: true
	EnableControl *bool `json:"enable_control,omitempty" yaml:"enable_control,omitempty"`
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

	if t.Enable == nil {
		tmp := false
		t.Enable = &tmp
	}

	if !*t.Enable {
		t.Server = nil
		t.ClientID = nil
		t.User = nil
		t.Pass = nil
		t.SslCa = nil
		t.TopicPrefix = nil
		t.RPCNtf = nil
		t.StatusNtf = nil
		t.UseClientCert = nil
		t.EnableRPC = nil
		t.EnableControl = nil
	}
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

	if !util.CompareBool(t.Enable, x.Enable) {
		zap.L().Info("Config Enable not equal")
		result = false
	}

	if !util.CompareString(t.Server, x.Server) {
		zap.L().Info("Config Server not equal")
		result = false
	}

	if !util.CompareString(t.SslCa, x.SslCa) {
		zap.L().Info("Config SslCa not equal")
		result = false
	}

	if !util.CompareBool(t.RPCNtf, x.RPCNtf) {
		zap.L().Info("Config RPCNtf not equal")
		result = false
	}

	if !util.CompareBool(t.StatusNtf, x.StatusNtf) {
		zap.L().Info("Config StatusNtf not equal")
		result = false
	}

	if !util.CompareBool(t.UseClientCert, x.UseClientCert) {
		zap.L().Info("Config UseClientCert not equal")
		result = false
	}

	if !util.CompareBool(t.EnableRPC, x.EnableRPC) {
		zap.L().Info("Config EnableRPC not equal")
		result = false
	}

	if !util.CompareBool(t.EnableControl, x.EnableControl) {
		zap.L().Info("Config EnableControl not equal")
		result = false
	}

	return result
}

func (t *Config) Merge(x *Config) {

	if x == nil {
		return
	}

	if t.Enable == nil {
		t.Enable = x.Enable
	}

	if t.Server == nil {
		t.Server = x.Server
	}

	if t.ClientID == nil {
		t.ClientID = x.ClientID
	}

	if t.User == nil {
		t.User = x.User
	}

	if t.Pass == nil {
		t.Pass = x.Pass
	}

	if t.SslCa == nil {
		t.SslCa = x.SslCa
	}

	if t.TopicPrefix == nil {
		t.TopicPrefix = x.TopicPrefix
	}

	if t.RPCNtf == nil {
		t.RPCNtf = x.RPCNtf
	}

	if t.StatusNtf == nil {
		t.StatusNtf = x.StatusNtf
	}

	if t.UseClientCert == nil {
		t.UseClientCert = x.UseClientCert
	}

	if t.EnableRPC == nil {
		t.EnableRPC = x.EnableRPC
	}

	if t.EnableControl == nil {
		t.EnableControl = x.EnableControl
	}

}
