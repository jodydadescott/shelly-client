package client

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/sdk/bluetooth"
	"github.com/jodydadescott/shelly-client/sdk/client/types"
	"github.com/jodydadescott/shelly-client/sdk/cloud"
	"github.com/jodydadescott/shelly-client/sdk/ethernet"
	"github.com/jodydadescott/shelly-client/sdk/input"
	"github.com/jodydadescott/shelly-client/sdk/light"
	"github.com/jodydadescott/shelly-client/sdk/mqtt"
	"github.com/jodydadescott/shelly-client/sdk/msghandlers"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/shelly"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
	"github.com/jodydadescott/shelly-client/sdk/switchx"
	"github.com/jodydadescott/shelly-client/sdk/system"
	"github.com/jodydadescott/shelly-client/sdk/websocket"
	"github.com/jodydadescott/shelly-client/sdk/wifi"
)

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler

type Config = types.Config
type ShellyStatus = shelly_types.Status
type ConfigReport = shelly_types.ConfigReport
type ShellyRPCMethods = shelly_types.RPCMethods
type ShellyConfig = types.ShellyConfig
type ShelllyDeviceInfo = shelly_types.DeviceInfo
type ShellyUpdateConfig = shelly_types.UpdateConfig
type UpdatesReport = shelly_types.UpdatesReport

type Client struct {
	_system    *system.Client
	shelly     *shelly.Client
	_wifi      *wifi.Client
	_bluetooth *bluetooth.Client
	_mqtt      *mqtt.Client
	_cloud     *cloud.Client
	_switch    *switchx.Client
	_light     *light.Client
	_input     *input.Client
	_websocket *websocket.Client
	_ethernet  *ethernet.Client
	MessageHandlerFactory
	config *Config
}

func New(config *Config) *Client {
	t := &Client{
		config: config,
		MessageHandlerFactory: msghandlers.NewWS(&msghandlers.Config{
			Hostname:    config.Hostname,
			Password:    config.Password,
			Username:    config.Username,
			SendTimeout: config.SendTimeout,
		}),
	}

	t.shelly = shelly.New(t)
	return t
}

func (t *Client) GetShellyConfigByName(name string) *ShellyConfig {

	if t.config.ShellyConfigs == nil {
		return nil
	}

	config := t.config.ShellyConfigs[name]
	if config != nil {
		config = config.Clone()
	}

	return config
}

func (t *Client) System() *system.Client {
	if t._system == nil {
		t._system = system.New(t)
	}
	return t._system
}

func (t *Client) Bluetooth() *bluetooth.Client {
	if t._bluetooth == nil {
		t._bluetooth = bluetooth.New(t)
	}
	return t._bluetooth
}

func (t *Client) Mqtt() *mqtt.Client {
	if t._mqtt == nil {
		t._mqtt = mqtt.New(t)
	}
	return t._mqtt
}

func (t *Client) Ethernet() *ethernet.Client {
	if t._ethernet == nil {
		t._ethernet = ethernet.New(t)
	}
	return t._ethernet
}

func (t *Client) Wifi() *wifi.Client {
	if t._wifi == nil {
		t._wifi = wifi.New(t)
	}
	return t._wifi
}

func (t *Client) Cloud() *cloud.Client {
	if t._cloud == nil {
		t._cloud = cloud.New(t)
	}
	return t._cloud
}

func (t *Client) Switch() *switchx.Client {
	if t._switch == nil {
		t._switch = switchx.New(t)
	}
	return t._switch
}

func (t *Client) Light() *light.Client {
	if t._light == nil {
		t._light = light.New(t)
	}
	return t._light
}

func (t *Client) Input() *input.Client {
	if t._input == nil {
		t._input = input.New(t)
	}
	return t._input
}

func (t *Client) Websocket() *websocket.Client {
	if t._websocket == nil {
		t._websocket = websocket.New(t)
	}
	return t._websocket
}

func (t *Client) Close() {
	zap.L().Debug("(*Client) Close()")
	t.MessageHandlerFactory.Close()
}

// GetStatus returns the status of all the components of the device.
func (t *Client) GetStatus(ctx context.Context) (*ShellyStatus, error) {
	return t.shelly.GetStatus(ctx)
}

// ListMethods lists all available RPC methods. It takes into account both ACL and authentication restrictions
// and only lists the methods allowed for the particular user/channel that's making the request.
func (t *Client) ListMethods(ctx context.Context) (*ShellyRPCMethods, error) {
	return t.shelly.ListMethods(ctx)
}

// GetConfig returns the configuration of all the components of the device.
func (t *Client) GetConfig(ctx context.Context, forceRefresh bool) (*ShellyConfig, error) {
	return t.shelly.GetConfig(ctx, forceRefresh)
}

// SetConfig sets the configuration for each component with non nil config. Note that this function
// calls into each componenet as necessary.
func (t *Client) SetConfig(ctx context.Context, config *ShellyConfig, force bool) (*ConfigReport, error) {
	return t.shelly.SetConfig(ctx, config, force)
}

// GetDeviceInfo returns information about the device.
func (t *Client) GetDeviceInfo(ctx context.Context) (*ShelllyDeviceInfo, error) {
	return t.shelly.GetDeviceInfo(ctx)
}

// CheckForUpdate checks for new firmware version for the device and returns information about it.
// If no update is available returns empty JSON object as result.
func (t *Client) CheckForUpdate(ctx context.Context) (*UpdatesReport, error) {
	return t.shelly.CheckForUpdate(ctx)
}

// Update updates the firmware version of the device.
func (t *Client) Update(ctx context.Context, config *ShellyUpdateConfig) error {
	return t.shelly.Update(ctx, config)
}

// FactoryReset resets the configuration to its default state
func (t *Client) FactoryReset(ctx context.Context) error {
	return t.shelly.FactoryReset(ctx)
}

// ResetWiFiConfig resets the WiFi configuration of the device
func (t *Client) ResetWiFiConfig(ctx context.Context) error {
	return t.shelly.ResetWiFiConfig(ctx)
}

// Reboot reboots the device
func (t *Client) Reboot(ctx context.Context) error {
	return t.shelly.Reboot(ctx)
}

func ExampleConfig() *Config {
	return &Config{
		Password:      "my password",
		SendTimeout:   time.Second * 30,
		ShellyConfigs: make(map[string]*types.ShellyConfig),
	}
}
