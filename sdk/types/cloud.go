package types

import (
	"github.com/jinzhu/copier"
)

// CloudStatus status of the Cloud component it can be checked whether the device is connected to the cloud.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cloud#status
type CloudStatus struct {
	// Connected true if the device is connected to the Shelly cloud, false otherwise
	Connected bool `json:"connected" yaml:"connected"`
}

// Clone return copy
func (t *CloudStatus) Clone() *CloudStatus {
	c := &CloudStatus{}
	copier.Copy(&c, &t)
	return c
}

// CloudConfig configuration of the Cloud component shows information about the connection to the cloud
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cloud#configuration
type CloudConfig struct {
	// Enable true if cloud connection is enabled, false otherwise
	Enable bool `json:"enable" yaml:"enable"`
	// Server name of the server to which the device is connected
	Server *string `json:"server,omitempty" yaml:"server,omitempty"`
}

// Clone return copy
func (t *CloudConfig) Clone() *CloudConfig {
	c := &CloudConfig{}
	copier.Copy(&c, &t)
	return c
}

// Markup markup config
func (t *CloudConfig) Markup() {

	if t == nil {
		return
	}

	t.Server = nil
}

// Sanatize sanatize config
func (t *CloudConfig) Sanatize() {

	if t == nil {
		return
	}

	t.Server = nil
}
