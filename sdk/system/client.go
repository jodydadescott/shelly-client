package system

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/system/types"
)

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler
type Request = msg_types.Request

type Config = types.Config
type Status = types.Status
type UDP = types.UDP
type GetStatusResponse = types.GetStatusResponse
type GetConfigResponse = types.GetConfigResponse
type Params = types.Params
type SetConfigResponse = types.SetConfigResponse

// New returns new instance of client
func New(messageHandlerFactory MessageHandlerFactory) *Client {
	return &Client{
		MessageHandlerFactory: messageHandlerFactory,
	}
}

// Client the component client
type Client struct {
	MessageHandlerFactory
	_messageHandler MessageHandler
}

func (t *Client) getMessageHandler() MessageHandler {
	if t._messageHandler != nil {
		return t._messageHandler
	}

	t._messageHandler = t.NewHandle(Component)
	return t._messageHandler
}

func getErr(method string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("component %s, method %s, error %w", Component, method, err)
}

// GetStatus returns status for component or error
func (t *Client) GetStatus(ctx context.Context) (*Status, error) {

	method := Component + ".GetStatus"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})

	if err != nil {
		return nil, getErr(method, err)
	}

	response := &GetStatusResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, getErr(method, response.Error)
	}

	if response.Result == nil {
		return nil, getErr(method, fmt.Errorf("result is missing from response"))
	}

	return response.Result, nil
}

// GetConfig returns component config or error
func (t *Client) GetConfig(ctx context.Context) (*Config, error) {

	method := Component + ".GetConfig"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})

	if err != nil {
		return nil, getErr(method, err)
	}

	response := &GetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, getErr(method, err)
	}

	if response.Error != nil {
		return nil, getErr(method, response.Error)
	}

	if response.Result == nil {
		return nil, getErr(method, fmt.Errorf("result is missing from response"))
	}

	return response.Result, nil
}

// SetConfig applies config to device component. Returns reboot required or error.
// If reboot is requred true is returned otherwise false.
func (t *Client) SetConfig(ctx context.Context, config *Config) (*bool, error) {

	method := Component + ".SetConfig"

	if config == nil {
		zap.L().Debug("System config is not present and will be disabled")
		config = &Config{}
	} else {
		zap.L().Debug("System config is present")
		config = config.Clone()
	}

	config.CfgRev = nil

	if config.Device != nil {
		config.Device.MAC = nil
		config.Device.FwID = nil
	}

	if config.Debug != nil {
		if config.Debug.Mqtt != nil {
			if config.Debug.Mqtt.Enable != nil {
				config.Debug.Mqtt.Enable = &falsepointer
			}
		}

		if config.Debug.Websocket != nil {
			if config.Debug.Websocket.Enable != nil {
				config.Debug.Websocket.Enable = &falsepointer
			}
		}

		if config.Debug.UDP == nil {
			config.Debug.UDP = &UDP{}
		}

		if !*config.Debug.Mqtt.Enable && !*config.Debug.Websocket.Enable {
			config.Debug.UDP.Addr = nil
		}

	}

	// Device *SystemDevice `json:"device,omitempty" yaml:"device,omitempty"`
	// // Location information about the current location of the device
	// Location *SystemLocation `json:"location,omitempty" yaml:"location,omitempty"`
	// // Debug configuration of the device's debug logs.
	// Debug *SystemDebug `json:"debug,omitempty" yaml:"debug,omitempty"`
	// // UIData user interface data
	// UIData *SystemUIData `json:"ui_data,omitempty" yaml:"ui_data,omitempty"`
	// // RPCUDP configuration for the RPC over UDP
	// RPCUDP *SystemRPCUDP `json:"rpc_udp,omitempty" yaml:"rpc_udp,omitempty"`
	// // Sntp configuration for the sntp server
	// Sntp *SystemSntp `json:"sntp,omitempty" yaml:"sntp,omitempty"`
	// // CfgRev Configuration revision. This number will be incremented for every configuration change of a device component.
	// // If the new config value is the same as the old one there will be no change of this property. Can not be modified
	// // explicitly by a call to Sys.SetConfig
	// CfgRev *int `json:"cfg_rev,omitempty" yaml:"cfg_rev,omitempty"`

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			Config: config,
		},
	})

	if err != nil {
		return nil, getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, getErr(method, response.Error)
	}

	if response.Result == nil {
		return nil, getErr(method, fmt.Errorf("result is missing from response"))
	}

	rebootRequired := false

	if response.Result.RestartRequired != nil {
		if *response.Result.RestartRequired {
			rebootRequired = true
		}
	}

	return &rebootRequired, nil
}
