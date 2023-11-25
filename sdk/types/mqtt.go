package types

import (
	"github.com/jinzhu/copier"
)

// MqttStatus MQTT component top level status
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Mqtt
type MqttStatus struct {
	Connected bool `json:"connected" yaml:"connected"`
}

// Clone return copy
func (t *MqttStatus) Clone() *MqttStatus {
	c := &MqttStatus{}
	copier.Copy(&c, &t)
	return c
}

// MqttConfig configuration of the MQTT component contains information about the credentials and prefix used and the
// protection and notifications settings of the MQTT connection.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Mqtt#configuration
type MqttConfig struct {
	// Enable true if MQTT connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
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
func (t *MqttConfig) Clone() *MqttConfig {
	c := &MqttConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *MqttConfig) Markup() {

	if t == nil {
		return
	}

	t.ClientID = nil

	if t.Enable {
		if t.User != nil {
			tmp := genericPassword
			t.Pass = &tmp
		}

		return
	}

	t.Server = nil
	t.Pass = nil
	t.SslCa = nil
	t.TopicPrefix = nil
	t.RPCNtf = nil
	t.StatusNtf = nil
	t.UseClientCert = nil
	t.EnableRPC = nil
	t.EnableControl = nil
}

// Sanatize sanatize config
func (t *MqttConfig) Sanatize() {

	if t == nil {
		return
	}

	if t.User == nil {
		t.Pass = nil
	}

}
