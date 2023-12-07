package util

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"

	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/cmd/types"
	sdkclient "github.com/jodydadescott/shelly-client/sdk"
	sdktypes "github.com/jodydadescott/shelly-client/sdk/types"
)

type ShellyDeviceInfo = sdktypes.ShelllyDeviceInfo
type ShellyDeviceStatus = sdktypes.ShellyStatus
type ShellyConfig = sdktypes.ShellyConfig

type ShellyClient = sdkclient.Client

type Config = types.Config

func NewShellyClient(config *Config, hostname string) *ShellyClient {

	sdkConfig := &sdkclient.Config{
		Hostname: hostname,
	}

	if config.Password != nil {
		sdkConfig.Password = *config.Password
	}

	if config.Timeout != nil {
		sdkConfig.SendTimeout = *config.Timeout
	}

	return sdkclient.New(sdkConfig)
}

func Process(ctx context.Context, config *Config, action string, getDevStatus bool, do func(context.Context, string, *ShellyClient, *ShellyDeviceInfo, *ShellyDeviceStatus) error) error {

	execute := func(ctx context.Context, hostname string, workerID int) error {

		cleanupHostname := func(hostname string) string {
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

		hostname = cleanupHostname(hostname)

		client := NewShellyClient(config, hostname)
		defer client.Close()

		deviceInfo, err := client.Shelly().GetDeviceInfo(ctx)
		if err != nil {
			return fmt.Errorf("workerID %d, hostname  %s, [deviceInfo] failed with error %w", workerID, hostname, err)
		}

		var devStatus *ShellyDeviceStatus

		if getDevStatus {
			tmp, err := client.Shelly().GetStatus(ctx)
			if err != nil {
				return fmt.Errorf("workerID %d, hostname  %s, [deviceStatus] failed with error %w", workerID, hostname, err)
			}

			devStatus = tmp
		}

		err = do(ctx, hostname, client, deviceInfo, devStatus)

		if err != nil {
			return fmt.Errorf("workerID %d, hostname %s, deviceInfoID %s, deviceInfoApp %s, [%s] failed with error %w", workerID, hostname, *deviceInfo.ID, *deviceInfo.App, action, err)
		}

		return nil
	}

	jobCount := len(config.Hostnames)

	if jobCount == 0 {
		return fmt.Errorf("one or more hostnames required")
	}

	if jobCount == 1 {
		zap.L().Debug("there is only 1 job to execute")
		return execute(ctx, config.Hostnames[0], 0)
	}

	zap.L().Debug(fmt.Sprintf("there are %d jobs to execute", jobCount))

	jobchan := make(chan string, jobCount)
	errchan := make(chan error, jobCount)

	wg := &sync.WaitGroup{}

	worker := func(id int) {

		wg.Add(1)

		go func() {

			defer wg.Done()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			for {
				select {

				case <-ctx.Done():
					zap.L().Debug("worker cancelled")
					return

				case hostname := <-jobchan:
					zap.L().Debug(fmt.Sprintf("worker %d processing hostname %s", id, hostname))
					err := execute(ctx, hostname, id)
					if err != nil {
						errchan <- err
					}

				default:
					zap.L().Debug("no more work")
					return

				}
			}

		}()
	}

	worker(1)

	if jobCount > 5 {
		worker(2)
		worker(3)
		worker(4)
		worker(5)
	}

	for _, hostname := range config.Hostnames {
		jobchan <- hostname
	}

	var errors *multierror.Error

	zap.L().Debug("waiting on jobs to complete")
	wg.Wait()
	zap.L().Debug("jobs have completed")

	for {
		select {

		case <-ctx.Done():
			zap.L().Debug("Context cancelled")
			errors = multierror.Append(errors, fmt.Errorf("cancelled"))
			return errors.ErrorOrNil()

		case err := <-errchan:
			if err != nil {
				zap.L().Debug("Adding error")
				errors = multierror.Append(errors, err)
			}

		default:
			zap.L().Debug("no more errors")
			return errors.ErrorOrNil()

		}
	}

}
