package light

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

var (
	truePointer  = true
	falsePointer = false
)

type Config = types.Config

type ShellyClient = sdk_client.Client

type ShellyDeviceInfo = shelly_types.DeviceInfo
type ShellyDeviceStatus = shelly_types.Status
type ShellyConfig = sdk_types.Config

type callback interface {
	GetConfig(context.Context) (*Config, error)
	GetCTX() (context.Context, context.CancelFunc)
	WriteStdout(input any) error
}

func New(t callback) *cobra.Command {

	var brightnessArg string

	getBrightness := func() (*float64, error) {

		if brightnessArg == "" {
			return nil, fmt.Errorf("brightness arg not set")
		}

		brightness, err := strconv.ParseFloat(brightnessArg, 32)
		if err != nil {
			return nil, err
		}

		return &brightness, nil
	}

	getIds := func(ctx context.Context, shellyClient *ShellyClient, args []string) ([]int, error) {

		if len(args) == 0 {
			return nil, fmt.Errorf("one or more IDs is required. They can be space of comma delineated. You can also use 'all'")
		}

		var results []int

		if len(args) == 1 {
			if strings.ToLower(args[0]) == "all" {

				shellyConfig, err := shellyClient.GetConfig(ctx, false)
				if err != nil {
					return nil, err
				}
				for _, lightConfig := range shellyConfig.Light {
					results = append(results, *lightConfig.ID)
				}
				return results, nil
			}
		}

		var errors *multierror.Error

		for _, arg := range args {
			for _, sub := range strings.Split((strings.TrimSpace(arg)), ",") {
				id, err := strconv.Atoi(sub)
				if err != nil {
					errors = multierror.Append(errors, err)
				} else {
					results = append(results, id)
				}
			}
		}

		return results, errors.ErrorOrNil()
	}

	rootCmd := &cobra.Command{
		Use:   "light",
		Short: "Turn light on, off, or set brigtness level",
	}

	setOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Turn light on",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "set on"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				ids, err := getIds(ctx, client, args)
				if err != nil {
					return err
				}

				var errors *multierror.Error

				for _, id := range ids {
					err := client.Light().Set(ctx, id, &truePointer, nil)
					if err != nil {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, id, action, err.Error()))
						errors = multierror.Append(errors, err)
					} else {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] completed", hostname, *deviceInfo.ID, *deviceInfo.App, id, action))
					}
				}

				return errors.ErrorOrNil()

			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	setOffCmd := &cobra.Command{
		Use:   "off",
		Short: "Turn light off",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "set off"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				ids, err := getIds(ctx, client, args)
				if err != nil {
					return err
				}

				var errors *multierror.Error

				for _, id := range ids {
					err := client.Light().Set(ctx, id, &falsePointer, nil)
					if err != nil {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, id, action, err.Error()))
						errors = multierror.Append(errors, err)
					} else {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] completed", hostname, *deviceInfo.ID, *deviceInfo.App, id, action))
					}
				}

				return errors.ErrorOrNil()
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	setBrightnessCmd := &cobra.Command{
		Use:   "bright",
		Short: "Sets light brightness",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			brightness, err := getBrightness()
			if err != nil {
				return err
			}

			action := "set brightness"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				ids, err := getIds(ctx, client, args)
				if err != nil {
					return err
				}

				var errors *multierror.Error

				for _, id := range ids {
					err := client.Light().Set(ctx, id, nil, brightness)
					if err != nil {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, id, action, err.Error()))
						errors = multierror.Append(errors, err)
					} else {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] completed", hostname, *deviceInfo.ID, *deviceInfo.App, id, action))
					}
				}

				return errors.ErrorOrNil()
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	toggleCmd := &cobra.Command{
		Use:   "toggle",
		Short: "Toggles switch",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "toggle"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				ids, err := getIds(ctx, client, args)
				if err != nil {
					return err
				}

				var errors *multierror.Error

				for _, id := range ids {
					err := client.Light().Toggle(ctx, id)
					if err != nil {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, id, action, err.Error()))
						errors = multierror.Append(errors, err)
					} else {
						t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s, lightID %d: [%s] completed", hostname, *deviceInfo.ID, *deviceInfo.App, id, action))
					}
				}

				return errors.ErrorOrNil()
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	rootCmd.AddCommand(toggleCmd, setOnCmd, setOffCmd, setBrightnessCmd)
	return rootCmd
}
