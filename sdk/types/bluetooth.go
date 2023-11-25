package types

import (
	"github.com/jinzhu/copier"
)

// BluetoothStatus status of the BLE component contains information about the bluetooth on/off state and
// does not own any status properties.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#status
type BluetoothStatus struct {
}

// Clone return copy
func (t *BluetoothStatus) Clone() *BluetoothStatus {
	c := &BluetoothStatus{}
	copier.Copy(&c, &t)
	return c
}

// BluetoothConfig configuration of the Bluetooth Low Energy component shows whether the bluetooth connection is enabled.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type BluetoothConfig struct {
	// Enable True if bluetooth is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// RPC configuration of the rpc service
	RPC *BluetoothRPC `json:"rpc,omitempty" yaml:"rpc,omitempty"`
	// Observer configuration of the BT LE observer
	Observer *BluetoothObserver `json:"observer,omitempty" yaml:"observer,omitempty"`
}

// Clone return copy
func (t *BluetoothConfig) Clone() *BluetoothConfig {
	c := &BluetoothConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *BluetoothConfig) Markup() {

	if t == nil {
		return
	}

	if t.Enable {
		t.RPC.Markup()
		t.Observer.Markup()
		return
	}

	t.RPC = nil
	t.Observer = nil
}

// Sanatize sanatize config
func (t *BluetoothConfig) Sanatize() {

	if t == nil {
		return
	}

	t.RPC.Sanatize()
	t.Observer.Sanatize()
}

// BluetoothRPC configuration of the rpc service
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type BluetoothRPC struct {
	// Enable True if rpc service is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
}

// Clone return copy
func (t *BluetoothRPC) Clone() *BluetoothRPC {
	c := &BluetoothRPC{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *BluetoothRPC) Markup() {

	if t == nil {
		return
	}

}

// Sanatize sanatize config
func (t *BluetoothRPC) Sanatize() {

	if t == nil {
		return
	}

}

// BluetoothObserver configuration of the BT LE observer
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type BluetoothObserver struct {
	// Enable true if BT LE observer is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
}

// Clone return copy
func (t *BluetoothObserver) Clone() *BluetoothObserver {
	c := &BluetoothObserver{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *BluetoothObserver) Markup() {

	if t == nil {
		return
	}

}

// Sanatize sanatize config
func (t *BluetoothObserver) Sanatize() {

	if t == nil {
		return
	}

}
