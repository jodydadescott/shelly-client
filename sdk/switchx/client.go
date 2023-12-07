package switchx

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
func (t *Client) GetStatus(ctx context.Context, id int) (*Status, error) {

	method := Component + ".GetStatus"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
		},
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
		return response.Result, getErr(method, fmt.Errorf("result is missing from response"))
	}

	return response.Result, nil
}

// GetConfig returns component config or error
func (t *Client) GetConfig(ctx context.Context, id int) (*Config, error) {

	method := Component + ".GetConfig"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
		},
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
		return response.Result, getErr(method, fmt.Errorf("result is missing from response"))
	}

	response.Result.Markup()

	return response.Result, nil
}

// SetConfig applies config to device component
func (t *Client) SetConfig(ctx context.Context, config *Config) error {

	method := Component + ".SetConfig"

	config = config.Clone()
	config.Sanatize()

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID:     config.ID,
			Config: config,
		},
	})

	if err != nil {
		return getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, err)
	}

	if response.Error != nil {
		return getErr(method, response.Error)
	}

	if response.Result == nil {
		return getErr(method, fmt.Errorf("result is missing from response"))
	}

	return nil
}

// Set sets switch to on/off
func (t *Client) Set(ctx context.Context, id int, on *bool) error {

	method := Component + ".Set"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
			On: on,
		},
	})

	if err != nil {
		return getErr(method, err)
	}

	rawResponse := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, rawResponse)

	if err != nil {
		return getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, err)
	}

	if response.Error != nil {
		return getErr(method, response.Error)
	}

	return nil
}

// Toggle toggles switch. If switch is on it will be turned off. If switch is off it will be turned on.
func (t *Client) Toggle(ctx context.Context, id int) error {

	method := Component + ".Toggle"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
		},
	})

	if err != nil {
		return getErr(method, err)
	}

	rawResponse := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, rawResponse)

	if err != nil {
		return getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, err)
	}

	if response.Error != nil {
		return getErr(method, response.Error)
	}

	return nil
}
