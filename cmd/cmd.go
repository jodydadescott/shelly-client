package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/go-multierror"
	"github.com/hokaccha/go-prettyjson"

	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/jodydadescott/openhab-go-sdk"
	"github.com/jodydadescott/unifi-go-sdk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-client/cmd/light"
	"github.com/jodydadescott/shelly-client/cmd/mqtt"
	"github.com/jodydadescott/shelly-client/cmd/switchx"
	"github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

type ShellyDeviceInfo = shelly_types.DeviceInfo
type ShellyStatus = shelly_types.Status
type ShellyConfig = sdk_types.Config
type ShellyClient = sdk_client.Client
type ShellyUpdateConfig = shelly_types.UpdateConfig
type ShellyRPCMethods = shelly_types.RPCMethods
type ShellyUpdateReport = shelly_types.UpdatesReport

type Config = types.Config
type UnifiConfig = unifi.Config
type Logger = logger.Config

type Cmd struct {
	ctx context.Context
	*cobra.Command
	configFileArg     string
	hostnameArg       []string
	passwordArg       string
	outputArg         string
	urlArg            string
	timeoutArg        string
	debugLevelArg     string
	unifiArg          string
	openhabArg        string
	rebootForceArg    bool
	setConfigForceArg bool
}

func NewCmd() *Cmd {

	t := &Cmd{}

	rootCmd := &cobra.Command{

		Use: BinaryName,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			debugLevel := t.debugLevelArg
			if debugLevel == "" {
				debugLevel = os.Getenv(DebugEnvVar)
			}
			if debugLevel != "" {
				loggerConfig := &logger.Config{}
				err := loggerConfig.ParseLogLevel(debugLevel)
				if err != nil {
					return err
				}
				logger.SetConfig(loggerConfig)
			}

			return nil
		},

		SilenceUsage: true,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manages device(s) configuration",
	}

	getConfigCmd := &cobra.Command{
		Use:   "get",
		Short: "Returns device config for specified device(s)",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
			defer client.Close()

			result, err := client.GetConfig(ctx, false)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	compareConfigCmd := &cobra.Command{
		Use:   "compare",
		Short: "Returns OS code 0 if the running config matches the config",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
			defer client.Close()

			runningConfig, err := client.GetConfig(ctx, false)
			if err != nil {
				return err
			}

			renderedConfig, err := sdk_client.GetConfig(ctx, client)
			if err != nil {
				return err
			}

			runningConfig.Sanatize()
			renderedConfig.Sanatize()

			if runningConfig.Equals(renderedConfig) {
				return t.WriteStdout("EQUAL")
			}

			return t.WriteStdout("NOT EQUAL")
		},
	}

	setConfigCmd := &cobra.Command{
		Use:   "set",
		Short: "Sets specified config for device(s)",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			action := "set config"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				shellyConfig, err := sdk_client.GetConfig(ctx, client)
				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				if t.setConfigForceArg {
					zap.L().Debug(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] set config force is true", hostname, *deviceInfo.ID, *deviceInfo.App, action))
				}

				configReport, err := client.SetConfig(ctx, shellyConfig, t.setConfigForceArg)
				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				if configReport.NoChange {
					return t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] no change in config", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				}

				if configReport.RebootRequired {
					return t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed and rebooted", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				}

				return t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	setConfigCmd.PersistentFlags().BoolVarP(&t.setConfigForceArg, "force", "f", false, "config will not be set if there are no changes; this flag forces set config")

	renderConfigCmd := &cobra.Command{
		Use:   "render",
		Short: "Shows the rendered config for the target device(s). This can be used to create desired config",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
			defer client.Close()

			shellyConfig, err := sdk_client.GetConfig(ctx, client)
			if err != nil {
				return err
			}

			return t.WriteStdout(shellyConfig)
		},
	}

	configCmd.AddCommand(getConfigCmd, setConfigCmd, renderConfigCmd, compareConfigCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Returns device(s) information",
	}

	infoStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Returns device status",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			if len(config.Hostnames) == 1 {
				client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
				defer client.Close()

				result, err := client.GetStatus(ctx)
				if err != nil {
					return err
				}

				return t.WriteStdout(result)
			}

			var results []*ShellyStatus

			for _, hostname := range config.Hostnames {
				client := util.NewShellyClient(config, util.CleanupHostname(hostname))
				defer client.Close()
				result, err := client.GetStatus(ctx)
				if err != nil {
					return err
				}
				results = append(results, result)
			}

			return t.WriteStdout(results)
		},
	}

	devInfoCmd := &cobra.Command{
		Use:   "dev",
		Short: "Returns device info",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			if len(config.Hostnames) == 1 {
				client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
				defer client.Close()

				result, err := client.GetDeviceInfo(ctx)
				if err != nil {
					return err
				}

				return t.WriteStdout(result)
			}

			var results []*ShellyDeviceInfo

			for _, hostname := range config.Hostnames {
				client := util.NewShellyClient(config, util.CleanupHostname(hostname))
				defer client.Close()
				result, err := client.GetDeviceInfo(ctx)
				if err != nil {
					return err
				}
				results = append(results, result)
			}

			return t.WriteStdout(results)
		},
	}

	methodsInfoCmd := &cobra.Command{
		Use:   "methods",
		Short: "Returns RPC methods for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			if len(config.Hostnames) == 1 {
				client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
				defer client.Close()

				result, err := client.ListMethods(ctx)
				if err != nil {
					return err
				}

				return t.WriteStdout(result)
			}

			var results []*ShellyRPCMethods

			for _, hostname := range config.Hostnames {
				client := util.NewShellyClient(config, util.CleanupHostname(hostname))
				defer client.Close()
				result, err := client.ListMethods(ctx)
				if err != nil {
					return err
				}
				results = append(results, result)
			}

			return t.WriteStdout(results)
		},
	}

	infoCmd.AddCommand(infoStatusCmd, devInfoCmd, methodsInfoCmd)

	firmwareCmd := &cobra.Command{
		Use:   "firmware",
		Short: "Manages device(s) firmware (updates)",
	}

	availableFirmwareCmd := &cobra.Command{
		Use:   "available",
		Short: "Returns available updates for specified device",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			if len(config.Hostnames) == 1 {
				client := util.NewShellyClient(config, util.CleanupHostname(config.Hostnames[0]))
				defer client.Close()

				result, err := client.CheckForUpdate(ctx)
				if err != nil {
					return err
				}

				return t.WriteStdout(result)
			}

			var results []*ShellyUpdateReport

			for _, hostname := range config.Hostnames {
				client := util.NewShellyClient(config, util.CleanupHostname(hostname))
				defer client.Close()
				result, err := client.CheckForUpdate(ctx)
				if err != nil {
					return err
				}
				results = append(results, result)
			}

			return t.WriteStdout(results)
		},
	}

	updateStableFirmwareCmd := &cobra.Command{
		Use:   "update-stable",
		Short: "Updates devices(s) firmware to latest stable",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			stage := "beta"
			action := fmt.Sprintf("update to %s", stage)

			params := &ShellyUpdateConfig{
				Stage: &stage,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				err := client.Update(ctx, params)

				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	updateBetaFirmwareCmd := &cobra.Command{
		Use:   "update-beta",
		Short: "Updates devices(s) firmware to latest stable",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			stage := "beta"
			action := fmt.Sprintf("update to %s", stage)

			params := &ShellyUpdateConfig{
				Stage: &stage,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				err := client.Update(ctx, params)

				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	updateURLFirmwareCmd := &cobra.Command{
		Use:   "update-url",
		Short: "Updates device(s) firmware to specified URL",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			if config.UpdateURL == nil {
				return fmt.Errorf("update URL is required")
			}

			url := *config.UpdateURL

			if !strings.HasPrefix(url, "http://") {
				if !strings.HasPrefix(url, "https://") {
					url = "https://" + url
				}
			}

			action := fmt.Sprintf("update with url %s", url)

			params := &ShellyUpdateConfig{
				Url: &url,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				err := client.Update(ctx, params)

				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	firmwareCmd.AddCommand(availableFirmwareCmd, updateStableFirmwareCmd, updateBetaFirmwareCmd, updateURLFirmwareCmd)

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reboots, Factory Resets or Wifi Resets specified device(s)",
	}

	rebootCmd := &cobra.Command{
		Use:   "reboot",
		Short: "Reboots device(s)",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "reboot"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				isRestartRequired := func() bool {

					if deviceStatus.System == nil {
						return false
					}

					if deviceStatus.System.RestartRequired == nil {
						return false
					}

					return *deviceStatus.System.RestartRequired
				}

				restart := func() bool {

					restartRequired := isRestartRequired()

					if restartRequired {
						zap.L().Debug(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] is required", hostname, *deviceInfo.ID, *deviceInfo.App, action))
						return true
					}

					if t.rebootForceArg {
						zap.L().Debug(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] is not required but force arg is true", hostname, *deviceInfo.ID, *deviceInfo.App, action))
						return true
					}

					zap.L().Debug(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] is not required and force arg is false", hostname, *deviceInfo.ID, *deviceInfo.App, action))
					return false
				}

				if !restart() {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] not required", hostname, *deviceInfo.ID, *deviceInfo.App, action))
					return nil
				}

				err := client.Reboot(ctx)
				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				return nil
			}

			return util.Process(ctx, config, action, true, do)
		},
	}

	rebootCmd.PersistentFlags().BoolVarP(&t.rebootForceArg, "force", "f", false, "reboots device even if device reboot is not required")

	// Need a sanity flag like force TODO
	factoryResetCmd := &cobra.Command{
		Use:   "factory",
		Short: "Resets device(s) to factory settings",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "factory reset"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				err := client.FactoryReset(ctx)

				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil

			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	// Need a sanity flag like force TODO
	resetWifiConfigCmd := &cobra.Command{
		Use:   "wifi",
		Short: "Resets device(s) Wifi config",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			action := "reset wifi"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyStatus) error {

				err := client.ResetWiFiConfig(ctx)

				if err != nil {
					t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStderr(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil

			}

			return util.Process(ctx, config, action, false, do)
		},
	}

	resetCmd.AddCommand(rebootCmd, factoryResetCmd, resetWifiConfigCmd)

	listHostnamesCmd := &cobra.Command{
		Use:   "list-hostnames",
		Short: "Returns list of Shelly hostnames",
	}

	unifiListHostnamesCmd := &cobra.Command{
		Use:   "unifi",
		Short: "Returns list of Shelly hostnames from Ubiquiti controller",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.Unifi == nil {
				return fmt.Errorf("config.Unifi is required")
			}

			hostnames, err := util.GetUnifiHostnames(config.Unifi)
			if err != nil {
				return err
			}

			for _, hostname := range hostnames {
				err := t.WriteStdout(hostname)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	openHabListHostnamesCmd := &cobra.Command{
		Use:   "openhab",
		Short: "Returns list of Shelly hostnames from Ubiquiti controller",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.OpenHAB == nil {
				return fmt.Errorf("config.OpenHAB is required")
			}

			hostnames, err := util.GetOpenHABHostnames(ctx, config.OpenHAB)
			if err != nil {
				return err
			}

			for _, hostname := range hostnames {
				err := t.WriteStdout(hostname)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	listHostnamesCmd.AddCommand(unifiListHostnamesCmd, openHabListHostnamesCmd)

	diffHostnamesCmd := &cobra.Command{
		Use:   "diff-hostnames",
		Short: "Returns list of Shelly hostnames that are present in one system but not the other",
	}

	openHabUnifiDiffHostnamesCmd := &cobra.Command{
		Use:   "openhab-unifi",
		Short: "Shows hostnames present in OpenHAB and not in Unifi",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.OpenHAB == nil {
				return fmt.Errorf("config.OpenHAB is required")
			}

			if config.OpenHAB == nil {
				return fmt.Errorf("config.OpenHAB is required")
			}

			openhabHostnames, err := util.GetOpenHABHostnames(ctx, config.OpenHAB)
			if err != nil {
				return err
			}

			unifiHostnames, err := util.GetUnifiHostnames(config.Unifi)
			if err != nil {
				return err
			}

			has := func(hostname string) bool {
				for _, h := range unifiHostnames {
					if hostname == h {
						return true
					}
				}
				return false
			}

			for _, h := range openhabHostnames {
				if !has(h) {
					t.WriteStdout(h)
				}
			}

			return nil
		},
	}

	unifiOpenHabDiffHostnamesCmd := &cobra.Command{
		Use:   "unifi-openhab",
		Short: "Shows hostnames present in Unifi and not in OpenHAB",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.OpenHAB == nil {
				return fmt.Errorf("config.OpenHAB is required")
			}

			if config.OpenHAB == nil {
				return fmt.Errorf("config.OpenHAB is required")
			}

			openhabHostnames, err := util.GetOpenHABHostnames(ctx, config.OpenHAB)
			if err != nil {
				return err
			}

			unifiHostnames, err := util.GetUnifiHostnames(config.Unifi)
			if err != nil {
				return err
			}

			has := func(hostname string) bool {
				for _, h := range openhabHostnames {
					if hostname == h {
						return true
					}
				}
				return false
			}

			for _, h := range unifiHostnames {
				if !has(h) {
					t.WriteStdout(h)
				}
			}

			return nil
		},
	}

	diffHostnamesCmd.AddCommand(openHabUnifiDiffHostnamesCmd, unifiOpenHabDiffHostnamesCmd)

	// Override the default help so we can use the -h for host
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	rootCmd.PersistentFlags().StringVarP(&t.debugLevelArg, "debug", "D", "", fmt.Sprintf("debug level (WIRE, DEBUG, INFO, WARN, ERROR) to STDERR; env var is %s", DebugEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.configFileArg, "config", "c", "", fmt.Sprintf("Config; optionally use env var '%s'", ShellyConfigEnvVar))
	rootCmd.PersistentFlags().StringSliceVarP(&t.hostnameArg, "hostname", "h", []string{}, fmt.Sprintf("Hostname; optionally use env var '%s'", ShellyHostnameEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.passwordArg, "password", "p", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyPasswordEnvVar))
	rootCmd.PersistentFlags().StringVar(&t.passwordArg, "update-url", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyURLEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.outputArg, "output", "o", ShellyOutputDefault, fmt.Sprintf("Output format. One of: prettyjson | json | jsonpath | yaml ; Optionally use env var '%s'", ShellyOutputEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.timeoutArg, "timeout", "t", "", "The timeout in seconds for the websocket call to the device")
	rootCmd.PersistentFlags().StringVar(&t.unifiArg, "unifi", "", "Use Unifi controller to get hostname(s). Must be enable, disable, true, or false")
	rootCmd.PersistentFlags().StringVar(&t.openhabArg, "openhab", "", "Use Openhab controller to get hostname(s). Must be enable, disable, true, or false")
	rootCmd.AddCommand(configCmd, infoCmd, resetCmd, firmwareCmd, listHostnamesCmd, diffHostnamesCmd, light.New(t), switchx.New(t), mqtt.New(t))
	t.Command = rootCmd

	return t
}

// WriteStdout writes any in desired format to STDOUT
func (t *Cmd) WriteStdout(input any) error {

	output, err := t.anyToOut(input)
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

// WriteStderr writes any in desired format to STDERR
func (t *Cmd) WriteStderr(input any) error {

	output, err := t.anyToOut(input)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, output)
	return nil
}

func (t *Cmd) anyToOut(input any) (string, error) {

	switch v := input.(type) {

	case nil:
		return "", nil

	case string:
		return fmt.Sprint(v), nil

	case *string:
		return fmt.Sprint(*v), nil

	}

	// Arg can be in the format of 'format' or 'format=...'.
	// For example jsonpath expects a second arg such as jsonpath=$.

	outputArgSplit := strings.Split(t.outputArg, "=")

	switch strings.ToLower(outputArgSplit[0]) {

	case "prettyjson":
		data, err := prettyjson.Marshal(input)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(strings.TrimSpace(string(data))), nil

	case "jsonpath":

		if len(outputArgSplit) > 1 {

			v := interface{}(nil)
			data, err := json.Marshal(input)
			if err != nil {
				return "", err
			}

			err = json.Unmarshal(data, &v)
			if err != nil {
				return "", err
			}

			data2, err := jsonpath.Get(outputArgSplit[1], v)
			if err != nil {
				return "", err
			}

			switch data3 := data2.(type) {
			case string:
				return fmt.Sprint(data3), nil

			case []interface{}:
				s := reflect.ValueOf(data3)
				result := ""
				for i := 0; i < s.Len(); i++ {
					if result == "" {
						result = fmt.Sprint(s.Index(i))
					} else {
						result = "/n" + fmt.Sprint(s.Index(i))
					}
				}
				return result, nil
			}

			panic("you should not be here")
		}

		return "", fmt.Errorf("missing jsonpath arg. Expect jsonpath")

	case "json":
		data, err := json.Marshal(input)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(strings.TrimSpace(string(data))), err

	case "yaml":
		data, err := yaml.Marshal(input)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(strings.TrimSpace(string(data))), err

	}

	return "", fmt.Errorf("format type %s is unknown", t.outputArg)
}

func (t *Cmd) GetCTX() (context.Context, context.CancelFunc) {

	if t.ctx == nil {

		zap.L().Debug("creating initial ctx")

		ctx, cancel := context.WithCancel(t.Context())
		t.ctx = ctx

		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt)

		go func() {

			select {

			case <-interruptChan: // first signal, cancel context
				zap.L().Debug("caught interrupt; calling cancel()")
				cancel()

			case <-ctx.Done():
				cancel()

			}

		}()

	}

	zap.L().Debug("returning child ctx")
	return context.WithCancel(t.ctx)
}

func (t *Cmd) GetConfig(ctx context.Context) (*Config, error) {

	var config *Config

	initLog := func() error {

		getDebugLevel := func() string {

			x := t.debugLevelArg
			if x != "" {
				zap.L().Debug(fmt.Sprintf("DebugLevel is %s from args", x))
				return x
			}

			x = os.Getenv(DebugEnvVar)
			if x != "" {
				zap.L().Debug(fmt.Sprintf("DebugLevel is %s from envvar %s", x, DebugEnvVar))
				return x
			}

			return ""
		}

		debugLevel := getDebugLevel()

		if debugLevel != "" {
			if config == nil || config.Logger == nil {
				loggerConfig := &logger.Config{}
				err := loggerConfig.ParseLogLevel(debugLevel)
				if err != nil {
					return err
				}
				logger.SetConfig(loggerConfig)
				return nil
			}
			err := config.Logger.ParseLogLevel(debugLevel)
			if err != nil {
				return err
			}
		}

		if config != nil && config.Logger != nil {
			logger.SetConfig(config.Logger)
		}

		return nil
	}

	getConfigFile := func() string {

		x := t.configFileArg
		if x != "" {
			zap.L().Debug(fmt.Sprintf("Config file is %s from args", x))
			return x
		}

		x = os.Getenv(ShellyConfigEnvVar)
		if x != "" {
			zap.L().Debug(fmt.Sprintf("Config file is %s from envvar %s", x, ShellyConfigEnvVar))
			return x
		}

		zap.L().Debug("Config file not set")

		return ""
	}

	initFromBytes := func(input []byte) error {

		var errs *multierror.Error
		err := json.Unmarshal(input, &config)
		if err == nil {
			return nil
		}

		errs = multierror.Append(errs, err)

		err = yaml.Unmarshal(input, &config)
		if err == nil {
			return nil
		}

		errs = multierror.Append(errs, err)

		return errs.ErrorOrNil()
	}

	initFromFile := func(filename string) error {

		if filename == "" {
			config = &Config{}
			return nil
		}

		_, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("problem with config file %s; %w", filename, err)
		}

		fileStats, err := os.Stat(filename)
		if err != nil {
			return err
		}

		permissions := fileStats.Mode().Perm()
		if permissions != SecureFilePerm {
			t.WriteStderr(fmt.Sprintf("WARNING: config file %s has overly promiscuous permissions", filename))
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		return initFromBytes(content)
	}

	compareSlices := func(a, b []string) bool {

		has := func(x string, y []string) bool {
			for _, v := range y {
				if x == v {
					return true
				}
			}
			return false
		}

		for _, v := range a {
			if !has(v, b) {
				return false
			}
		}

		for _, v := range b {
			if !has(v, a) {
				return false
			}
		}

		return true
	}

	loadBase := func() error {

		if t.outputArg != "" {
			config.Output = &t.outputArg
			zap.L().Debug(fmt.Sprintf("Output %s loaded from args", t.outputArg))
		}

		if t.timeoutArg != "" {
			timeout, err := time.ParseDuration(t.timeoutArg)
			if err != nil {
				return err
			}
			config.Timeout = &timeout
			zap.L().Debug(fmt.Sprintf("Timeout %s loaded from args", t.timeoutArg))
		}

		return nil
	}

	loadHostnames := func() {

		if len(t.hostnameArg) > 0 {
			if len(config.Hostnames) > 0 {
				if compareSlices(t.hostnameArg, config.Hostnames) {
					zap.L().Debug(fmt.Sprintf("hostnames %s in config are the same as from arg", config.Hostnames))
					return
				}
				zap.L().Debug(fmt.Sprintf("hostnames %s in config overwritten by %s from args", config.Hostnames, t.hostnameArg))
				config.Hostnames = t.hostnameArg
				return
			}

			zap.L().Debug(fmt.Sprintf("hostnames is %s from args", t.hostnameArg))
			config.Hostnames = t.hostnameArg
			return
		}

		x := os.Getenv(ShellyHostnameEnvVar)
		if x != "" {
			x = strings.Trim(x, " ")
			hostnames := strings.Split(x, ",")

			if len(hostnames) > 0 {
				if len(config.Hostnames) > 0 {
					if compareSlices(hostnames, config.Hostnames) {
						zap.L().Debug(fmt.Sprintf("hostnames %s in config are the same as from envvar %s", config.Hostnames, ShellyHostnameEnvVar))
						return
					}
					zap.L().Debug(fmt.Sprintf("hostnames %s in config overwritten by %s from envvar %s", config.Hostnames, hostnames, ShellyHostnameEnvVar))
					config.Hostnames = hostnames
					return
				}

				zap.L().Debug(fmt.Sprintf("hostnames is %s from envvar %s", hostnames, ShellyHostnameEnvVar))
				config.Hostnames = hostnames
				return
			}

		}

		if len(config.Hostnames) > 0 {
			zap.L().Debug(fmt.Sprintf("hostnames is %s from config", config.Hostnames))
			return
		}

		zap.L().Debug(fmt.Sprintf("Hostnames not set in arg, envvar %s or config", ShellyHostnameEnvVar))
	}

	loadShelly := func() {

		if config.Shelly == nil {
			config.Shelly = &ShellyConfig{}
		}

		x := t.passwordArg
		if x != "" {
			if config.Shelly.Password != "" {
				if config.Shelly.Password == x {
					zap.L().Debug(fmt.Sprintf("Password %s in config is the same as from arg", x))
					return
				}
				zap.L().Debug(fmt.Sprintf("Password %s overwritten with %s from args", config.Shelly.Password, x))
				config.Shelly.Password = x
				return
			}

			zap.L().Debug(fmt.Sprintf("Password is %s from args", x))
			config.Shelly.Password = x
			return
		}

		x = os.Getenv(ShellyPasswordEnvVar)
		if x != "" {
			if config.Shelly.Password != "" {
				if config.Shelly.Password == x {
					zap.L().Debug(fmt.Sprintf("Password in config is the same as from envvar %s", ShellyPasswordEnvVar))
					return
				}
				zap.L().Debug(fmt.Sprintf("Password overwritten with from envvar %s", ShellyPasswordEnvVar))
				config.Shelly.Password = x
				return
			}

			zap.L().Debug(fmt.Sprintf("Password in envvar %s", ShellyPasswordEnvVar))
			config.Shelly.Password = x
			return
		}

		if config.Shelly.Password != "" {
			zap.L().Debug("password configured is config")
			return
		}

		zap.L().Debug(fmt.Sprintf("Password not set in arg, envvar %s or config", ShellyPasswordEnvVar))
	}

	loadOutput := func() {

		x := t.outputArg
		if x != "" {
			if config.Output != nil {
				if *config.Output == x {
					zap.L().Debug(fmt.Sprintf("Output %s in config is the same as from arg", x))
					return
				}
				zap.L().Debug(fmt.Sprintf("Output %s overwritten with %s from args", *config.Output, x))
				config.Output = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("Output is %s from args", x))
			config.Output = &x
			return
		}

		x = os.Getenv(ShellyOutputEnvVar)
		if x != "" {
			if config.Output != nil {
				if *config.Output == x {
					zap.L().Debug(fmt.Sprintf("Output %s in config is the same as from envvar %s", x, ShellyOutputEnvVar))
					return
				}
				zap.L().Debug(fmt.Sprintf("Output %s overwritten with %s from envvar %s", *config.Output, x, ShellyOutputEnvVar))
				config.Output = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("Output is %s from envvar %s", x, ShellyOutputEnvVar))
			config.Output = &x
			return
		}

		if config.Output != nil {
			zap.L().Debug(fmt.Sprintf("Output is %s from config", *config.Output))
			return
		}

		zap.L().Debug(fmt.Sprintf("Output not set in arg, envvar %s or config", ShellyOutputEnvVar))

	}

	loadUpdateURL := func() {

		x := t.urlArg
		if x != "" {
			if config.UpdateURL != nil {
				if *config.UpdateURL == x {
					zap.L().Debug(fmt.Sprintf("UpdateURL %s in config is the same as from arg", x))
					return
				}
				zap.L().Debug(fmt.Sprintf("UpdateURL %s overwritten with %s from args", *config.UpdateURL, x))
				config.UpdateURL = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("UpdateURL is %s from args", x))
			config.UpdateURL = &x
			return
		}

		x = os.Getenv(ShellyURLEnvVar)
		if x != "" {
			if config.UpdateURL != nil {
				if *config.UpdateURL == x {
					zap.L().Debug(fmt.Sprintf("Url %s in config is the same as from envvar %s", x, ShellyURLEnvVar))
					return
				}
				zap.L().Debug(fmt.Sprintf("Url %s overwritten with %s from envvar %s", *config.UpdateURL, x, ShellyURLEnvVar))
				config.UpdateURL = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("Url is %s from envvar %s", x, ShellyURLEnvVar))
			config.UpdateURL = &x
			return
		}

		if config.UpdateURL != nil {
			zap.L().Debug(fmt.Sprintf("UpdateURL is %s from config", *config.UpdateURL))
			return
		}

		zap.L().Debug(fmt.Sprintf("UpdateURL not set in arg, envvar %s or config", ShellyURLEnvVar))

	}

	loadTimeout := func() error {

		x := t.timeoutArg
		if x != "" {

			d, err := time.ParseDuration(x)
			if err != nil {
				return err
			}

			if config.Timeout != nil {
				if *config.Timeout == d {
					zap.L().Debug(fmt.Sprintf("Timeout %s in config is the same as from arg", d.String()))
					return nil
				}
				zap.L().Debug(fmt.Sprintf("Timeout %s overwritten with %s from args", config.Timeout.String(), d.String()))
				config.Timeout = &d
				return nil
			}

			zap.L().Debug(fmt.Sprintf("Timeout is %s from args", x))
			config.Timeout = &d
			return nil
		}

		x = os.Getenv(ShellyTimeoutEnvVar)
		if x != "" {

			d, err := time.ParseDuration(x)
			if err != nil {
				return err
			}

			if config.Timeout != nil {
				if *config.Timeout == d {
					zap.L().Debug(fmt.Sprintf("Timeout %s in config is the same as from envvar %s", d.String(), ShellyTimeoutEnvVar))
					return nil
				}
				zap.L().Debug(fmt.Sprintf("Timeout %s overwritten with %s from envvar %s", config.Timeout.String(), d.String(), ShellyTimeoutEnvVar))
				config.Timeout = &d
				return nil
			}

			zap.L().Debug(fmt.Sprintf("Timeout is %s from envvar %s", x, ShellyTimeoutEnvVar))
			config.Timeout = &d
			return nil
		}

		if config.Timeout != nil {
			zap.L().Debug(fmt.Sprintf("Timeout is %s from config", config.Timeout.String()))
			return nil
		}

		zap.L().Debug(fmt.Sprintf("Timeout not set in arg, envvar %s or config", ShellyTimeoutEnvVar))

		return nil
	}

	loadUnifi := func() error {

		if config.Unifi == nil {
			zap.L().Debug("unifi config is not present")
			return nil
		}

		unifiArg := strings.ToLower(t.unifiArg)

		switch unifiArg {

		case "":
			if config.Unifi.Enabled {
				zap.L().Debug("unifi config is enabled by config")
			} else {
				zap.L().Debug("unifi config is disabled by config")
			}

		case "enable", "enabled", "true":
			if config.Unifi.Enabled {
				zap.L().Debug("unifi config is enabled by arg and config")
			} else {
				zap.L().Debug("unifi config is disabled by config but enabled by arg")
				config.Unifi.Enabled = true
			}

		case "disable", "disabled", "false":
			if config.Unifi.Enabled {
				zap.L().Debug("unifi is enabled in config but disabled by arg")
				config.Unifi.Enabled = false
			} else {
				zap.L().Debug("unifi config is disabled by arg and config")
			}

		default:
			return fmt.Errorf("unifiArg value of %s is not valid, expecting enable, true, disable, or false", t.unifiArg)

		}

		if !config.Unifi.Enabled {
			zap.L().Debug("unifi config is present but disabled")
			return nil
		}

		zap.L().Debug("unifi config is present and enabled")

		hostnames, err := util.GetUnifiHostnames(config.Unifi)
		if err != nil {
			return err
		}

		config.Hostnames = append(config.Hostnames, hostnames...)

		return nil
	}

	loadOpenhab := func() error {

		if config.OpenHAB == nil {
			zap.L().Debug("openHAB config is not present")
			return nil
		}

		openhabArg := strings.ToLower(t.openhabArg)

		switch openhabArg {

		case "":
			if config.OpenHAB.Enabled {
				zap.L().Debug("openHAB config is enabled by config")
			} else {
				zap.L().Debug("openHAB config is disabled by config")
			}

		case "enable", "enabled", "true":
			if config.OpenHAB.Enabled {
				zap.L().Debug("openHAB config is enabled by arg and config")
			} else {
				zap.L().Debug("openHAB config is disabled by config but enabled by arg")
				config.OpenHAB.Enabled = true
			}

		case "disable", "disabled", "false":
			if config.Unifi.Enabled {
				zap.L().Debug("openHAB is enabled in config but disabled by arg")
				config.OpenHAB.Enabled = false
			} else {
				zap.L().Debug("openHAB config is disabled by arg and config")
			}

		default:
			return fmt.Errorf("openHAB value of %s is not valid, expecting enable, true, disable, or false", t.openhabArg)

		}

		if !config.OpenHAB.Enabled {
			zap.L().Debug("openHAB config is present but disabled")
			return nil
		}

		zap.L().Debug("openHAB config is present and enabled")

		hostnames, err := util.GetOpenHABHostnames(ctx, config.OpenHAB)
		if err != nil {
			return err
		}

		config.Hostnames = append(config.Hostnames, hostnames...)

		return nil
	}

	if err := initLog(); err != nil {
		return nil, err
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		err = initFromBytes(content)
		if err != nil {
			return nil, err
		}

	} else {
		err = initFromFile(getConfigFile())
		if err != nil {
			return nil, err
		}
	}

	if err := loadBase(); err != nil {
		return nil, err
	}

	loadHostnames()
	loadShelly()
	loadOutput()
	loadUpdateURL()

	if err := loadTimeout(); err != nil {
		return nil, err
	}

	if err := loadUnifi(); err != nil {
		return nil, err
	}

	if err := loadOpenhab(); err != nil {
		return nil, err
	}

	return config, nil
}

func ExampleConfig() *Config {

	output := "pretty-json"
	timeout := time.Second * 60

	notes := "Depending on the command Hostname may be required. For commands that\n"
	notes += "use Hostnames hostname will be prepended if it is set. If Unifi\n"
	notes += "config is present then hostnames will be loaded from Unifi. Shelly\n"
	notes += "configs can be specified in the map. The name should be the device\n"
	notes += "ID or App. Device ID takes precedence\n"

	unifiConfig := unifi.ExampleConfig()
	unifiConfig.Enabled = true

	return &Config{
		Notes:   notes,
		Shelly:  sdk_client.ExampleConfig(),
		Unifi:   unifiConfig,
		Output:  &output,
		Timeout: &timeout,
		Logger: &Logger{
			LogLevel: logger.DebugLevel,
		},
		Mqtt:    mqtt.ExampleConfig(),
		OpenHAB: openhab.ExampleConfig(),
	}
}
