package sdk

import (
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/sdk/bluetooth"
	"github.com/jodydadescott/shelly-client/sdk/cloud"
	"github.com/jodydadescott/shelly-client/sdk/ethernet"
	"github.com/jodydadescott/shelly-client/sdk/input"
	"github.com/jodydadescott/shelly-client/sdk/light"
	"github.com/jodydadescott/shelly-client/sdk/mqtt"
	"github.com/jodydadescott/shelly-client/sdk/msghandlers"
	"github.com/jodydadescott/shelly-client/sdk/shelly"
	"github.com/jodydadescott/shelly-client/sdk/switchx"
	"github.com/jodydadescott/shelly-client/sdk/system"
	"github.com/jodydadescott/shelly-client/sdk/types"
	"github.com/jodydadescott/shelly-client/sdk/websocket"
	"github.com/jodydadescott/shelly-client/sdk/wifi"
)

type Config struct {
	Hostname    string
	Password    string
	SendTimeout time.Duration
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

type Client struct {
	_system    *system.Client
	_shelly    *shelly.Client
	_wifi      *wifi.Client
	_bluetooth *bluetooth.Client
	_mqtt      *mqtt.Client
	_cloud     *cloud.Client
	_switch    *switchx.Client
	_light     *light.Client
	_input     *input.Client
	_websocket *websocket.Client
	_ethernet  *ethernet.Client
	types.MessageHandlerFactory
}

func New(config *Config) *Client {
	return &Client{
		MessageHandlerFactory: msghandlers.NewWS(&msghandlers.Config{
			Hostname:    config.Hostname,
			Password:    config.Password,
			Username:    types.ShellyUser,
			SendTimeout: config.SendTimeout,
		}),
	}
}

func (t *Client) System() *system.Client {
	if t._system == nil {
		t._system = system.New(t)
	}
	return t._system
}

func (t *Client) Shelly() *shelly.Client {
	if t._shelly == nil {
		t._shelly = shelly.New(t)
	}
	return t._shelly
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
