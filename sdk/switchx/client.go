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

	t._messageHandler = t.NewHandle()
	return t._messageHandler
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
		return nil, err
	}

	response := &GetStatusResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	if response.Result == nil {
		return nil, fmt.Errorf("Result is missing from response")
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
		return nil, err
	}

	response := &GetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	if response.Result == nil {
		return nil, fmt.Errorf("Result is missing from response")
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
		return err
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return response.Error
	}

	if response.Result == nil {
		return fmt.Errorf("Result is missing from response")
	}

	return nil
}

func setOrToggle(respBytes []byte, err error) error {

	if err != nil {
		return err
	}

	rawResponse := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, rawResponse)

	if err != nil {
		return err
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return response.Error
	}

	return nil
}

// Set sets switch to on/off
func (t *Client) Set(ctx context.Context, id int, on *bool) error {

	method := Component + ".Set"

	return setOrToggle(t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
			On: on,
		},
	}))
}

// Toggle toggles switch. If switch is on it will be turned off. If switch is off it will be turned on.
func (t *Client) Toggle(ctx context.Context, id int) error {

	method := Component + ".Toggle"

	return setOrToggle(t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
		},
	}))
}

// Close closes messange handler
func (t *Client) Close() {
	if t._messageHandler != nil {
		t._messageHandler.Close()
	}
}
