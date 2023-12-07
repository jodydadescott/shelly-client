package mqtt

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdkclient "github.com/jodydadescott/shelly-client/sdk"
	sdktypes "github.com/jodydadescott/shelly-client/sdk/types"
)

type Config = types.Config

type ShellyClient = sdkclient.Client

type ShellyDeviceInfo = sdktypes.ShelllyDeviceInfo
type ShellyDeviceStatus = sdktypes.ShellyStatus
type ShellyConfig = sdktypes.ShellyConfig

type callback interface {
	GetConfig() (*Config, error)
	GetCTX() (context.Context, context.CancelFunc)
	WriteStdout(input any) error
}

// {"src":"shellypluswdus-441793ccc4d0","dst":"shelly/events","method":"NotifyStatus","params":{"ts":1701980272.34,"light:0":{"id":0,"brightness":3.00,"output":false,"source":"ui"}}}

func New(t callback) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "mqtt",
		Short: "Poles all ",
	}

	setOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Turn light on",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

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

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

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

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

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

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

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
