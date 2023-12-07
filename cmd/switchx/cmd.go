package switchx

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdkclient "github.com/jodydadescott/shelly-client/sdk"
	sdktypes "github.com/jodydadescott/shelly-client/sdk/types"
)

var (
	truePointer  = true
	falsePointer = false
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

func New(t callback) *cobra.Command {

	getIds := func(ctx context.Context, shellyClient *ShellyClient, args []string) ([]int, error) {

		if len(args) == 0 {
			return nil, fmt.Errorf("one or more IDs is required. They can be space of comma delineated. You can also use 'all'")
		}

		var results []int

		if len(args) == 1 {
			if strings.ToLower(args[0]) == "all" {

				shellyConfig, err := shellyClient.Shelly().GetConfig(ctx)
				if err != nil {
					return nil, err
				}
				for _, lightConfig := range shellyConfig.Light {
					results = append(results, lightConfig.ID)
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
		Use:   "switch",
		Short: "Turn switch on or off",
	}

	setOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Turn switch on",
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
					err := client.Switch().Set(ctx, id, &truePointer)
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
		Short: "Turn switch off",
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
					err := client.Switch().Set(ctx, id, &falsePointer)
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
					err := client.Switch().Toggle(ctx, id)
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

	rootCmd.AddCommand(setOnCmd, setOffCmd, toggleCmd)
	return rootCmd
}
