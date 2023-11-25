package shelly

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-client/sdk/shelly"
)

type callback interface {
	WriteStdout(any) error
	Shelly() (*shelly.Client, error)
}

const (
	ShellyConfigVar = "SHELLY_CONFIG"
)

func NewCmd(callback callback) *cobra.Command {

	var stageArg string
	var urlArg string
	var markupArg bool
	var configArg string

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

			config := configArg
			if config == "" {
				config = os.Getenv(ShellyConfigVar)
			}

			if config == "" {
				return fmt.Errorf("config is required")
			}

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			report, err := client.SetConfigFromFile(cmd.Context(), config)
			if err != nil {
				return err
			}

			return callback.WriteStdout(report)
		},
	}

	setConfigCmd.PersistentFlags().StringVarP(&configArg, "config", "c", "", fmt.Sprintf("Config file or Directory name. Directory will look for file name device-id then app-id; env var is %s", ShellyConfigVar))

	rootCmd.AddCommand(getConfigCmd, getStatusCmd, getInfoCmd, getMethodsCmd,
		getUpdatesCmd, rebootCmd, updateCmd,
		factoryResetCmd, resetWifiConfigCmd, setConfigCmd)
	return rootCmd
}
