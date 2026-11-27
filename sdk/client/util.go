package client

import (
	"context"
	"fmt"

	"github.com/jodydadescott/shelly-client/sdk/devconfig"
)

const (
	commonConfig = "common"
)

// GetConfig returns the desired config for the device: the "common" config, with the config named
// after the device's app (model) layered on top, then the config named after its device ID.
// At least one of the app or device ID configs must exist.
func GetConfig(ctx context.Context, client *Client) (devconfig.Config, error) {

	deviceInfo, err := client.GetDeviceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetDeviceInfo] failed with error %w", err)
	}

	app := client.GetShellyConfigByName(*deviceInfo.App)
	device := client.GetShellyConfigByName(*deviceInfo.ID)
	if app == nil && device == nil {
		return nil, fmt.Errorf("no config for deviceID %s, deviceApp %s found", *deviceInfo.ID, *deviceInfo.App)
	}

	config := client.GetShellyConfigByName(commonConfig)
	config = devconfig.Merge(config, app)
	config = devconfig.Merge(config, device)
	return config, nil
}
