package types

import (
	"strings"
	"time"

	"github.com/jinzhu/copier"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/jodydadescott/unifi-go-sdk"

	"github.com/jodydadescott/shelly-client/sdk/types"
)

type Logger = logger.Config
type ShellyConfig = types.ShellyConfig
type UnifiConfig = unifi.Config

type Config struct {
	Notes         string                   `json:"notes,omitempty" yaml:"notes,omitempty"`
	Hostnames     []string                 `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	Password      *string                  `json:"password,omitempty" yaml:"password,omitempty"`
	Output        *string                  `json:"output,omitempty" yaml:"output,omitempty"`
	UpdateURL     *string                  `json:"updateURL,omitempty" yaml:"updateURL,omitempty"`
	Timeout       *time.Duration           `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger        *Logger                  `json:"logger,omitempty" yaml:"logger,omitempty"`
	Unifi         *UnifiConfig             `json:"unifi,omitempty" yaml:"unifi,omitempty"`
	ShellyConfigs map[string]*ShellyConfig `json:"shellyConfigs,omitempty" yaml:"shellyConfigs,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

func (t *Config) AddHostname(hostname string) {
	if hostname == "" {
		return
	}

	for _, v := range t.Hostnames {
		if v == hostname {
			return
		}
	}
	t.Hostnames = append(t.Hostnames, hostname)
}

func (t *Config) AddShellyConfig(name string, config *ShellyConfig) {
	if t.ShellyConfigs == nil {
		t.ShellyConfigs = make(map[string]*types.ShellyConfig)
	}
	t.ShellyConfigs[name] = config
}

func (t *Config) GetMatchingShellyConfig(match string) *ShellyConfig {
	if t.ShellyConfigs == nil {
		t.ShellyConfigs = make(map[string]*types.ShellyConfig)
		return nil
	}

	match = strings.ToLower(match)

	for name, config := range t.ShellyConfigs {
		if strings.HasPrefix(strings.ToLower(name), match) {
			return config
		}
	}

	return nil
}

func ExampleConfig() *Config {

	var hostnames []string
	hostnames = append(hostnames, "host1")
	hostnames = append(hostnames, "host2")
	shellyPassword := "not really the secret"
	output := "pretty-json"
	timeout := time.Second * 60

	notes := "Depending on the command Hostname may be required. For commands that\n"
	notes += "use Hostnames hostname will be prepended if it is set. If Unifi\n"
	notes += "config is present then hostnames will be loaded from Unifi. Shelly\n"
	notes += "configs can be specified in the map. The name should be the device\n"
	notes += "ID or App. Device ID takes precedence\n"

	unifiConfig := unifi.ExampleConfig()
	unifiConfig.Enabled = true

	return &Config{
		Notes:     notes,
		Hostnames: hostnames,
		Unifi:     unifiConfig,
		Password:  &shellyPassword,
		Output:    &output,
		Timeout:   &timeout,
		Logger: &Logger{
			LogLevel: logger.DebugLevel,
		},
		ShellyConfigs: make(map[string]*types.ShellyConfig),
	}
}
