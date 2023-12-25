package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/jodydadescott/shelly-client/cmd/mqtt/types"
	cmdtypes "github.com/jodydadescott/shelly-client/cmd/types"
	"github.com/jodydadescott/shelly-client/cmd/util"
	sdk_client "github.com/jodydadescott/shelly-client/sdk/client"
	sdk_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	light_types "github.com/jodydadescott/shelly-client/sdk/light/types"
	shelly_types "github.com/jodydadescott/shelly-client/sdk/shelly/types"
)

type LightStatus = light_types.Status
type Config = cmdtypes.Config
type ShellyClient = sdk_client.Client
type ShellyDeviceInfo = shelly_types.DeviceInfo
type ShellyDeviceStatus = shelly_types.Status
type ShellyConfig = sdk_types.Config
type MqttStatus = types.MqttStatus
type MqttConfig = types.Config
type Params = types.Params

type callback interface {
	GetConfig(context.Context) (*Config, error)
	GetCTX() (context.Context, context.CancelFunc)
	WriteStdout(input any) error
}

func New(t callback) *cobra.Command {

	action := "mqtt poll"

	loadMQTT := func(config *MqttConfig) error {

		if config == nil {
			zap.L().Debug("MQTT config is not present")
			return nil
		}

		if config.Broker == "" {
			return fmt.Errorf("broker is required")
		}

		if config.Topic == "" {
			return fmt.Errorf("topic is required")
		}

		if config.ClientID != "" {
			zap.L().Debug(fmt.Sprintf("clientID is %s (config)", config.ClientID))
		} else {
			config.ClientID = defaultClientID
			zap.L().Debug(fmt.Sprintf("clientID is %s (default)", config.ClientID))
		}

		if config.LightSource != "" {
			zap.L().Debug(fmt.Sprintf("lightSource is %s (config)", config.LightSource))
		} else {
			config.LightSource = defaultLightSource
			zap.L().Debug(fmt.Sprintf("lightSource is %s (default)", config.LightSource))
		}

		if config.KeepAlive > 0 {
			zap.L().Debug(fmt.Sprintf("keepAlive is %s (config)", config.KeepAlive.String()))
		} else {
			config.KeepAlive = defaultKeepAlive
			zap.L().Debug(fmt.Sprintf("keepAlive is %s (default)", config.KeepAlive.String()))
		}

		if config.PingTimeout > 0 {
			zap.L().Debug(fmt.Sprintf("pingTimeout is %s (config)", config.PingTimeout.String()))
		} else {
			config.PingTimeout = defaultPingTimeout
			zap.L().Debug(fmt.Sprintf("pingTimeout is %s (default)", config.PingTimeout.String()))
		}

		if config.PublishTimeout > 0 {
			zap.L().Debug(fmt.Sprintf("publishTimeout is %s (config)", config.PublishTimeout.String()))
		} else {
			config.PublishTimeout = defaultPublishTimeout
			zap.L().Debug(fmt.Sprintf("publishTimeout is %s (default)", config.PublishTimeout.String()))
		}

		if config.DaemonInterval > 0 {
			zap.L().Debug(fmt.Sprintf("daemonInterval is %s (config)", config.DaemonInterval.String()))
		} else {
			config.DaemonInterval = defaultDaemonInterval
			zap.L().Debug(fmt.Sprintf("daemonInterval is %s (default)", config.DaemonInterval.String()))
		}

		return nil
	}

	process := func(ctx context.Context, config *Config) error {

		zap.L().Debug("process")

		opts := mqtt.NewClientOptions().AddBroker(config.Mqtt.Broker).SetClientID(config.Mqtt.ClientID)
		opts.SetKeepAlive(2 * time.Second)
		opts.SetPingTimeout(1 * time.Second)

		mqttClient := mqtt.NewClient(opts)

		if token := mqttClient.Connect(); token.WaitTimeout(time.Second*30) && token.Error() != nil {
			return token.Error()
		}

		do := func(ctx context.Context, hostname string, client *ShellyClient, deviceInfo *ShellyDeviceInfo, deviceStatus *ShellyDeviceStatus) error {

			if deviceInfo.App == nil {
				return fmt.Errorf("deviceInfo.App is nil")
			}

			if *deviceInfo.App != shellyPlusWallDimmer {
				if logger.Trace {
					zap.L().Debug(fmt.Sprintf("hostname %s, deviceInfoID %s, deviceInfoApp %s, [ignored]", hostname, *deviceInfo.ID, *deviceInfo.App))
				}
				return nil
			}

			newResults := func(topic string) []*MqttStatus {

				ts := float64(time.Now().Unix())

				var results []*MqttStatus

				for i, light := range deviceStatus.Light {

					m := &MqttStatus{
						Src:    *deviceInfo.ID,
						Dst:    topic,
						Method: "NotifyStatus",
						Params: &Params{},
					}

					m.Params.Ts = &ts
					light.Source = &config.Mqtt.LightSource

					switch i {

					case 0:
						m.Params.Light0 = light
					case 1:
						m.Params.Light1 = light
					case 2:
						m.Params.Light2 = light
					case 3:
						m.Params.Light3 = light
					case 4:
						m.Params.Light4 = light
					case 5:
						m.Params.Light5 = light
					case 6:
						m.Params.Light6 = light
					case 7:
						m.Params.Light7 = light

					}

					results = append(results, m)
				}

				return results
			}

			for _, result := range newResults(config.Mqtt.Topic) {
				b, _ := json.Marshal(result)
				msg := string(b)

				go func() {
					token := mqttClient.Publish(config.Mqtt.Topic, 0, false, msg)

					if token.WaitTimeout(time.Second * 30) {
						if logger.Trace {
							zap.L().Debug(fmt.Sprintf("Publish success; message %s, topic %s", msg, config.Mqtt.Topic))
						}
					} else {
						zap.L().Error(fmt.Sprintf("Publish fail; message %s, topic %s", msg, config.Mqtt.Topic))
					}

				}()

			}
			return nil
		}

		return util.Process(ctx, config, action, true, do)
	}

	rootCmd := &cobra.Command{
		Use:   "mqtt",
		Short: "Polls each shelly light device and publishes status to specified MQTT broker",
	}

	mqttOnceCmd := &cobra.Command{
		Use:   "once",
		Short: "Runs MQTT polling once and then exits",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.Mqtt == nil {
				return fmt.Errorf("mqtt is required")
			}

			if err := loadMQTT(config.Mqtt); err != nil {
				return err
			}

			return process(ctx, config)
		},
	}

	mqttDaemonCmd := &cobra.Command{
		Use:   "daemon",
		Short: "Runs MQTT polling as a daemon",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := t.GetCTX()
			defer cancel()

			config, err := t.GetConfig(ctx)
			if err != nil {
				return err
			}

			if config.Mqtt == nil {
				return fmt.Errorf("mqtt is required")
			}

			if err := loadMQTT; err != nil {
				return nil
			}

			run := func() error {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()
				return process(ctx, config)
			}

			ticker := time.NewTicker(config.Mqtt.DaemonInterval)

			for {

				select {

				case <-ctx.Done():
					return nil

				case <-ticker.C:
					err := run()
					if err != nil {
						if config.Mqtt.FailOnError {
							return err
						} else {
							zap.L().Error(fmt.Sprintf("failOnError is false; error is %s", err.Error()))
						}
					}

				}

			}

		},
	}

	rootCmd.AddCommand(mqttOnceCmd, mqttDaemonCmd)
	return rootCmd
}

func ExampleConfig() *MqttConfig {
	return &MqttConfig{
		Notes:          "MQTT config is only required if running MQTT polling",
		Broker:         "mybroker:1883",
		Topic:          "my/topic",
		Username:       "username",
		Password:       "password",
		ClientID:       "clientID",
		KeepAlive:      defaultKeepAlive,
		PingTimeout:    defaultPingTimeout,
		PublishTimeout: defaultPublishTimeout,
	}
}
