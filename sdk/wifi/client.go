package wifi

import (
	"context"
	"encoding/json"
	"fmt"
)

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

	response.Result.Markup()

	return response.Result, nil
}

// SetConfig applies config to device component. Returns reboot required or error.
// If reboot is requred true is returned otherwise false.
func (t *Client) SetConfig(ctx context.Context, config *Config) (*bool, error) {

	method := Component + ".SetConfig"

	config = config.Clone()
	config.Sanatize()

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
func (t *Client) Scan(ctx context.Context) (*WifiScanResults, error) {

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
