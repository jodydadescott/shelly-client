package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	gorilla "github.com/gorilla/websocket"
	logger "github.com/jodydadescott/jody-go-logger"
	"go.uber.org/zap"

	client_types "github.com/jodydadescott/shelly-client/sdk/client/types"
	msg_types "github.com/jodydadescott/shelly-client/sdk/msghandlers/types"
)

type Config = client_types.Config
type Response = msg_types.Response

type MessageHandlerFactory = msg_types.MessageHandlerFactory
type MessageHandler = msg_types.MessageHandler
type Request = msg_types.Request
type AuthResponse = msg_types.AuthResponse
type AuthRequest = msg_types.AuthRequest

type Client struct {
	config            *Config
	authResponseMutex sync.RWMutex
	handleMutex       sync.RWMutex
	handleMap         map[int]*Handle
	egressMessages    chan []byte
	uniqID            int
	wg                sync.WaitGroup
	cancel            context.CancelFunc
	authResponse      *AuthResponse
}

func New(config *Config) MessageHandlerFactory {
	zap.L().Debug("New")

	config = config.Clone()

	if config.Hostname == "" {
		panic("hostname is required")
	}

	if config.Username == "" {
		config.Username = defaultShellyUser
		zap.L().Debug(fmt.Sprintf("username is %s (default)", config.Username))
	} else {
		zap.L().Debug(fmt.Sprintf("username is %s (config)", config.Username))
	}

	if config.SendTimeout <= 0 {
		config.SendTimeout = defaultSendTimeout
		zap.L().Debug(fmt.Sprintf("sendTimeout is %s (default)", config.SendTimeout.String()))
	} else {
		zap.L().Debug(fmt.Sprintf("sendTimeout is %s (config)", config.SendTimeout.String()))
	}

	if config.RetryWait <= 0 {
		config.RetryWait = defaultRetryWait
		zap.L().Debug(fmt.Sprintf("retryWait is %s (default)", config.RetryWait.String()))
	} else {
		zap.L().Debug(fmt.Sprintf("retryWait is %s (config)", config.RetryWait.String()))
	}

	if config.SendTrys <= 0 {
		config.SendTrys = defaultSendTrys
		zap.L().Debug(fmt.Sprintf("sendTrys is %d (default)", config.SendTrys))
	} else {
		zap.L().Debug(fmt.Sprintf("sendTrys is %d (config)", config.SendTrys))
	}

	if config.Password == "" {
		zap.L().Debug("password is NOT set")
	} else {
		zap.L().Debug("password is set")
	}

	t := &Client{
		config:         config,
		handleMap:      make(map[int]*Handle),
		egressMessages: make(chan []byte, 50),
	}

	t.run()
	return t
}

func (t *Client) IsAuthEnabled() bool {
	return t.getAuthResponse() != nil
}

func (t *Client) getAuthResponse() *AuthResponse {
	t.authResponseMutex.RLock()
	defer t.authResponseMutex.RUnlock()
	return t.authResponse
}

func (t *Client) setAuthResponse(authResponse *AuthResponse) {
	t.authResponseMutex.Lock()
	defer t.authResponseMutex.Unlock()
	t.authResponse = authResponse
}

func (t *Client) Close() {
	zap.L().Debug("(*Client) Close()")

	for _, handle := range t.handleMap {
		handle.close()
	}

	t.cancel()
	t.wg.Wait()
}

func (t *Client) run() {

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	routeMessage := func(b []byte) {

		msg := &Response{}
		err := json.Unmarshal(b, msg)
		if err != nil {
			zap.L().Error(fmt.Sprintf("routeMessage error %v", err))
			return
		}

		zap.L().Debug("getting handle mutex")
		t.handleMutex.RLock()
		defer t.handleMutex.RUnlock()
		defer zap.L().Debug("releasing handle mutex")

		handle := t.handleMap[*msg.ID]
		if handle == nil {
			zap.L().Error(fmt.Sprintf("handle lookup ID %d failure", msg.ID))
			return
		}

		handle.receive <- &responseWrapper{
			response: msg,
			rawBytes: b,
		}
	}

	handleEgress := func(conn *gorilla.Conn, errChan *errChan) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					zap.L().Debug("handleEgress closed by context")
					conn.WriteMessage(gorilla.CloseMessage, gorilla.FormatCloseMessage(gorilla.CloseNormalClosure, ""))
					return

				case b := <-t.egressMessages:

					if logger.Wire {
						zap.L().Debug(fmt.Sprintf("TX->%s", string(b)))
					}

					err := conn.WriteMessage(gorilla.BinaryMessage, b)
					if err != nil {
						zap.L().Debug("handleEgress closed by error")
						errChan.putError(err)
						return
					}

				}
			}

		}()
	}

	handleIngress := func(conn *gorilla.Conn, errChan *errChan) {
		go func() {
			for {
				_, b, err := conn.ReadMessage()

				if logger.Wire {
					zap.L().Debug(fmt.Sprintf("RX->%s", string(b)))
				}

				if err != nil {
					zap.L().Debug("handleIngress closed by error")
					errChan.putError(err)
					return
				}

				routeMessage(b)
			}
		}()
	}

	handle := func(conn *gorilla.Conn) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		errChan := newErrChan()

		handleIngress(conn, errChan)
		handleEgress(conn, errChan)

		err := errChan.getError(ctx)
		errChan.close()
		return err
	}

	connect := func() error {
		theURL := url.URL{Scheme: WsScheme, Host: t.config.Hostname, Path: wsPath}
		conn, _, err := gorilla.DefaultDialer.Dial(theURL.String(), nil)

		if err != nil {
			return err
		}

		zap.L().Debug("Connected")

		return handle(conn)
	}

	t.wg.Add(1)

	go func() {

		defer t.wg.Done()
		defer close(t.egressMessages)

		for {

			zap.L().Debug(fmt.Sprintf("Connecting to %s", t.config.Hostname))

			err := connect()

			if err == nil {
				return
			}

			zap.L().Debug(fmt.Sprintf("Connect error %v; will try again in %v", err, defaultRetryWait.String()))

			select {

			case <-ctx.Done():
				zap.L().Debug("Connect cancelled")
				return

			case <-time.After(defaultRetryWait):
				continue
			}

		}

	}()

}

func (t *Client) NewHandle(name string) MessageHandler {

	zap.L().Debug(fmt.Sprintf("(*Client) NewHandle(%s)", name))

	zap.L().Debug("getting handle mutex")
	t.handleMutex.Lock()
	defer t.handleMutex.Unlock()
	defer zap.L().Debug("releasing handle mutex")

	t.uniqID = t.uniqID + 1

	handle := &Handle{
		client:  t,
		id:      t.uniqID,
		receive: make(chan *responseWrapper, 30),
		done:    make(chan struct{}),
	}

	t.handleMap[handle.id] = handle

	zap.L().Debug(fmt.Sprintf("(*Client) NewHandle(%s) returned", name))

	return handle
}

func (t *Handle) close() {
	zap.L().Debug("(*Handle) close()")
	close(t.done)
	close(t.receive)
}

type responseWrapper struct {
	response *Response
	rawBytes []byte
}

type Handle struct {
	client  *Client
	id      int
	receive chan *responseWrapper
	done    chan struct{}
}

func (t *Handle) Send(ctx context.Context, request *Request) ([]byte, error) {

	zap.L().Debug("(*Handle) Send(ctx, *Request)")

	request = request.Clone()
	request.ID = &t.id

	request.Auth = t.client.getAuthResponse()

	if request.Auth != nil {
		zap.L().Debug("Using previous auth")
	} else {
		zap.L().Debug("Auth is not set")
	}

	requestBytes, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	waitOnResponse := func() (*responseWrapper, error) {

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		counter := 0

		for {

			t.client.egressMessages <- requestBytes

			select {

			case response := <-t.receive:
				return response, nil

			case <-t.done:
				return nil, fmt.Errorf("channel closed shutdown")

			case <-ctx.Done():
				return nil, fmt.Errorf("channel closed by caller")

			case <-time.After(t.client.config.SendTimeout):
				if counter >= t.client.config.SendTrys {
					zap.L().Debug(fmt.Sprintf("try %d of %d; giving up", counter, t.client.config.SendTrys))
					return nil, fmt.Errorf("timeout waiting for response")
				}
				zap.L().Debug(fmt.Sprintf("try %d of %d; will try again", counter, t.client.config.SendTrys))
				counter++

			}

		}

	}

	response, err := waitOnResponse()

	if err != nil {
		return nil, err
	}

	if response.response.Error != nil {

		if response.response.Error.Code == 401 {

			zap.L().Debug("server responded with auth required")

			if t.client.config.Username == "" {
				return nil, fmt.Errorf("username is required")
			}

			if t.client.config.Password == "" {
				return nil, fmt.Errorf("password is required")
			}

			authRequest := &AuthRequest{}
			err = json.Unmarshal([]byte(response.response.Error.Message), authRequest)
			if err != nil {
				return nil, err
			}

			authRequest.Username = t.client.config.Username
			authRequest.Password = t.client.config.Password
			authResponse, err := authRequest.ToAuthResponse()

			if err != nil {
				return nil, err
			}

			t.client.setAuthResponse(authResponse)

			request.Auth = authResponse

			requestBytes, err := json.Marshal(request)
			if err != nil {
				return nil, err
			}

			t.client.egressMessages <- requestBytes

			response, err := waitOnResponse()

			if err != nil {
				return nil, err
			}

			return response.rawBytes, nil

		}

		return nil, response.response.Error
	}

	return response.rawBytes, nil
}

type errChan struct {
	closed bool
	mutex  sync.Mutex
	errs   chan error
}

func newErrChan() *errChan {
	return &errChan{
		errs: make(chan error, 2),
	}
}

func (t *errChan) putError(err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.closed {
		return
	}

	t.closed = true
	t.errs <- err
}

func (t *errChan) getError(ctx context.Context) error {

	select {

	case <-ctx.Done():
		return nil

	case err := <-t.errs:
		return err

	}
}

func (t *errChan) close() {

	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.closed {
		return
	}

	close(t.errs)
	t.closed = true
}
