package types

import (
	"fmt"
)

// Error Shelly Error
type Error struct {
	Code    int    `json:"code,omitempty" yaml:"code,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

func (t *Error) Error() string {

	if t == nil {
		panic("error is nil")
	}

	m := errorCodeMap[t.Code]

	if m != "" {
		return fmt.Sprintf("status %d: err %s; %s", t.Code, t.Message, m)
	}

	return fmt.Sprintf("status %d: err %s", t.Code, t.Message)
}

type ErrorCode int

var getErrorCodeMap = func() map[int]string {
	return map[int]string{
		ErrorCodeInvalidArgument:    "parameters sent in the request do not match the ones specified by the method in the request",
		ErrorCodeDeadLineExceeded:   "request timeout",
		ErrorCodeNotFound:           "instance with specified ID not found",
		ErrorCodeResourceExhausted:  "resource has reached its limit. For example, when you try to create 21 schedule jobs on one Shelly device (the limit is 20)",
		ErrorCodeFailedPrecondition: "precondition for a requested action is not satisfied. For example, when you try to turn a switch on in a situation of overpower condition, or when a reboot has been scheduled and the device is shutting down",
		ErrorCodeUnAvailable:        "service is unavailable. The service can be internal - a sensor could be unreachable, or external. External services are - timezone information, firmware update or HTTP requests in Scripts.",
		ErrorCodeNotImplemented:     "method is not implemented on this device or caller is not authorized",
	}
}

const (
	ErrorCodeInvalidArgument    = -103
	ErrorCodeDeadLineExceeded   = -104
	ErrorCodeNotFound           = -105
	ErrorCodeResourceExhausted  = -108
	ErrorCodeFailedPrecondition = -109
	ErrorCodeUnAvailable        = -114
	ErrorCodeNotImplemented     = 404
)

var errorCodeMap = getErrorCodeMap()
