package shelly

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-client/cmd/files"
	"github.com/jodydadescott/shelly-client/sdk/shelly"
)

type callback interface {
	WriteStdout(any) error
	WriteStderr(s string)
	Shelly() (*shelly.Client, error)
	RebootDevice(ctx context.Context) error
	GetFiles() (*files.Files, error)
}

func NewCmd(callback callback) *cobra.Command {

	var stageArg string
	var urlArg string
	var disableAutoRebootArg bool
	var markupArg bool

	rootCmd := &cobra.Command{
		Use:   "shelly",
		Short: "Shelly Component",
	}

	getConfigCmd := &cobra.Command{
		Use:   "get-config",
		Short: "Returns config",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			config, err := client.GetConfig(cmd.Context(), markupArg)
			if err != nil {
				return err
			}

			return callback.WriteStdout(config)
		},
	}

	getConfigCmd.PersistentFlags().BoolVar(&markupArg, "markup", false, "returns config that can be used as a template")

	getStatusCmd := &cobra.Command{
		Use:   "get-status",
		Short: "Returns status",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.GetStatus(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getInfoCmd := &cobra.Command{
		Use:   "get-info",
		Short: "Returns device info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.GetDeviceInfo(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getMethodsCmd := &cobra.Command{
		Use:   "get-methods",
		Short: "Returns all available RPC methods for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.ListMethods(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getUpdatesCmd := &cobra.Command{
		Use:   "get-updates",
		Short: "Returns available update info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.CheckForUpdate(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	rebootCmd := &cobra.Command{
		Use:   "reboot",
		Short: "Executes device reboot",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			err = client.Reboot(cmd.Context())
			if err != nil {
				return err
			}

			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Returns available update info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			params := &shelly.ShellyUpdateConfig{}

			if stageArg != "" {
				params.Stage = &stageArg
			}

			if urlArg != "" {
				params.Url = &urlArg
			}

			return client.Update(cmd.Context(), params)
		},
	}

	updateCmd.PersistentFlags().StringVar(&stageArg, "stage", "", "The type of the new version - either stable or beta. By default updates to stable version. Optional")
	updateCmd.PersistentFlags().StringVar(&urlArg, "url", "", "Url address of the update. Optional")

	factoryResetCmd := &cobra.Command{
		Use:   "factory-reset",
		Short: "Executes factory reset",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			return client.FactoryReset(cmd.Context())
		},
	}

	resetWifiConfigCmd := &cobra.Command{
		Use:   "reset-wifi-config",
		Short: "Executes Wifi config reset",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			return client.ResetWiFiConfig(cmd.Context())
		},
	}

	setConfigCmd := &cobra.Command{
		Use:   "set-config",
		Short: "Sets config",
		RunE: func(cmd *cobra.Command, args []string) error {

			var config *shelly.ShellyConfig

			setConfig := func(file *files.File) error {

				var errors *multierror.Error

				callback.WriteStderr(fmt.Sprintf("Using file %s", file.FullName))

				err := json.Unmarshal(file.Bytes, &config)

				if err != nil {
					errors = multierror.Append(errors, err)
					err = yaml.Unmarshal(file.Bytes, &config)

					if err != nil {
						errors = multierror.Append(errors, err)
						errors = multierror.Append(errors, fmt.Errorf("invalid format. Expect JSON or YAML"))
						return errors.ErrorOrNil()
					}
				}

				return nil
			}

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			files, err := callback.GetFiles()
			if err != nil {
				return err
			}

			file := files.GetNamedFile()
			if file != nil {
				return setConfig(file)
			}

			file = files.GetSTDIN()
			if file != nil {
				return setConfig(file)
			}

			device, err := client.GetDeviceInfo(cmd.Context())
			if err != nil {
				return err
			}

			file = files.GetFile(*device.ID)
			if file != nil {
				return setConfig(file)
			}

			file = files.GetFile(*device.App)
			if file != nil {
				return setConfig(file)
			}

			report := client.SetConfig(cmd.Context(), config)

			err = report.Error()

			if err != nil {
				return err
			}

			if report.RebootRequired() {
				if disableAutoRebootArg {
					callback.WriteStderr("reboot is required; autoreboot is disabled")
					return nil
				}

				callback.WriteStderr("rebooting")
				return callback.RebootDevice(cmd.Context())
			}

			return nil
		},
	}

	setConfigCmd.PersistentFlags().BoolVar(&disableAutoRebootArg, "disable-autoreboot", false, "disable automatic reboot (if reboot is necessary)")

	rootCmd.AddCommand(getConfigCmd, getStatusCmd, getInfoCmd, getMethodsCmd,
		getUpdatesCmd, rebootCmd, updateCmd,
		factoryResetCmd, resetWifiConfigCmd, setConfigCmd)
	return rootCmd
}
