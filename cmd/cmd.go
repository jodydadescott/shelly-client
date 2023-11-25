package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hokaccha/go-prettyjson"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-client/cmd/files"
	lightcmd "github.com/jodydadescott/shelly-client/cmd/light"
	shellycmd "github.com/jodydadescott/shelly-client/cmd/shelly"
	switchxcmd "github.com/jodydadescott/shelly-client/cmd/switchx"
	wificmd "github.com/jodydadescott/shelly-client/cmd/wifi"
	"github.com/jodydadescott/shelly-client/sdk"
	"github.com/jodydadescott/shelly-client/sdk/input"
	"github.com/jodydadescott/shelly-client/sdk/light"
	"github.com/jodydadescott/shelly-client/sdk/shelly"
	"github.com/jodydadescott/shelly-client/sdk/switchx"
	"github.com/jodydadescott/shelly-client/sdk/system"
	"github.com/jodydadescott/shelly-client/sdk/wifi"
)

type Cmd struct {
	*cobra.Command
	_client       *sdk.Client
	hostnameArg   string
	passwordArg   string
	outputArg     string
	configArg     string
	timeoutArg    int
	debugLevelArg string
}

func NewCmd() *Cmd {

	t := &Cmd{}

	command := &cobra.Command{

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

		PersistentPostRun: func(cmd *cobra.Command, args []string) {

			if t._client != nil {
				t._client.Close()
			}

		},
	}

	t.Command = command

	t.PersistentFlags().StringVarP(&t.hostnameArg, "hostname", "H", "", fmt.Sprintf("Hostname; optionally use env var '%s'", ShellyHostnameEnvVar))
	t.PersistentFlags().StringVarP(&t.hostnameArg, "password", "p", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyPasswordEnvVar))
	t.PersistentFlags().StringVarP(&t.outputArg, "output", "o", ShellyOutputDefault, fmt.Sprintf("Output format. One of: prettyjson | json | jsonpath | yaml ; Optionally use env var '%s'", ShellyOutputEnvVar))
	t.PersistentFlags().StringVarP(&t.configArg, "config", "c", "", "Config file or Directory name. Directory will look for file name device-id then app-id")
	t.PersistentFlags().IntVarP(&t.timeoutArg, "timeout", "t", 0, "The timeout in seconds for the websocket call to the device")
	t.PersistentFlags().StringVarP(&t.debugLevelArg, "debug", "D", "", fmt.Sprintf("debug level (TRACE, DEBUG, INFO, WARN, ERROR) to STDERR; env var is %s", DebugEnvVar))

	t.AddCommand(shellycmd.NewCmd(t), wificmd.NewCmd(t), switchxcmd.NewCmd(t), lightcmd.NewCmd(t))

	return t
}

func (t *Cmd) client() (*sdk.Client, error) {

	if t._client != nil {
		return t._client, nil
	}

	config := &sdk.Config{
		Hostname: t.hostnameArg,
		Password: t.passwordArg,
	}

	if t.timeoutArg > 0 {
		config.SendTimeout = time.Duration(t.timeoutArg) * time.Second
	}

	if config.Hostname == "" {
		config.Hostname = os.Getenv(ShellyHostnameEnvVar)
	}

	if config.Password == "" {
		config.Password = os.Getenv(ShellyPasswordEnvVar)
	}

	client, err := sdk.New(config)
	if err != nil {
		return nil, err
	}

	t._client = client

	return t._client, nil
}

func (t *Cmd) System() (*system.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.System(), nil
}

func (t *Cmd) Shelly() (*shelly.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.Shelly(), nil
}

func (t *Cmd) Wifi() (*wifi.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.Wifi(), nil
}

func (t *Cmd) Switch() (*switchx.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.Switch(), nil
}

func (t *Cmd) Light() (*light.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.Light(), nil
}

func (t *Cmd) Input() (*input.Client, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	return client.Input(), nil
}

func (t *Cmd) RebootDevice(ctx context.Context) error {
	client, err := t.client()
	if err != nil {
		return err
	}
	return client.Shelly().Reboot(ctx)
}

// WriteObject writes object in desired format to STDOUT
func (t *Cmd) WriteStdout(input any) error {

	switch v := input.(type) {

	case nil:
		return nil

	case string:
		fmt.Println(v)
		return nil

	case *string:
		fmt.Println(*v)
		return nil

	}

	// Arg can be in the format of 'format' or 'format=...'.
	// For example jsonpath expects a second arg such as jsonpath=$.

	outputArgSplit := strings.Split(t.outputArg, "=")

	switch strings.ToLower(outputArgSplit[0]) {

	case "prettyjson":
		data, err := prettyjson.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	case "jsonpath":

		if len(outputArgSplit) > 1 {

			v := interface{}(nil)
			data, err := json.Marshal(input)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &v)
			if err != nil {
				return err
			}

			data2, err := jsonpath.Get(outputArgSplit[1], v)
			if err != nil {
				return err
			}

			switch data2.(type) {
			case string:
				fmt.Println(data2)

			case []interface{}:
				s := reflect.ValueOf(data2)
				for i := 0; i < s.Len(); i++ {
					fmt.Println(s.Index(i))
				}
			}

			return nil
		}

		return fmt.Errorf("missing jsonpath arg. Expect jsonpath")

	case "json":
		data, err := json.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	case "yaml":
		data, err := yaml.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	}

	return fmt.Errorf("format type %s is unknown", t.outputArg)
}

func (t *Cmd) WriteStderr(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func (t *Cmd) GetFiles() (*files.Files, error) {
	return files.NewFiles(t.configArg)
}
