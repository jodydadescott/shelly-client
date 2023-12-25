package types

import (
	"time"

	"github.com/jinzhu/copier"

	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

type ShellyConfig = shelly_types.Config

type Config struct {
	Hostname      string                   `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	ShellyConfigs map[string]*ShellyConfig `json:"shellyConfigs,omitempty" yaml:"shellyConfigs,omitempty"`
	SendTimeout   time.Duration            `json:"SendTimeout,omitempty" yaml:"SendTimeout,omitempty"`
	Username      string                   `json:"username,omitempty" yaml:"username,omitempty"`
	Password      string                   `json:"password,omitempty" yaml:"password,omitempty"`
	RetryWait     time.Duration            `json:"retryWait,omitempty" yaml:"retryWait,omitempty"`
	SendTrys      int                      `json:"sendTrys,omitempty" yaml:"sendTrys,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

func (t *Config) AddShellyConfig(name string, config *ShellyConfig) {
	if t.ShellyConfigs == nil {
		t.ShellyConfigs = make(map[string]*ShellyConfig)
	}
	t.ShellyConfigs[name] = config
}
