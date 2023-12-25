package types

import (
	"time"

	"github.com/jinzhu/copier"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/jodydadescott/openhab-go-sdk"
	"github.com/jodydadescott/unifi-go-sdk"

	mqtt "github.com/jodydadescott/shelly-client/cmd/mqtt/types"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client"
)

type Logger = logger.Config
type ShellyConfig = sdk_types.Config
type UnifiConfig = unifi.Config
type MqttConfig = mqtt.Config
type OpenHAB = openhab.Config

type Config struct {
	Notes     string         `json:"notes,omitempty" yaml:"notes,omitempty"`
	Hostnames []string       `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	Output    *string        `json:"output,omitempty" yaml:"output,omitempty"`
	UpdateURL *string        `json:"updateURL,omitempty" yaml:"updateURL,omitempty"`
	Timeout   *time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger    *Logger        `json:"logger,omitempty" yaml:"logger,omitempty"`
	Unifi     *UnifiConfig   `json:"unifi,omitempty" yaml:"unifi,omitempty"`
	Shelly    *ShellyConfig  `json:"shelly,omitempty" yaml:"shelly,omitempty"`
	Mqtt      *MqttConfig    `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	OpenHAB   *OpenHAB       `json:"openHAB,omitempty" yaml:"openHAB,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}
