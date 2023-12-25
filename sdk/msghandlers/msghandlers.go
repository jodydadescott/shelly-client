package msghandlers

import (
	client_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
	"github.com/jodydadescott/shelly-client/sdk/msghandlers/ws"
)

type Config = client_types.Config
type Request = msg_types.Request
type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler

func NewWS(config *Config) MessageHandlerFactory {
	return ws.New(config)
}
