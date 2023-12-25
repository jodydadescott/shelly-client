package types

import (
	bluetooth_types "github.com/jodydadescott/shelly-client/sdk/bluetooth/types"
	cloud_types "github.com/jodydadescott/shelly-client/sdk/cloud/types"
	ethernet_types "github.com/jodydadescott/shelly-client/sdk/ethernet/types"
	input_types "github.com/jodydadescott/shelly-client/sdk/input/types"
	light_types "github.com/jodydadescott/shelly-client/sdk/light/types"
	mqtt_types "github.com/jodydadescott/shelly-client/sdk/mqtt/types"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	switch_types "github.com/jodydadescott/shelly-client/sdk/switchx/types"
	system_types "github.com/jodydadescott/shelly-client/sdk/system/types"
	websocket_types "github.com/jodydadescott/shelly-client/sdk/websocket/types"
	wifi_types "github.com/jodydadescott/shelly-client/sdk/wifi/types"
)

type Request = msg_types.Request
type Response = msg_types.Response
type Error = msg_types.Error

type BluetoothStatus = bluetooth_types.Status
type BluetoothConfig = bluetooth_types.Config

type CloudStatus = cloud_types.Status
type CloudConfig = cloud_types.Config

type MqttStatus = mqtt_types.Status
type MqttConfig = mqtt_types.Config

type EthernetStatus = ethernet_types.Status
type EthernetConfig = ethernet_types.Config

type SystemStatus = system_types.Status
type SystemConfig = system_types.Config

type WifiStatus = wifi_types.Status
type WifiConfig = wifi_types.Config

type WebsocketStatus = websocket_types.Status
type WebsocketConfig = websocket_types.Config

type LightStatus = light_types.Status
type LightConfig = light_types.Config

type InputStatus = input_types.Status
type InputConfig = input_types.Config

type SwitchStatus = switch_types.Status
type SwitchConfig = switch_types.Config

type SystemAvailableUpdates = system_types.SystemAvailableUpdates
