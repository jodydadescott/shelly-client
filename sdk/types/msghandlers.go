package types

import (
	"context"
)

type MessageHandlerFactory interface {
	NewHandle(string) MessageHandler
	IsAuthEnabled() bool
	Close()
}

type MessageHandler interface {
	Send(ctx context.Context, request *Request) ([]byte, error)
}
