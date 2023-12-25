package wifi

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/wifi/types"
)

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler
type Request = msg_types.Request

type Config = types.Config
type Status = types.Status
type ScanResponse = types.ScanResponse
type ScanResults = types.ScanResults
type APClients = types.APClients
type GetStatusResponse = types.GetStatusResponse
type GetConfigResponse = types.GetConfigResponse
type Params = types.Params
type SetConfigResponse = types.SetConfigResponse
type ListAPClientsResponse = types.ListAPClientsResponse

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
		zap.L().Debug("Wifi config is not present and will be disabled")
		config = &Config{}
	} else {
		zap.L().Debug("Wifi config is present")
		config = config.Clone()
	}

	if config.Ap != nil {
		if config.Ap.Enable == nil {
			config.Ap.Enable = &falsex
		}
		if !*config.Ap.Enable {
			config.Ap.SSID = nil
			config.Ap.Pass = nil
			config.Ap.IsOpen = nil
			config.Ap.RangeExtender = nil
		}
	}

	if config.Sta != nil {
		if config.Sta.Enable == nil {
			config.Sta.Enable = &falsex
		}
		if !*config.Sta.Enable {
			config.Sta.SSID = nil
			config.Sta.Pass = nil
			config.Sta.IsOpen = nil
			config.Sta.Ipv4Mode = nil
			config.Sta.IP = nil
			config.Sta.Netmask = nil
			config.Sta.Gateway = nil
			config.Sta.Nameserver = nil
		}
	}

	if config.Sta1 != nil {
		if config.Sta1.Enable == nil {
			config.Sta1.Enable = &falsex
		}
		if !*config.Sta.Enable {
			config.Sta1.SSID = nil
			config.Sta1.Pass = nil
			config.Sta1.IsOpen = nil
			config.Sta1.Ipv4Mode = nil
			config.Sta1.IP = nil
			config.Sta1.Netmask = nil
			config.Sta1.Gateway = nil
			config.Sta1.Nameserver = nil
		}
	}

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
		return nil, getErr(method, err)
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

// Scan scans for Wifi networks and returns results or an error
func (t *Client) Scan(ctx context.Context) (*ScanResults, error) {

	method := Component + ".Scan"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})

	if err != nil {
		return nil, getErr(method, err)
	}

	response := &ScanResponse{}
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

// ListAPClients returns list of AP Clients or an error
func (t *Client) ListAPClients(ctx context.Context) (*APClients, error) {

	method := Component + ".ListAPClients"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})

	if err != nil {
		return nil, getErr(method, err)
	}

	response := &ListAPClientsResponse{}
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
