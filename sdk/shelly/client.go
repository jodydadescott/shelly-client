package shelly

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	bluetooth_client "github.com/jodydadescott/shelly-client/sdk/bluetooth"
	bluetooth_types "github.com/jodydadescott/shelly-client/sdk/bluetooth/types"
	cloud_client "github.com/jodydadescott/shelly-client/sdk/cloud"
	cloud_types "github.com/jodydadescott/shelly-client/sdk/cloud/types"
	ethernet_client "github.com/jodydadescott/shelly-client/sdk/ethernet"
	ethernet_types "github.com/jodydadescott/shelly-client/sdk/ethernet/types"
	input_client "github.com/jodydadescott/shelly-client/sdk/input"
	input_types "github.com/jodydadescott/shelly-client/sdk/input/types"
	light_client "github.com/jodydadescott/shelly-client/sdk/light"
	light_types "github.com/jodydadescott/shelly-client/sdk/light/types"
	mqtt_client "github.com/jodydadescott/shelly-client/sdk/mqtt"
	mqtt_types "github.com/jodydadescott/shelly-client/sdk/mqtt/types"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
	switch_client "github.com/jodydadescott/shelly-client/sdk/switchx"
	switch_types "github.com/jodydadescott/shelly-client/sdk/switchx/types"
	system_client "github.com/jodydadescott/shelly-client/sdk/system"
	system_types "github.com/jodydadescott/shelly-client/sdk/system/types"
	websocket_client "github.com/jodydadescott/shelly-client/sdk/websocket"
	websocket_types "github.com/jodydadescott/shelly-client/sdk/websocket/types"
	wifi_client "github.com/jodydadescott/shelly-client/sdk/wifi"
	wifi_types "github.com/jodydadescott/shelly-client/sdk/wifi/types"
)

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler
type Request = msg_types.Request

type Config = shelly_types.Config
type ConfigReport = shelly_types.ConfigReport
type DeviceInfo = shelly_types.DeviceInfo
type ListMethodsResponse = shelly_types.ListMethodsResponse
type TLSConfig = shelly_types.TLSConfig
type Status = shelly_types.Status
type SetConfigResponse = shelly_types.SetConfigResponse
type RawTLSConfig = shelly_types.RawTLSConfig
type AuthConfig = shelly_types.AuthConfig
type RawAuthConfig = shelly_types.RawAuthConfig
type GetConfigResponse = shelly_types.GetConfigResponse
type GetStatusResponse = shelly_types.GetStatusResponse
type RPCMethods = shelly_types.RPCMethods
type DeviceInfoResponse = shelly_types.DeviceInfoResponse
type UpdatesReport = shelly_types.UpdatesReport
type UpdateConfig = shelly_types.UpdateConfig
type CheckForUpdateResponse = shelly_types.CheckForUpdateResponse
type EthernetConfig = ethernet_types.Config
type CloudConfig = cloud_types.Config
type BluetoothConfig = bluetooth_types.Config
type MqttConfig = mqtt_types.Config
type SystemConfig = system_types.Config
type WebsocketConfig = websocket_types.Config
type WifiConfig = wifi_types.Config
type LightConfig = light_types.Config
type InputConfig = input_types.Config
type SwitchConfig = switch_types.Config

type clientContract interface {
	MessageHandlerFactory
	System() *system_client.Client
	Bluetooth() *bluetooth_client.Client
	Mqtt() *mqtt_client.Client
	Wifi() *wifi_client.Client
	Cloud() *cloud_client.Client
	Switch() *switch_client.Client
	Input() *input_client.Client
	Light() *light_client.Client
	Websocket() *websocket_client.Client
	Ethernet() *ethernet_client.Client
}

func New(clientContract clientContract) *Client {
	return &Client{
		clientContract: clientContract,
	}
}

type Client struct {
	clientContract
	deviceInfo      *DeviceInfo
	shellyConfig    *Config
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

// GetStatus returns the status of all the components of the device.
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

	return response.Result.Convert(), nil
}

// ListMethods lists all available RPC methods. It takes into account both ACL and authentication restrictions
// and only lists the methods allowed for the particular user/channel that's making the request.
func (t *Client) ListMethods(ctx context.Context) (*RPCMethods, error) {

	// Do NOT validate command here because it would be recursive

	method := Component + ".ListMethods"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})
	if err != nil {
		return nil, getErr(method, err)
	}

	response := &ListMethodsResponse{}
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

// GetConfig returns the configuration of all the components of the device.
func (t *Client) GetConfig(ctx context.Context, forceRefresh bool) (*Config, error) {

	method := Component + ".GetConfig"

	if forceRefresh {
		zap.L().Debug("forceRefresh is true")
	} else {
		if t.shellyConfig == nil {
			zap.L().Debug("forceRefresh is false but there is no existing ShellyConfig")
		} else {
			zap.L().Debug("forceRefresh is false and there is an existing ShellyConfig")
			return t.shellyConfig.Clone(), nil
		}
	}

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
		return nil, fmt.Errorf("result is missing from response")
	}

	config := response.Result.Convert()

	config.Auth = &AuthConfig{}

	authEnabled := false

	if t.clientContract.IsAuthEnabled() {
		authEnabled = true
	}

	config.Auth.Enable = &authEnabled

	t.shellyConfig = config
	return config.Clone(), nil
}

// GetDeviceInfo returns information about the device.
func (t *Client) GetDeviceInfo(ctx context.Context) (*DeviceInfo, error) {

	method := Component + ".GetDeviceInfo"

	if t.deviceInfo != nil {
		return t.deviceInfo.Clone(), nil
	}

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})
	if err != nil {
		return nil, getErr(method, err)
	}

	response := &DeviceInfoResponse{}
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

	t.deviceInfo = response.Result
	return t.deviceInfo.Clone(), nil
}

// CheckForUpdate checks for new firmware version for the device and returns information about it.
// If no update is available returns empty JSON object as result.
func (t *Client) CheckForUpdate(ctx context.Context) (*UpdatesReport, error) {

	method := Component + ".CheckForUpdate"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})
	if err != nil {
		return nil, getErr(method, err)
	}

	response := &CheckForUpdateResponse{}
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

	return &UpdatesReport{
		Src:              response.Src,
		AvailableUpdates: response.Result,
	}, nil
}

// Update updates the firmware version of the device.
func (t *Client) Update(ctx context.Context, config *UpdateConfig) error {

	method := Component + ".Update"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: config,
	})
	if err != nil {
		return getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return getErr(method, response.Error)
	}

	return nil
}

// FactoryReset resets the configuration to its default state
func (t *Client) FactoryReset(ctx context.Context) error {

	method := Component + ".FactoryReset"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
	})
	if err != nil {
		return getErr(method, err)
	}

	response := &SetConfigResponse{}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return getErr(method, response.Error)
	}

	return nil
}

// ResetWiFiConfig resets the WiFi configuration of the device
func (t *Client) ResetWiFiConfig(ctx context.Context) error {

	method := Component + ".ResetWiFiConfig"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
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

	return nil
}

// Reboot reboots the device
func (t *Client) Reboot(ctx context.Context) error {

	method := Component + ".Reboot"

	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
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

	return nil
}
