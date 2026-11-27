package util

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/jodydadescott/openhab-go-sdk"
	"github.com/jodydadescott/unifi-go-sdk"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/cmd/types"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

var space = regexp.MustCompile(`\s+`)
var shellyPrefixes = []string{"shellypluswdus", "shellyplus1pm"}

type ShellyDeviceInfo = shelly_types.DeviceInfo
type ShellyDeviceStatus = shelly_types.Status
type ShellyConfig = sdk_types.Config
type UnifiConfig = unifi.Config
type ShellyClient = sdk_client.Client
type OpenHABConfig = openhab.Config

type Config = types.Config

func NewShellyClient(config *Config, hostname string) *ShellyClient {

	if config.Shelly == nil {
		panic("shelly config is required")
	}

	newShellyConfig := &ShellyConfig{
		Hostname:      hostname,
		Username:      config.Shelly.Username,
		Password:      config.Shelly.Password,
		ShellyConfigs: config.Shelly.ShellyConfigs,
	}

	return sdk_client.New(newShellyConfig)
}

func CleanupHostname(hostname string) string {
	if hostname == "" {
		return ""
	}
	split := strings.Split(hostname, "http://")
	if len(split) > 1 {
		return split[1]
	} else {
		split = strings.Split(hostname, "htts://")
		if len(split) > 1 {
			return split[1]
		}
	}
	return hostname
}

func Process(ctx context.Context, config *Config, action string, getDevStatus bool, do func(ctx context.Context, hostname string, shellyClient *ShellyClient, shellyDeviceInfo *ShellyDeviceInfo, shellyDeviceStatus *ShellyDeviceStatus) error) error {

	execute := func(ctx context.Context, hostname string, workerID int) error {

		hostname = CleanupHostname(hostname)

		client := NewShellyClient(config, hostname)
		defer client.Close()

		deviceInfo, err := client.GetDeviceInfo(ctx)
		if err != nil {
			return fmt.Errorf("workerID %d, hostname %s, [deviceInfo] failed with error %w", workerID, hostname, err)
		}

		var devStatus *ShellyDeviceStatus

		if getDevStatus {
			tmp, err := client.GetStatus(ctx)
			if err != nil {
				return fmt.Errorf("workerID %d, hostname %s, [deviceStatus] failed with error %w", workerID, hostname, err)
			}

			devStatus = tmp
		}

		err = do(ctx, hostname, client, deviceInfo, devStatus)

		if err != nil {
			return fmt.Errorf("workerID %d, hostname %s, deviceInfoID %s, deviceInfoApp %s, [%s] failed with error %w", workerID, hostname, *deviceInfo.ID, *deviceInfo.App, action, err)
		}

		return nil
	}

	getHostnames := func() ([]string, error) {

		var hostnames []string

		addHostname := func(hostname string) {
			for _, v := range hostnames {
				if v == hostname {
					return
				}
			}
			zap.L().Debug(fmt.Sprintf("Adding hostname %s from Unifi", hostname))
			hostnames = append(hostnames, hostname)
		}

		if len(config.Hostnames) > 0 {
			zap.L().Debug(fmt.Sprintf("there are %d hostnames in the config", len(config.Hostnames)))
			for _, hostname := range config.Hostnames {
				addHostname(hostname)
			}
		} else {
			zap.L().Debug("there are no hostnames in the config")
		}

		if config.Unifi == nil {
			zap.L().Debug("Unifi config is not present")
			return hostnames, nil
		}

		if config.Unifi.Enabled {
			zap.L().Debug("unifi is enabled")
		} else {
			zap.L().Debug("unifi is disabled")
			return hostnames, nil
		}

		unifiHostnames, err := GetUnifiHostnames(config.Unifi)
		if err != nil {
			return nil, err
		}

		if len(unifiHostnames) > 0 {
			zap.L().Debug(fmt.Sprintf("there are %d hostnames from unifi", len(unifiHostnames)))
		} else {
			zap.L().Debug("there are no hostnames from unifi")
		}

		hostnames = append(hostnames, unifiHostnames...)

		return hostnames, nil
	}

	hostnames, err := getHostnames()
	if err != nil {
		return err
	}

	jobCount := len(hostnames)

	if jobCount == 0 {
		return fmt.Errorf("one or more hostnames required")
	}

	// At most 4 devices at a time; every device is processed and every error is returned.
	sem := make(chan struct{}, 4)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	var errs *multierror.Error

	for i, hostname := range hostnames {
		wg.Add(1)
		go func(id int, hostname string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if err := execute(ctx, hostname, id); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}
		}(i, hostname)
	}

	wg.Wait()
	return errs.ErrorOrNil()
}

func hasShellyPrefix(input string) bool {
	input = strings.ToLower(input)
	for _, shellyPrefix := range shellyPrefixes {
		if strings.HasPrefix(input, shellyPrefix) {
			return true
		}
	}
	return false
}

func normalizeHostname(input string) string {
	input = strings.ToLower(input)
	input = space.ReplaceAllString(input, "-")
	return strings.Split(input, ".")[0]
}

func GetUnifiHostnames(config *UnifiConfig) ([]string, error) {

	var hostnames []string

	addHostname := func(hostname string) {
		for _, v := range hostnames {
			if v == hostname {
				return
			}
		}
		hostnames = append(hostnames, hostname)
	}

	unifiClient := unifi.New(config)

	clients, err := unifiClient.GetClients()
	if err != nil {
		return nil, err
	}

	for _, client := range clients {

		rawName := client.Name

		if rawName == "" {
			rawName = client.Hostname
		}

		if rawName == "" {
			rawName = client.DisplayName
		}

		hostname := normalizeHostname(rawName)

		if hasShellyPrefix(hostname) {
			addHostname(hostname)
		} else {
			if logger.Trace {
				zap.L().Debug(fmt.Sprintf("Ignoring hostname %s", hostname))
			}
		}
	}

	return hostnames, nil
}

func GetOpenHABHostnames(ctx context.Context, config *OpenHABConfig) ([]string, error) {

	var hostnames []string

	addHostname := func(hostname string) {
		if hostname == "" {
			return
		}

		hostname = normalizeHostname(hostname)
		for _, v := range hostnames {
			if v == hostname {
				return
			}
		}
		hostnames = append(hostnames, hostname)
	}

	openhabClient := openhab.New(config)

	things, err := openhabClient.GetThings(ctx)
	if err != nil {
		return nil, err
	}

	for _, thing := range *things {
		hostname := thing.Properties.ServiceName
		if hasShellyPrefix(hostname) {
			addHostname(hostname)
		} else {
			if logger.Trace {
				zap.L().Debug(fmt.Sprintf("Ignoring hostname %s", hostname))
			}
		}
	}

	return hostnames, nil
}
