package types

// RawShellyStatus internal use only
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly
type RawShellyStatus struct {
	Bluetooth *BluetoothStatus `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud     *CloudStatus     `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt      *MqttStatus      `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet  *EthernetStatus  `json:"eth,omitempty" yaml:"eth,omitempty"`
	System    *SystemStatus    `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi      *WifiStatus      `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Websocket *WebsocketStatus `json:"ws,omitempty" yaml:"ws,omitempty"`
	Light0    *LightStatus     `json:"light:0,omitempty" yaml:"light:0,omitempty"`
	Light1    *LightStatus     `json:"light:1,omitempty" yaml:"light:1,omitempty"`
	Light2    *LightStatus     `json:"light:2,omitempty" yaml:"light:2,omitempty"`
	Light3    *LightStatus     `json:"light:3,omitempty" yaml:"light:3,omitempty"`
	Light4    *LightStatus     `json:"light:4,omitempty" yaml:"light:4,omitempty"`
	Light5    *LightStatus     `json:"light:5,omitempty" yaml:"light:5,omitempty"`
	Light6    *LightStatus     `json:"light:6,omitempty" yaml:"light:6,omitempty"`
	Light7    *LightStatus     `json:"light:7,omitempty" yaml:"light:7,omitempty"`
	Input0    *InputStatus     `json:"input:0,omitempty" yaml:"input:0,omitempty"`
	Input1    *InputStatus     `json:"input:1,omitempty" yaml:"input:1,omitempty"`
	Input2    *InputStatus     `json:"input:2,omitempty" yaml:"input:2,omitempty"`
	Input3    *InputStatus     `json:"input:3,omitempty" yaml:"input:3,omitempty"`
	Input4    *InputStatus     `json:"input:4,omitempty" yaml:"input:4,omitempty"`
	Input5    *InputStatus     `json:"input:5,omitempty" yaml:"input:5,omitempty"`
	Input6    *InputStatus     `json:"input:6,omitempty" yaml:"input:6,omitempty"`
	Input7    *InputStatus     `json:"input:7,omitempty" yaml:"input:7,omitempty"`
	Switch0   *SwitchStatus    `json:"switch:0,omitempty" yaml:"switch:0,omitempty"`
	Switch1   *SwitchStatus    `json:"switch:1,omitempty" yaml:"switch:1,omitempty"`
	Switch2   *SwitchStatus    `json:"switch:2,omitempty" yaml:"switch:2,omitempty"`
	Switch3   *SwitchStatus    `json:"switch:3,omitempty" yaml:"switch:3,omitempty"`
	Switch4   *SwitchStatus    `json:"switch:4,omitempty" yaml:"switch:4,omitempty"`
	Switch5   *SwitchStatus    `json:"switch:5,omitempty" yaml:"switch:5,omitempty"`
	Switch6   *SwitchStatus    `json:"switch:6,omitempty" yaml:"switch:6,omitempty"`
	Switch7   *SwitchStatus    `json:"switch:7,omitempty" yaml:"switch:7,omitempty"`
}

func (t *RawShellyStatus) Convert() *Status {

	c := &Status{
		Bluetooth: t.Bluetooth,
		Cloud:     t.Cloud,
		Mqtt:      t.Mqtt,
		Ethernet:  t.Ethernet,
		System:    t.System,
		Wifi:      t.Wifi,
		Light:     make(map[int]*LightStatus),
		Input:     make(map[int]*InputStatus),
		Switch:    make(map[int]*SwitchStatus),
	}

	if t.Light0 != nil {
		c.Light[0] = t.Light0
	}
	if t.Light1 != nil {
		c.Light[1] = t.Light1
	}
	if t.Light2 != nil {
		c.Light[2] = t.Light2
	}
	if t.Light3 != nil {
		c.Light[3] = t.Light3
	}
	if t.Light4 != nil {
		c.Light[4] = t.Light4
	}
	if t.Light5 != nil {
		c.Light[5] = t.Light5
	}
	if t.Light6 != nil {
		c.Light[6] = t.Light6
	}
	if t.Light7 != nil {
		c.Light[7] = t.Light7
	}

	if t.Input0 != nil {
		c.Input[0] = t.Input0
	}
	if t.Input1 != nil {
		c.Input[1] = t.Input1
	}
	if t.Input2 != nil {
		c.Input[2] = t.Input2
	}
	if t.Input3 != nil {
		c.Input[3] = t.Input3
	}
	if t.Input4 != nil {
		c.Input[4] = t.Input4
	}
	if t.Input5 != nil {
		c.Input[5] = t.Input5
	}
	if t.Input6 != nil {
		c.Input[6] = t.Input6
	}
	if t.Input7 != nil {
		c.Input[7] = t.Input7
	}

	if t.Switch0 != nil {
		c.Switch[0] = t.Switch0
	}
	if t.Switch1 != nil {
		c.Switch[1] = t.Switch1
	}
	if t.Switch2 != nil {
		c.Switch[2] = t.Switch2
	}
	if t.Switch3 != nil {
		c.Switch[3] = t.Switch3
	}
	if t.Switch4 != nil {
		c.Switch[4] = t.Switch4
	}
	if t.Switch5 != nil {
		c.Switch[5] = t.Switch5
	}
	if t.Switch6 != nil {
		c.Switch[6] = t.Switch6
	}
	if t.Switch7 != nil {
		c.Switch[7] = t.Switch7
	}

	return c
}

// RawConfig internal use only
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type RawConfig struct {
	Bluetooth *BluetoothConfig `json:"ble,omitempty" yaml:"ble,omitempty"`
	Cloud     *CloudConfig     `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Mqtt      *MqttConfig      `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	Ethernet  *EthernetConfig  `json:"eth,omitempty" yaml:"eth,omitempty"`
	System    *SystemConfig    `json:"sys,omitempty" yaml:"sys,omitempty"`
	Wifi      *WifiConfig      `json:"wifi,omitempty" yaml:"wifi,omitempty"`
	Websocket *WebsocketConfig `json:"ws,omitempty" yaml:"ws,omitempty"`
	Light0    *LightConfig     `json:"light:0,omitempty" yaml:"light:0,omitempty"`
	Light1    *LightConfig     `json:"light:1,omitempty" yaml:"light:1,omitempty"`
	Light2    *LightConfig     `json:"light:2,omitempty" yaml:"light:2,omitempty"`
	Light3    *LightConfig     `json:"light:3,omitempty" yaml:"light:3,omitempty"`
	Light4    *LightConfig     `json:"light:4,omitempty" yaml:"light:4,omitempty"`
	Light5    *LightConfig     `json:"light:5,omitempty" yaml:"light:5,omitempty"`
	Light6    *LightConfig     `json:"light:6,omitempty" yaml:"light:6,omitempty"`
	Light7    *LightConfig     `json:"light:7,omitempty" yaml:"light:7,omitempty"`
	Input0    *InputConfig     `json:"input:0,omitempty" yaml:"input:0,omitempty"`
	Input1    *InputConfig     `json:"input:1,omitempty" yaml:"input:1,omitempty"`
	Input2    *InputConfig     `json:"input:2,omitempty" yaml:"input:2,omitempty"`
	Input3    *InputConfig     `json:"input:3,omitempty" yaml:"input:3,omitempty"`
	Input4    *InputConfig     `json:"input:4,omitempty" yaml:"input:4,omitempty"`
	Input5    *InputConfig     `json:"input:5,omitempty" yaml:"input:5,omitempty"`
	Input6    *InputConfig     `json:"input:6,omitempty" yaml:"input:6,omitempty"`
	Input7    *InputConfig     `json:"input:7,omitempty" yaml:"input:7,omitempty"`
	Switch0   *SwitchConfig    `json:"switch:0,omitempty" yaml:"switch:0,omitempty"`
	Switch1   *SwitchConfig    `json:"switch:1,omitempty" yaml:"switch:1,omitempty"`
	Switch2   *SwitchConfig    `json:"switch:2,omitempty" yaml:"switch:2,omitempty"`
	Switch3   *SwitchConfig    `json:"switch:3,omitempty" yaml:"switch:3,omitempty"`
	Switch4   *SwitchConfig    `json:"switch:4,omitempty" yaml:"switch:4,omitempty"`
	Switch5   *SwitchConfig    `json:"switch:5,omitempty" yaml:"switch:5,omitempty"`
	Switch6   *SwitchConfig    `json:"switch:6,omitempty" yaml:"switch:6,omitempty"`
	Switch7   *SwitchConfig    `json:"switch:7,omitempty" yaml:"switch:7,omitempty"`
}

func (t *RawConfig) Convert() *Config {

	c := &Config{
		Bluetooth: t.Bluetooth,
		Cloud:     t.Cloud,
		Mqtt:      t.Mqtt,
		Ethernet:  t.Ethernet,
		System:    t.System,
		Wifi:      t.Wifi,
		Websocket: t.Websocket,
		Light:     make(map[int]*LightConfig),
		Input:     make(map[int]*InputConfig),
		Switch:    make(map[int]*SwitchConfig),
	}

	if t.Light0 != nil {
		c.Light[0] = t.Light0
	}
	if t.Light1 != nil {
		c.Light[1] = t.Light1
	}
	if t.Light2 != nil {
		c.Light[2] = t.Light2
	}
	if t.Light3 != nil {
		c.Light[3] = t.Light3
	}
	if t.Light4 != nil {
		c.Light[4] = t.Light4
	}
	if t.Light5 != nil {
		c.Light[5] = t.Light5
	}
	if t.Light6 != nil {
		c.Light[6] = t.Light6
	}
	if t.Light7 != nil {
		c.Light[7] = t.Light7
	}

	if t.Input0 != nil {
		c.Input[0] = t.Input0
	}
	if t.Input1 != nil {
		c.Input[1] = t.Input1
	}
	if t.Input2 != nil {
		c.Input[2] = t.Input2
	}
	if t.Input3 != nil {
		c.Input[3] = t.Input3
	}
	if t.Input4 != nil {
		c.Input[4] = t.Input4
	}
	if t.Input5 != nil {
		c.Input[5] = t.Input5
	}
	if t.Input6 != nil {
		c.Input[6] = t.Input6
	}
	if t.Input7 != nil {
		c.Input[7] = t.Input7
	}

	if t.Switch0 != nil {
		c.Switch[0] = t.Switch0
	}
	if t.Switch1 != nil {
		c.Switch[1] = t.Switch1
	}
	if t.Switch2 != nil {
		c.Switch[2] = t.Switch2
	}
	if t.Switch3 != nil {
		c.Switch[3] = t.Switch3
	}
	if t.Switch4 != nil {
		c.Switch[4] = t.Switch4
	}
	if t.Switch5 != nil {
		c.Switch[5] = t.Switch5
	}
	if t.Switch6 != nil {
		c.Switch[6] = t.Switch6
	}
	if t.Switch7 != nil {
		c.Switch[7] = t.Switch7
	}

	return c
}

// ShellyAuthConfig internal use only
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#configuration
type RawAuthConfig struct {
	// User is used by the following methods:
	// SetAuth: Must be set to admin. Only one user is supported. Required
	User *string `json:"user" yaml:"user"`
	// Realm is used by the following methods:
	// SetAuth : Must be the id of the device. Only one realm is supported. Required
	Realm *string `json:"realm" yaml:"realm"`
	// Ha1 is used by the following methods:
	// SetAuth : "user:realm:password" encoded in SHA256 (null to disable authentication). Required
	Ha1 *string `json:"ha1" yaml:"ha1"`
}

// RawTLSConfig internal use only
type RawTLSConfig struct {
	// Data is used by the following methods:
	// PutUserCA : Contents of the PEM file (null if you want to delete the existing data). Required
	// PutTLSClientCert : Contents of the client.crt file (null if you want to delete the existing data). Required
	// PutTLSClientKey : Contents of the client.key file (null if you want to delete the existing data). Required
	Data *string `json:"data" yaml:"data"`
	// Append is used by the following methods:
	// PutUserCA : true if more data will be appended afterwards, default false.
	// PutTLSClientCert : true if more data will be appended afterwards, default false
	// PutTLSClientKey : true if more data will be appended afterwards, default false
	Append *bool `json:"append,omitempty" yaml:"append,omitempty"`
}
