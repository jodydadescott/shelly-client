package light

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/sdk/light/types"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler
type Request = msg_types.Request
type DeviceInfo = shelly_types.DeviceInfo

type Config = types.Config
type Status = types.Status
type GetStatusResponse = types.GetStatusResponse
type GetConfigResponse = types.GetConfigResponse
type Params = types.Params
type SetConfigResponse = types.SetConfigResponse

type clientContract interface {
	MessageHandlerFactory
	GetDeviceInfo(ctx context.Context) (*DeviceInfo, error)
}

// New returns new instance of client
func New(clientContract clientContract) *Client {
	return &Client{
		clientContract: clientContract,
	}
}

// Client the component client
type Client struct {
	clientContract
	_messageHandler MessageHandler
}

func (t *Client) getMessageHandler() MessageHandler {
	if t._messageHandler != nil {
		return t._messageHandler
	}

	t._messageHandler = t.NewHandle(Component)
	return t._messageHandler
}

func getErr(method string, id *int, err error) error {
	if err == nil {
		return nil
	}

	if id == nil {
		return fmt.Errorf("component %s, method %s, error %w", Component, method, err)
	}

	return fmt.Errorf("component %s, method %s, id %d, error %w", Component, method, *id, err)
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
		return nil, getErr(method, &id, err)
	}

	response := &GetStatusResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, getErr(method, &id, err)
	}

	if response.Error != nil {
		return nil, getErr(method, &id, response.Error)
	}

	if response.Result == nil {
		return nil, getErr(method, &id, fmt.Errorf("result is missing from response"))
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
		return nil, getErr(method, &id, err)
	}

	response := &GetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, getErr(method, &id, err)
	}

	if response.Error != nil {
		return nil, getErr(method, &id, response.Error)
	}

	if response.Result == nil {
		return nil, getErr(method, &id, fmt.Errorf("result is missing from response"))
	}

	return response.Result, nil
}

// SetConfig applies config to device component.
func (t *Client) SetConfig(ctx context.Context, config *Config) error {

	method := Component + ".SetConfig"

	if config == nil {
		zap.L().Debug("Light config is not present and will be disabled")
		config = &Config{}
	} else {
		zap.L().Debug("Light config is present")
		config = config.Clone()
	}

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID:     *config.ID,
			Config: config,
		},
	})

	if err != nil {
		return getErr(method, nil, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, nil, err)
	}

	if response.Error != nil {
		return getErr(method, nil, response.Error)
	}

	if response.Result == nil {
		return getErr(method, nil, fmt.Errorf("result is missing from response"))
	}

	return nil
}

func (t *Client) Set(ctx context.Context, id int, on *bool, brightness *float64) error {

	method := Component + ".Set"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID:         id,
			On:         on,
			Brightness: brightness,
		},
	})

	if err != nil {
		return getErr(method, &id, err)
	}

	rawResponse := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, rawResponse)

	if err != nil {
		return getErr(method, &id, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, &id, err)
	}

	if response.Error != nil {
		return getErr(method, &id, response.Error)
	}

	return nil
}

func (t *Client) Toggle(ctx context.Context, id int) error {

	method := Component + ".Toggle"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: &Params{
			ID: id,
		},
	})

	if err != nil {
		return getErr(method, &id, err)
	}

	rawResponse := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, rawResponse)

	if err != nil {
		return getErr(method, &id, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return getErr(method, &id, err)
	}

	if response.Error != nil {
		return getErr(method, &id, response.Error)
	}

	return nil
}
