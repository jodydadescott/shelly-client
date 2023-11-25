package msghandlers

import (
	"github.com/jodydadescott/shelly-client/sdk/msghandlers/ws"
	"github.com/jodydadescott/shelly-client/sdk/types"
)

type Config = types.ClientConfig
type Request = types.Request
type MessageHandlerFactory = types.MessageHandlerFactory
type MessageHandler = types.MessageHandler

func NewWS(config *Config) (MessageHandlerFactory, error) {
	return ws.New(config)
}
