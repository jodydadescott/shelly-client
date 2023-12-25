package types

import (
	"time"

	"github.com/jinzhu/copier"

	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	light_types "github.com/jodydadescott/shelly-client/sdk/light/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

type LightStatus = light_types.Status

type ShellyClient = sdk_client.Client

type ShellyDeviceInfo = shelly_types.DeviceInfo
type ShellyDeviceStatus = shelly_types.Status
type ShellyConfig = sdk_types.ShellyConfig

type Config struct {
	Notes          string        `json:"notes,omitempty" yaml:"notes,omitempty"`
	Broker         string        `json:"broker,omitempty" yaml:"broker,omitempty"`
	Topic          string        `json:"topic,omitempty" yaml:"topic,omitempty"`
	Username       string        `json:"username,omitempty" yaml:"username,omitempty"`
	Password       string        `json:"password,omitempty" yaml:"password,omitempty"`
	ClientID       string        `json:"clientID,omitempty" yaml:"clientID,omitempty"`
	LightSource    string        `json:"lightSource,omitempty" yaml:"lightSource,omitempty"`
	KeepAlive      time.Duration `json:"keepAlive,omitempty" yaml:"keepAlive,omitempty"`
	PingTimeout    time.Duration `json:"pingTimeout,omitempty" yaml:"pingTimeout,omitempty"`
	PublishTimeout time.Duration `json:"publishTimeout,omitempty" yaml:"publishTimeout,omitempty"`
	DaemonInterval time.Duration `json:"daemonInterval,omitempty" yaml:"daemonInterval,omitempty"`
	FailOnError    bool          `json:"failOnError,omitempty" yaml:"failOnError,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

type MqttStatus struct {
	Src    string  `json:"src,omitempty" yaml:"src,omitempty"`
	Dst    string  `json:"dst,omitempty" yaml:"dst,omitempty"`
	Method string  `json:"method,omitempty" yaml:"method,omitempty"`
	Params *Params `json:"params,omitempty" yaml:"params,omitempty"`
}

type Params struct {
	Ts     *float64     `json:"ts,omitempty" yaml:"ts,omitempty"`
	Light0 *LightStatus `json:"light:0,omitempty" yaml:"light:0,omitempty"`
	Light1 *LightStatus `json:"light:1,omitempty" yaml:"light:1,omitempty"`
	Light2 *LightStatus `json:"light:2,omitempty" yaml:"light:2,omitempty"`
	Light3 *LightStatus `json:"light:3,omitempty" yaml:"light:3,omitempty"`
	Light4 *LightStatus `json:"light:4,omitempty" yaml:"light:4,omitempty"`
	Light5 *LightStatus `json:"light:5,omitempty" yaml:"light:5,omitempty"`
	Light6 *LightStatus `json:"light:6,omitempty" yaml:"light:6,omitempty"`
	Light7 *LightStatus `json:"light:7,omitempty" yaml:"light:7,omitempty"`
}
