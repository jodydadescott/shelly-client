package shelly

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/hashicorp/go-multierror"
	"github.com/jodydadescott/shelly-client/sdk/bluetooth"
	"github.com/jodydadescott/shelly-client/sdk/cloud"
	"github.com/jodydadescott/shelly-client/sdk/ethernet"
	"github.com/jodydadescott/shelly-client/sdk/input"
	"github.com/jodydadescott/shelly-client/sdk/light"
	"github.com/jodydadescott/shelly-client/sdk/mqtt"
	"github.com/jodydadescott/shelly-client/sdk/switchx"
	"github.com/jodydadescott/shelly-client/sdk/system"
	"github.com/jodydadescott/shelly-client/sdk/websocket"
	"github.com/jodydadescott/shelly-client/sdk/wifi"
)

type clientContract interface {
	MessageHandlerFactory
	System() *system.Client
	Bluetooth() *bluetooth.Client
	Mqtt() *mqtt.Client
	Wifi() *wifi.Client
	Cloud() *cloud.Client
	Switch() *switchx.Client
	Input() *input.Client
	Light() *light.Client
	Websocket() *websocket.Client
	Ethernet() *ethernet.Client
}

func New(clientContract clientContract) *Client {
	return &Client{
		clientContract: clientContract,
	}
}

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

func getErr(method string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("component %s, method %s, error %w", Component, method, err)
}

// GetStatus returns the status of all the components of the device.
func (t *Client) GetStatus(ctx context.Context) (*ShellyStatus, error) {

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

	return response.Result.convert(), nil
}

// ListMethods lists all available RPC methods. It takes into account both ACL and authentication restrictions
// and only lists the methods allowed for the particular user/channel that's making the request.
func (t *Client) ListMethods(ctx context.Context) (*ShellyRPCMethods, error) {

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
func (t *Client) GetConfig(ctx context.Context) (*ShellyConfig, error) {

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
		return nil, fmt.Errorf("Result is missing from response")
	}

	config := response.Result.convert()

	config.Auth = &ShellyAuthConfig{}

	if t.clientContract.IsAuthEnabled() {
		config.Auth.Enable = true
	}

	config.Markup()

	return config, nil
}

// SetConfig sets the configuration for each component with non nil config. Note that this function
// calls into each componenet as necessary.
func (t *Client) SetConfig(ctx context.Context, config *ShellyConfig, deviceInfo *ShelllyDeviceInfo) error {

	config = config.Clone()
	config.Sanatize()

	rebootRequired := false

	if deviceInfo == nil {
		tmp, err := t.GetDeviceInfo(ctx)
		if err != nil {
			return err
		}
		deviceInfo = tmp
	}

	var errors *multierror.Error

	if config.UserCA != nil {
		err := t.setUserCA(ctx, config.UserCA)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if config.TLSClientCert != nil {
		err := t.setTLSClientCert(ctx, config.TLSClientCert)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if config.TLSClientKey != nil {
		err := t.setTLSClientKey(ctx, config.TLSClientKey)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if config.Bluetooth != nil {
		tmpRebootRequired, err := t.Bluetooth().SetConfig(ctx, config.Bluetooth)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Cloud != nil {
		tmpRebootRequired, err := t.Cloud().SetConfig(ctx, config.Cloud)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Ethernet != nil {
		tmpRebootRequired, err := t.Ethernet().SetConfig(ctx, config.Ethernet)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Mqtt != nil {
		if config.Mqtt.ClientID == nil {
			config.Mqtt.ClientID = deviceInfo.ID
			zap.L().Debug(fmt.Sprintf("Mqtt.ClientID set set to %s", *config.Mqtt.ClientID))
		} else if config.Mqtt.ClientID != deviceInfo.ID {
			zap.L().Info(fmt.Sprintf("Mqtt.ClientID is set to %s; is this what you really want?", *config.Mqtt.ClientID))
		}

		tmpRebootRequired, err := t.Mqtt().SetConfig(ctx, config.Mqtt)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.System != nil {
		tmpRebootRequired, err := t.System().SetConfig(ctx, config.System)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Websocket != nil {
		tmpRebootRequired, err := t.Websocket().SetConfig(ctx, config.Websocket)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Wifi != nil {
		tmpRebootRequired, err := t.Wifi().SetConfig(ctx, config.Wifi)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			if tmpRebootRequired != nil && *tmpRebootRequired {
				rebootRequired = true
			}
		}
	}

	if config.Light != nil {
		for _, v := range config.Light {
			err := t.Light().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	if config.Input != nil {
		for _, v := range config.Input {
			err := t.Input().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	if config.Switch != nil {
		for _, v := range config.Switch {
			err := t.Switch().SetConfig(ctx, v)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	if config.Auth != nil {
		err := t.setAuth(ctx, config.Auth, deviceInfo)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if rebootRequired {
		zap.L().Debug("reboot is required; rebooting")
		err := t.Reboot(ctx)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	} else {
		zap.L().Debug("reboot is NOT required")
	}

	return errors.ErrorOrNil()
}

// GetDeviceInfo returns information about the device.
func (t *Client) GetDeviceInfo(ctx context.Context) (*ShelllyDeviceInfo, error) {

	method := Component + ".GetDeviceInfo"

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

	return response.Result, nil
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
func (t *Client) Update(ctx context.Context, config *ShellyUpdateConfig) error {

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

func (t *Client) setAuth(ctx context.Context, config *ShellyAuthConfig, deviceInfo *ShelllyDeviceInfo) error {

	method := Component + ".SetAuth"

	rawUser := ShellyUser

	raw := &RawShellyAuthConfig{
		User:  &rawUser,
		Realm: deviceInfo.ID,
	}

	if config.Enable {

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

func (t *Client) send(ctx context.Context, request *Request) error {

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

func (t *Client) setUserCA(ctx context.Context, config *ShellyUserCAConfig) error {

	method := Component + ".PutUserCA"

	if config.Enable {

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

			err := t.send(ctx, &Request{
				Method: &method,
				Params: &RawShellyTLSConfig{
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

	return getErr(method, t.send(ctx, &Request{
		Method: &method,
		Params: &RawShellyTLSConfig{
			Append: &append,
		},
	}))
}

func (t *Client) setTLSClientCert(ctx context.Context, config *ShellyTLSClientCertConfig) error {

	method := Component + ".PutTLSClientCert"

	if config.Enable {

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

			err := t.send(ctx, &Request{
				Method: &method,
				Params: &RawShellyTLSConfig{
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

	return getErr(method, t.send(ctx, &Request{
		Method: &method,
		Params: &RawShellyTLSConfig{
			Append: &append,
		},
	}))
}

func (t *Client) setTLSClientKey(ctx context.Context, config *ShellyTLSClientKeyConfig) error {

	method := Component + ".PutTLSClientKey"

	if config.Enable {

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

			err := t.send(ctx, &Request{
				Method: &method,
				Params: &RawShellyTLSConfig{
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

	return getErr(method, t.send(ctx, &Request{
		Method: &method,
		Params: &RawShellyTLSConfig{
			Append: &append,
		},
	}))
}
