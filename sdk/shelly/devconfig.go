package shelly

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jodydadescott/shelly-client/sdk/devconfig"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
)

// Call sends one RPC and returns its result.
func (t *Client) Call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
		Method: &method,
		Params: params,
	})
	if err != nil {
		return nil, getErr(method, err)
	}

	var response struct {
		Result json.RawMessage  `json:"result"`
		Error  *msg_types.Error `json:"error"`
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, getErr(method, err)
	}
	if response.Error != nil {
		return nil, getErr(method, response.Error)
	}
	return response.Result, nil
}

// GetDeviceConfig returns the running config in the device's own shape (Shelly.GetConfig),
// plus devconfig.AuthKey with {"enable": <auth enabled>}.
func (t *Client) GetDeviceConfig(ctx context.Context) (devconfig.Config, error) {
	raw, err := t.Call(ctx, Component+".GetConfig", nil)
	if err != nil {
		return nil, err
	}
	config, err := devconfig.FromJSON(raw)
	if err != nil {
		return nil, getErr(Component+".GetConfig", err)
	}

	deviceInfo, err := t.GetDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}
	config[devconfig.AuthKey] = map[string]any{"enable": deviceInfo.AuthEnabled}
	return config, nil
}

// ApplyDeviceConfig makes the device match desired. Only components with drift are set, unless
// force is true, in which case every desired component is set. Components are set in
// devconfig.ApplyOrder (auth, then wifi last), then the device reboots if any change requires it.
// Enabling auth uses password, the same password this client connects with.
func (t *Client) ApplyDeviceConfig(ctx context.Context, desired devconfig.Config, password string, force bool) (*ConfigReport, error) {
	running, err := t.GetDeviceConfig(ctx)
	if err != nil {
		return nil, err
	}

	result, err := devconfig.Compare(running, desired)
	if err != nil {
		return nil, err
	}

	keys := result.Changed
	if force {
		keys = nil
		for k := range desired {
			keys = append(keys, k)
		}
	} else if result.Equal() {
		return &ConfigReport{NoChange: true}, nil
	}

	methods, err := t.ListMethods(ctx)
	if err != nil {
		return nil, err
	}

	rebootRequired := false
	for _, key := range devconfig.ApplyOrder(keys) {
		if key == devconfig.AuthKey {
			if err := t.SetAuth(ctx, desired[key], password); err != nil {
				return nil, err
			}
			continue
		}
		method, params, err := devconfig.SetConfigCall(key, desired[key], methods.Methods)
		if err != nil {
			return nil, err
		}
		raw, err := t.Call(ctx, method, params)
		if err != nil {
			return nil, fmt.Errorf("set %s: %w", key, err)
		}
		var res struct {
			RestartRequired bool `json:"restart_required"`
		}
		if err := json.Unmarshal(raw, &res); err != nil {
			return nil, getErr(method, err)
		}
		if res.RestartRequired {
			rebootRequired = true
		}
	}

	if rebootRequired {
		if err := t.Reboot(ctx); err != nil {
			return nil, err
		}
	}

	return &ConfigReport{RebootRequired: rebootRequired, Drift: result.Drift}, nil
}

// SetAuth applies desired auth, {"enable": bool}, with Shelly.SetAuth; enabling uses password.
func (t *Client) SetAuth(ctx context.Context, desired any, password string) error {
	m, ok := desired.(map[string]any)
	if !ok {
		return fmt.Errorf("%s must be an object like {enable: true}", devconfig.AuthKey)
	}
	enable, ok := m["enable"].(bool)
	if !ok {
		return fmt.Errorf("%s.enable must be true or false", devconfig.AuthKey)
	}

	deviceInfo, err := t.GetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	params := map[string]any{"user": ShellyUser, "realm": *deviceInfo.ID, "ha1": nil}
	if enable {
		if password == "" {
			return fmt.Errorf("enabling %s needs the connection password to be set", devconfig.AuthKey)
		}
		sum := sha256.Sum256([]byte(ShellyUser + ":" + *deviceInfo.ID + ":" + password))
		params["ha1"] = hex.EncodeToString(sum[:])
	}

	_, err = t.Call(ctx, Component+".SetAuth", params)
	return err
}
