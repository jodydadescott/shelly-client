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
	"github.com/jodydadescott/unifi-go-sdk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-client/cmd/light"
	"github.com/jodydadescott/shelly-client/cmd/switchx"
	"github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdkclient "github.com/jodydadescott/shelly-client/sdk"
	"github.com/jodydadescott/shelly-client/sdk/shelly"
	sdktypes "github.com/jodydadescott/shelly-client/sdk/types"
)

type ShellyDeviceInfo = sdktypes.ShelllyDeviceInfo
type ShellyDeviceStatus = sdktypes.ShellyStatus
type ShellyConfig = sdktypes.ShellyConfig

type ShellyClient = sdkclient.Client

type Config = types.Config

type UnifiConfig = unifi.Config

type Cmd struct {
	ctx context.Context
	*cobra.Command
	configFileArg  string
	hostnameArg    []string
	passwordArg    string
	outputArg      string
	urlArg         string
	timeoutArg     string
	debugLevelArg  string
	unifiArg       string
	rebootForceArg bool
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

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Returns specified config, status or info for device",
	}

	getConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Returns device config",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, config.Hostnames[0]).Shelly()
			defer client.Close()

			ctx, cancel := t.GetCTX()
			defer cancel()

			result, err := client.GetConfig(ctx)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	getStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Returns device status",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, config.Hostnames[0]).Shelly()
			defer client.Close()

			ctx, cancel := t.GetCTX()
			defer cancel()

			result, err := client.GetStatus(ctx)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	getInfoCmd := &cobra.Command{
		Use:   "info",
		Short: "Returns device info",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, config.Hostnames[0]).Shelly()
			defer client.Close()

			ctx, cancel := t.GetCTX()
			defer cancel()

			result, err := client.GetDeviceInfo(ctx)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	getMethodsCmd := &cobra.Command{
		Use:   "methods",
		Short: "Returns RPC methods for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, config.Hostnames[0]).Shelly()
			defer client.Close()

			ctx, cancel := t.GetCTX()
			defer cancel()

			result, err := client.ListMethods(ctx)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	getUpdatesCmd := &cobra.Command{
		Use:   "updates",
		Short: "Returns available firmware updates for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) != 1 {
				return fmt.Errorf("one and only one hostname is required for this command")
			}

			client := util.NewShellyClient(config, config.Hostnames[0]).Shelly()
			defer client.Close()

			ctx, cancel := t.GetCTX()
			defer cancel()

			result, err := client.CheckForUpdate(ctx)
			if err != nil {
				return err
			}

			return t.WriteStdout(result)
		},
	}

	getCmd.AddCommand(getConfigCmd, getStatusCmd, getInfoCmd, getMethodsCmd, getUpdatesCmd)

	unifiCmd := &cobra.Command{
		Use:   "unifi",
		Short: "Ubiquiti Networks",
	}

	unifiGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Returns list of Shelly device hostname(s)",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if config.Unifi == nil {
				return fmt.Errorf("config.Unifi is required")
			}

			hostnames, err := getUnifiHostnames(config.Unifi)
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

	unifiCmd.AddCommand(unifiGetCmd)

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates firmware",
	}

	updateStableCmd := &cobra.Command{
		Use:   "stable",
		Short: "Updates devices(s) firmware to latest stable",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			stage := "beta"
			action := fmt.Sprintf("update to %s", stage)

			params := &shelly.ShellyUpdateConfig{
				Stage: &stage,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				err := client.Shelly().Update(ctx, params)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	updateBetaCmd := &cobra.Command{
		Use:   "beta",
		Short: "Updates devices(s) firmware to latest stable",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			stage := "beta"
			action := fmt.Sprintf("update to %s", stage)

			params := &shelly.ShellyUpdateConfig{
				Stage: &stage,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				err := client.Shelly().Update(ctx, params)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	updateURLCmd := &cobra.Command{
		Use:   "url",
		Short: "Updates device(s) firmware to specified URL",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
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

			params := &shelly.ShellyUpdateConfig{
				Url: &url,
			}

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				err := client.Shelly().Update(ctx, params)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	updateCmd.AddCommand(updateStableCmd, updateBetaCmd, updateURLCmd)

	doCmd := &cobra.Command{
		Use:  "do",
		Long: "Executes actions on specified device(s)",
	}

	rebootCmd := &cobra.Command{
		Use:   "reboot",
		Short: "Reboots device(s)",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			action := "reboot"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

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
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] not required", hostname, *deviceInfo.ID, *deviceInfo.App, action))
					return nil
				}

				err := client.Shelly().Reboot(ctx)
				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, true, do)
		},
	}

	rebootCmd.PersistentFlags().BoolVarP(&t.rebootForceArg, "force", "f", false, "reboots device even if device reboot is not required")

	// Need a sanity flag like force TODO
	factoryResetCmd := &cobra.Command{
		Use:   "factory",
		Short: "Resets device(s) to factory settings",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			action := "factory reset"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				err := client.Shelly().FactoryReset(ctx)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil

			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	// Need a sanity flag like force TODO
	resetWifiConfigCmd := &cobra.Command{
		Use:   "wifi",
		Short: "Resets device(s) Wifi config",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			action := "reset wifi"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				err := client.Shelly().ResetWiFiConfig(ctx)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil

			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	doCmd.AddCommand(rebootCmd, factoryResetCmd, resetWifiConfigCmd)

	setCmd := &cobra.Command{
		Use: "set",
	}

	setConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Sets specified config for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			action := "set config"

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

				getConfig := func() *ShellyConfig {

					shellyConfig := config.GetMatchingShellyConfig(*deviceInfo.ID)
					if shellyConfig != nil {
						zap.L().Debug(fmt.Sprintf("Using config %s based on matching ID", *deviceInfo.ID))
						return shellyConfig
					}

					shellyConfig = config.GetMatchingShellyConfig(*deviceInfo.App)
					if shellyConfig != nil {
						zap.L().Debug(fmt.Sprintf("Using config %s based on matching App", *deviceInfo.App))
						return shellyConfig
					}

					shellyConfig = config.GetMatchingShellyConfig("default")
					if shellyConfig != nil {
						zap.L().Debug("Using default config")
						return shellyConfig
					}

					return nil
				}

				shellyConfig := getConfig()
				if shellyConfig == nil {
					return fmt.Errorf("no matching config found")
				}

				err := client.Shelly().SetConfig(ctx, shellyConfig, deviceInfo)

				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))
				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			return util.Process(ctx, config, action, false, do)
		},
	}

	exampleCmd := &cobra.Command{
		Use:   "example",
		Short: "writes example to STDOUT",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := t.GetConfig()
			if err != nil {
				return err
			}

			if len(config.Hostnames) == 0 {
				return fmt.Errorf("one or more hostnames required")
			}

			action := "example"

			result := types.ExampleConfig()

			do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {
				deviceConfig, err := client.Shelly().GetConfig(ctx)
				if err != nil {
					t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] failed with error %s", hostname, *deviceInfo.ID, *deviceInfo.App, action, err.Error()))
					return err
				}

				t.WriteStdout(fmt.Sprintf("hostname %s, deviceID %s, deviceApp %s: [%s] completed", hostname, *deviceInfo.ID, action, *deviceInfo.App))

				result.AddShellyConfig(*deviceInfo.App, deviceConfig)
				return nil
			}

			ctx, cancel := t.GetCTX()
			defer cancel()

			defer t.WriteStdout(result)

			return util.Process(ctx, config, action, false, do)
		},
	}

	setCmd.AddCommand(setConfigCmd)

	// Override the default help so we can use the -h for host
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	rootCmd.PersistentFlags().StringVarP(&t.debugLevelArg, "debug", "D", "", fmt.Sprintf("debug level (TRACE, DEBUG, INFO, WARN, ERROR) to STDERR; env var is %s", DebugEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.configFileArg, "config", "c", "", fmt.Sprintf("Config; optionally use env var '%s'", ShellyConfigEnvVar))
	rootCmd.PersistentFlags().StringSliceVarP(&t.hostnameArg, "hostname", "h", []string{}, fmt.Sprintf("Hostname; optionally use env var '%s'", ShellyHostnameEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.passwordArg, "password", "p", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyPasswordEnvVar))
	rootCmd.PersistentFlags().StringVar(&t.passwordArg, "update-url", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyURLEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.outputArg, "output", "o", ShellyOutputDefault, fmt.Sprintf("Output format. One of: prettyjson | json | jsonpath | yaml ; Optionally use env var '%s'", ShellyOutputEnvVar))
	rootCmd.PersistentFlags().StringVarP(&t.timeoutArg, "timeout", "t", "", "The timeout in seconds for the websocket call to the device")
	rootCmd.PersistentFlags().StringVar(&t.unifiArg, "unifi", "", "Use Unifi controller to get hostname(s). Must be enable or disable")
	rootCmd.AddCommand(getCmd, setCmd, doCmd, updateCmd, exampleCmd, unifiCmd, light.New(t), switchx.New(t))

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
		fmt.Println(strings.TrimSpace(string(data)))
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

func (t *Cmd) GetConfig() (*Config, error) {

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
					config.Hostnames = t.hostnameArg
					return
				}

				zap.L().Debug(fmt.Sprintf("hostnames is %s from envvar %s", hostnames, ShellyHostnameEnvVar))
				config.Hostnames = t.hostnameArg
				return
			}

		}

		if len(config.Hostnames) > 0 {
			zap.L().Debug(fmt.Sprintf("hostnames is %s from config", config.Hostnames))
			return
		}

		zap.L().Debug(fmt.Sprintf("Hostnames not set in arg, envvar %s or config", ShellyHostnameEnvVar))
	}

	loadPassword := func() {

		x := t.passwordArg
		if x != "" {
			if config.Password != nil {
				if *config.Password == x {
					zap.L().Debug(fmt.Sprintf("Password %s in config is the same as from arg", x))
					return
				}
				zap.L().Debug(fmt.Sprintf("Password %s overwritten with %s from args", *config.Password, x))
				config.Password = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("Password is %s from args", x))
			config.Password = &x
			return
		}

		x = os.Getenv(ShellyPasswordEnvVar)
		if x != "" {
			if config.Password != nil {
				if *config.Password == x {
					zap.L().Debug(fmt.Sprintf("Password %s in config is the same as from envvar %s", x, ShellyPasswordEnvVar))
					return
				}
				zap.L().Debug(fmt.Sprintf("Password %s overwritten with %s from envvar %s", *config.Password, x, ShellyPasswordEnvVar))
				config.Password = &x
				return
			}

			zap.L().Debug(fmt.Sprintf("Password is %s from envvar %s", x, ShellyPasswordEnvVar))
			config.Password = &x
			return
		}

		if config.Password != nil {
			zap.L().Debug(fmt.Sprintf("Password is %s from config", *config.Password))
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
			zap.L().Debug("Unifi config is not present")
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
			zap.L().Debug("Unifi config is present but diabled")
			return nil
		}

		zap.L().Debug("Unifi config is present and enabled")

		hostnames, err := getUnifiHostnames(config.Unifi)
		if err != nil {
			return err
		}

		config.Hostnames = append(config.Hostnames, hostnames...)
		return nil
	}

	err := initLog()
	if err != nil {
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

	loadHostnames()
	loadPassword()
	loadOutput()
	loadUpdateURL()
	err = loadTimeout()
	if err != nil {
		return nil, err
	}

	err = loadUnifi()
	if err != nil {
		return nil, err
	}

	if t.passwordArg != "" {
		config.Password = &t.passwordArg
		zap.L().Debug("Password loaded from args")
	}

	if t.passwordArg != "" {
		config.Password = &t.passwordArg
		zap.L().Debug("Password loaded from args")
	}

	if t.outputArg != "" {
		config.Output = &t.outputArg
		zap.L().Debug(fmt.Sprintf("Output %s loaded from args", t.outputArg))
	}

	if t.timeoutArg != "" {
		timeout, err := time.ParseDuration(t.timeoutArg)
		if err != nil {
			return nil, err
		}
		config.Timeout = &timeout
		zap.L().Debug(fmt.Sprintf("Timeout %s loaded from args", t.timeoutArg))
	}

	return config, nil
}

func getUnifiHostnames(config *UnifiConfig) ([]string, error) {

	normalizeHostname := func(input string) string {
		input = strings.ToLower(input)
		input = space.ReplaceAllString(input, "-")
		return strings.Split(input, ".")[0]
	}

	hasShellyPrefix := func(input string) bool {
		input = strings.ToLower(input)
		for _, shellyPrefix := range knownShellyHostnamePrefixes {
			if strings.HasPrefix(input, shellyPrefix) {
				return true
			}
		}

		return false
	}

	unifiClient := unifi.New(config)

	clients, err := unifiClient.GetClients()
	if err != nil {
		return nil, err
	}

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
