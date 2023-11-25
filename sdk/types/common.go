package types

import (
	"github.com/jinzhu/copier"
)

// FirmwareStatus is common for components Sys and Shelly
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Sys#status &
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Shelly#status
type FirmwareStatus struct {
	// Version of the new firmware
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`
	// BuildID Id of the new build
	BuildID *string `json:"build_id,omitempty" yaml:"build_id,omitempty"`
}

// Request generic request
type Request struct {
	ID     *int          `json:"id,omitempty" yaml:"id,omitempty"`
	Method *string       `json:"method,omitempty" yaml:"method,omitempty"`
	Params interface{}   `json:"params,omitempty" yaml:"params,omitempty"`
	Auth   *AuthResponse `json:"auth,omitempty" yaml:"auth,omitempty"`
}

// Clone return copy
func (t *Request) Clone() *Request {
	c := &Request{}
	copier.Copy(&c, &t)
	return c
}

// Response generic response
type Response struct {
	ID    *int    `json:"id,omitempty" yaml:"id,omitempty"`
	Src   *string `json:"src,omitempty" yaml:"src,omitempty"`
	Error *Error  `json:"error,omitempty" yaml:"error,omitempty"`
}

// Clone return copy
func (t *Response) Clone() *Response {
	c := &Response{}
	copier.Copy(&c, &t)
	return c
}
