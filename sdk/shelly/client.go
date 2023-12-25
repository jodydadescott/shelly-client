package shelly

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
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
	GetShellyConfigByName(name string) *Config
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

// SetConfig sets the configuration for each component with non nil config. Note that this function
// calls into each componenet as necessary.
func (t *Client) SetConfig(ctx context.Context, config *Config, force bool) (*ConfigReport, error) {

	config.Sanatize()

	deviceInfo, err := t.GetDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	if deviceInfo.ID == nil {
		return nil, fmt.Errorf("deviceInfo.ID is nil")
	}

	if deviceInfo.App == nil {
		return nil, fmt.Errorf("deviceInfo.App is nil")
	}

	existingConfig, err := t.GetConfig(ctx, false)
	if err != nil {
		return nil, err
	}

	existingConfig.Sanatize()

	if force {
		zap.L().Debug("force is enabled")
	} else {
		zap.L().Debug("force is not enabled")
		if existingConfig.Equals(config) {
			zap.L().Debug("no change to config")
			return &ConfigReport{
				NoChange: true,
			}, nil
		}
	}

	rebootRequired := false

	send := func(request *Request) error {

		respBytes, err := t.getMessageHandler().Send(ctx, request)
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

	setTLSClientKey := func(config *TLSConfig) error {

		if existingConfig.TLSClientKey == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support TLSClientKey config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		method := Component + ".PutTLSClientKey"

		if config == nil {
			config = &TLSConfig{}
		}

		if config.Enable == nil {
			tmp := false
			config.Enable = &tmp
		}

		if *config.Enable {

			if config.Data == nil {
				return fmt.Errorf("missing required data")
			}

			data := splitByWidth(*config.Data, maxRPCChunkSize)
			counter := 0
			append := true

			for _, chunk := range data {

				counter++

				if len(data) <= counter {
					append = false
				}

				err := send(&Request{
					Method: &method,
					Params: &RawTLSConfig{
						Data:   &chunk,
						Append: &append,
					},
				})

				if err != nil {
					return err
				}

			}

			return nil
		}

		append := false

		return getErr(method, send(&Request{
			Method: &method,
			Params: &RawTLSConfig{
				Append: &append,
			},
		}))
	}

	setUserCA := func(config *TLSConfig) error {

		if existingConfig.UserCA == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support UserCA config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		method := Component + ".PutUserCA"

		if config == nil {
			config = &TLSConfig{}
		}

		if config.Enable == nil {
			tmp := false
			config.Enable = &tmp
		}

		if *config.Enable {

			if config.Data == nil {
				return fmt.Errorf("missing required data")
			}

			data := splitByWidth(*config.Data, maxRPCChunkSize)
			counter := 0
			append := true

			for _, chunk := range data {

				counter++

				if len(data) <= counter {
					append = false
				}

				err := send(&Request{
					Method: &method,
					Params: &RawTLSConfig{
						Data:   &chunk,
						Append: &append,
					},
				})

				if err != nil {
					return getErr(method, err)
				}

			}

			return nil
		}

		append := false

		return getErr(method, send(&Request{
			Method: &method,
			Params: &RawTLSConfig{
				Append: &append,
			},
		}))
	}

	setTLSClientCert := func(config *TLSConfig) error {

		if existingConfig.TLSClientCert == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support TLSClientCert config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		method := Component + ".PutTLSClientCert"

		if config == nil {
			config = &TLSConfig{}
		}

		if config.Enable == nil {
			tmp := false
			config.Enable = &tmp
		}

		if *config.Enable {

			if config.Data == nil {
				return fmt.Errorf("missing required data")
			}

			data := splitByWidth(*config.Data, maxRPCChunkSize)
			counter := 0
			append := true

			for _, chunk := range data {

				counter++

				if len(data) <= counter {
					append = false
				}

				err := send(&Request{
					Method: &method,
					Params: &RawTLSConfig{
						Data:   &chunk,
						Append: &append,
					},
				})

				if err != nil {
					return err
				}

			}

			return nil
		}

		append := false

		return getErr(method, send(&Request{
			Method: &method,
			Params: &RawTLSConfig{
				Append: &append,
			},
		}))
	}

	setAuth := func(config *AuthConfig) error {

		method := Component + ".SetAuth"

		if config == nil {
			config = &AuthConfig{}
		}

		if config.Enable == nil {
			tmp := false
			config.Enable = &tmp
		}

		rawUser := ShellyUser

		raw := &RawAuthConfig{
			User:  &rawUser,
			Realm: deviceInfo.ID,
		}

		if *config.Enable {

			zap.L().Debug("Auth is enabled")

			if config.Pass == nil {
				return fmt.Errorf("pass is required")
			}

			raw.User = &rawUser
			raw.Realm = deviceInfo.ID
			raw.Ha1 = config.Pass

			hashInput := *raw.User + ":" + *raw.Realm + ":" + *config.Pass
			h := sha256.New()
			h.Write([]byte(hashInput))
			b := h.Sum(nil)

			tmp := hex.EncodeToString(b)
			raw.Ha1 = &tmp
		} else {
			zap.L().Debug("Auth is disabled")
		}

		respBytes, err := t.getMessageHandler().Send(ctx, &Request{
			Method: &method,
			Params: raw,
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
			return response.Error
		}

		return nil
	}

	setBluetooth := func(config *BluetoothConfig) error {

		if existingConfig.Bluetooth == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support Bluetooth config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		tmp, err := t.Bluetooth().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if tmp != nil && *tmp {
			rebootRequired = true
		}
		return nil
	}

	setCloud := func(config *CloudConfig) error {

		if existingConfig.Cloud == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support Cloud config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		tmp, err := t.Cloud().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if tmp != nil && *tmp {
			rebootRequired = true
		}
		return nil
	}

	setEthernet := func(config *EthernetConfig) error {

		if existingConfig.Ethernet == nil {
			zap.L().Warn(fmt.Sprintf("deviceID %s, deviceApp %s does not support Ethernet config; ignoring", *deviceInfo.ID, *deviceInfo.App))
			return nil
		}

		tmp, err := t.Ethernet().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if tmp != nil && *tmp {
			rebootRequired = true
		}
		return nil
	}

	setMqtt := func(config *MqttConfig) error {

		tmp, err := t.Mqtt().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if tmp != nil && *tmp {
			rebootRequired = true
		}
		return nil
	}

	setSystem := func(config *SystemConfig) error {

		tmp, err := t.System().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if tmp != nil && *tmp {
			rebootRequired = true
		}
		return nil
	}

	setWebsocket := func(config *WebsocketConfig) error {

		tmp, err := t.Websocket().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if *tmp {
			rebootRequired = true
		}
		return nil
	}

	setWifi := func(config *WifiConfig) error {

		tmp, err := t.Wifi().SetConfig(ctx, config)
		if err != nil {
			return err
		}
		if *tmp {
			rebootRequired = true
		}
		return nil
	}

	setLight := func(config map[int]*LightConfig) error {

		var errors *multierror.Error

		for _, v := range config {
			zap.L().Debug(fmt.Sprintf("Setting config for light %d", *v.ID))
			err := t.Light().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}

		return errors.ErrorOrNil()
	}

	setInput := func(config map[int]*InputConfig) error {

		var errors *multierror.Error

		zap.L().Debug("Input config is present")

		for _, v := range config {
			zap.L().Debug(fmt.Sprintf("Setting config for input %d", *v.ID))
			err := t.Input().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}

		return errors.ErrorOrNil()
	}

	setSwitch := func(config map[int]*SwitchConfig) error {

		var errors *multierror.Error

		for _, v := range config {
			zap.L().Debug(fmt.Sprintf("Setting config for switch %d", *v.ID))
			err := t.Switch().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}

		return errors.ErrorOrNil()
	}

	var errors *multierror.Error

	addError := func(err error) {
		if err == nil {
			return
		}
		errors = multierror.Append(errors, err)
	}

	config = config.Clone()

	addError(setUserCA(config.UserCA))
	addError(setTLSClientCert(config.TLSClientCert))
	addError(setTLSClientKey(config.TLSClientKey))
	addError(setBluetooth(config.Bluetooth))
	addError(setCloud(config.Cloud))
	addError(setEthernet(config.Ethernet))
	addError(setMqtt(config.Mqtt))
	addError(setSystem(config.System))
	addError(setWebsocket(config.Websocket))
	addError(setWifi(config.Wifi))
	addError(setLight(config.Light))
	addError(setInput(config.Input))
	addError(setSwitch(config.Switch))
	addError(setAuth(config.Auth))

	if rebootRequired {
		zap.L().Debug("reboot is required; rebooting")
		err := t.Reboot(ctx)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	} else {
		zap.L().Debug("reboot is NOT required")
	}

	err = errors.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	return &ConfigReport{
		RebootRequired: rebootRequired,
	}, nil

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

func splitByWidth(str string, size int) []string {
	strLength := len(str)
	var splited []string
	var stop int
	for i := 0; i < strLength; i += size {
		stop = i + size
		if stop > strLength {
			stop = strLength
		}
		splited = append(splited, str[i:stop])
	}
	return splited
}
