package types

import (
	"time"

	"github.com/jinzhu/copier"
)

type ClientConfig struct {
	Hostname    string        `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Password    string        `json:"password,omitempty" yaml:"password,omitempty"`
	Username    string        `json:"username,omitempty" yaml:"username,omitempty"`
	SendTimeout time.Duration `json:"sendTimeout,omitempty" yaml:"sendTimeout,omitempty"`
	RetryWait   time.Duration `json:"retryWait,omitempty" yaml:"retryWait,omitempty"`
	SendTrys    int           `json:"sendTrys,omitempty" yaml:"sendTrys,omitempty"`
}

// Clone return copy
func (t *ClientConfig) Clone() *ClientConfig {
	c := &ClientConfig{}
	copier.Copy(&c, &t)
	return c
}
