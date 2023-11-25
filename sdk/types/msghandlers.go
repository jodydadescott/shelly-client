package types

import (
	"context"
)

type MessageHandlerFactory interface {
	NewHandle() MessageHandler
	Close()
}

type MessageHandler interface {
	Send(ctx context.Context, request *Request) ([]byte, error)
	Close()
	IsAuthEnabled() bool
}
