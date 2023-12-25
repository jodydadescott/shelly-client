package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

const (
	commonConfig = "common"
)

func GetConfig(ctx context.Context, client *Client) (*ShellyConfig, error) {

	deviceInfo, err := client.GetDeviceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetDeviceInfo] failed with error %w", err)
	}

	commonConfig := client.GetShellyConfigByName(commonConfig)
	if commonConfig != nil {
		zap.L().Debug("common config exist")
	}

	config := client.GetShellyConfigByName(*deviceInfo.ID)
	if config != nil {
		zap.L().Debug(fmt.Sprintf("retrieved config with deviceID %s", *deviceInfo.ID))
		return config.Merge(commonConfig), nil
	}

	config = client.GetShellyConfigByName(*deviceInfo.App)
	if config != nil {
		zap.L().Debug(fmt.Sprintf("retrieved config with deviceApp %s", *deviceInfo.App))
		return config.Merge(commonConfig), nil
	}

	return nil, fmt.Errorf("no config for deviceID %s, deviceApp %s found", *deviceInfo.ID, *deviceInfo.App)
}
