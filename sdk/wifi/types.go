package wifi

import (
	"github.com/jodydadescott/shelly-client/sdk/types"
)

type Request = types.Request
type Response = types.Response
type Error = types.Error
type MessageHandlerFactory = types.MessageHandlerFactory
type MessageHandler = types.MessageHandler

type Config = types.WifiConfig
type Status = types.WifiStatus
type WifiScanResults = types.WifiScanResults
type APClients = types.WifiAPClients
type WifiNet = types.WifiNet
type APClient = types.WifiAPClient
type APConfig = types.WifiAPConfig
type StaConfig = types.WifiSTAConfig
type RoamConfig = types.WifiRoamConfig
type RangeExtenderConfig = types.WifiRangeExtenderConfig

// Result internal use only
type Result struct {
	RestartRequired *bool  `json:"restart_required,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

// Params internal use only
type Params struct {
	Config *Config `json:"config,omitempty"`
}

// GetConfigResponse internal use only
type GetConfigResponse struct {
	Response
	Result *Config `json:"result,omitempty"`
}

// SetConfigResponse internal use only
type SetConfigResponse struct {
	Response
	Result *Result `json:"result,omitempty"`
}

// GetStatusResponse internal use only
type GetStatusResponse struct {
	Response
	Result *Status `json:"result,omitempty"`
}

// ScanResponse internal use only
type ScanResponse struct {
	Response
	Result *WifiScanResults `json:"result,omitempty"`
}

// ListAPClientsResponse internal use only
type ListAPClientsResponse struct {
	Response
	Result *APClients `json:"result,omitempty"`
}
