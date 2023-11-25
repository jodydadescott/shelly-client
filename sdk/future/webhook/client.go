package webhook

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// )

// // GenericResponse internal use only
// type GenericResponse struct {
// 	Response
// }

// // ListSupportedResponse internal use only
// type ListSupportedResponse struct {
// 	Response
// 	Result *ListSupportedResult `json:"result,omitempty"`
// }

// // ListSupportedResult internal use only
// // https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Webhook#webhooklistsupported
// type ListSupportedResult struct {
// 	HookTypes []string `json:"hook_types"`
// }

// type GetConfigResponse struct {
// 	Response
// 	Result *WebhookConfig `json:"result,omitempty"`
// }

// // CreateResult internal use only
// // https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Webhook#webhookcreate

// type CreateResult struct {
// 	ID       *int `json:"id"`
// 	Revision *int `json:"rev"`
// }

// type CreateResponse struct {
// 	Response
// 	Result *CreateResult `json:"result,omitempty"`
// }

// func New(messageHandlerFactory MessageHandlerFactory) *Client {
// 	return &Client{
// 		MessageHandlerFactory: messageHandlerFactory,
// 	}
// }

// type Client struct {
// 	MessageHandlerFactory
// 	_messageHandler MessageHandler
// }

// func (t *Client) getMessageHandler() MessageHandler {
// 	if t._messageHandler != nil {
// 		return t._messageHandler
// 	}

// 	t._messageHandler = t.NewHandle()
// 	return t._messageHandler
// }

// // ListSupported lists all supported events that can be used to trigger a Webhook
// func (t *Client) ListSupported(ctx context.Context) ([]string, error) {

// 	method := Component + ".ListSupported"

// 	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
// 		Method: &method,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &ListSupportedResponse{}
// 	err = json.Unmarshal(respBytes, response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if response.Error != nil {
// 		return nil, response.Error
// 	}

// 	if response.Result == nil {
// 		return nil, fmt.Errorf("Result is missing from response")
// 	}

// 	return response.Result.HookTypes, nil
// }

// // List lists all existing Webhooks for this device.
// func (t *Client) GetConfig(ctx context.Context) (*WebhookConfig, error) {

// 	method := Component + ".List"

// 	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
// 		Method: &method,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &GetConfigResponse{}
// 	err = json.Unmarshal(respBytes, response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if response.Error != nil {
// 		return nil, response.Error
// 	}

// 	if response.Result == nil {
// 		return nil, fmt.Errorf("Result is missing from response")
// 	}

// 	return response.Result, nil
// }

// func (t *Client) createOrUpdate(ctx context.Context, webhook *Webhook) (*Webhook, error) {

// 	webhook.Clone()
// 	webhook.Sanatize()

// 	method := ".Create"
// 	if webhook.ID != nil {
// 		method = ".Update"
// 	}

// 	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
// 		Method: &method,
// 		Params: webhook,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &CreateResponse{}
// 	err = json.Unmarshal(respBytes, response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if response.Error != nil {
// 		return nil, response.Error
// 	}

// 	if response.Result == nil {
// 		return nil, fmt.Errorf("Result is missing from response")
// 	}

// 	config, err := t.GetConfig(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range config.Webhooks {

// 		if *v.ID == *response.Result.ID {
// 			return v, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("Webhook was successfully created but not found by ID")
// }

// // Create creates a Webhook instance
// func (t *Client) Create(ctx context.Context, webhook *Webhook) (*Webhook, error) {
// 	return t.createOrUpdate(ctx, webhook)
// }

// // Update creates a Webhook instance
// func (t *Client) Update(ctx context.Context, webhook *Webhook) (*Webhook, error) {

// 	if webhook.ID == nil {
// 		return nil, fmt.Errorf("ID is required for update")
// 	}

// 	return t.createOrUpdate(ctx, webhook)
// }

// // Delete deletes an existing Webhook instance
// func (t *Client) Delete(ctx context.Context, id int) error {

// 	method := Component + ".Delete"

// 	webhook := &Webhook{
// 		ID: &id,
// 	}

// 	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
// 		Method: &method,
// 		Params: webhook,
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	response := &GenericResponse{}
// 	err = json.Unmarshal(respBytes, response)
// 	if err != nil {
// 		return err
// 	}

// 	if response.Error != nil {
// 		return response.Error
// 	}

// 	return nil
// }

// // Delete deletes all existing Webhooks
// func (t *Client) DeleteAll(ctx context.Context) error {

// 	method := Component + ".DeleteAll"

// 	respBytes, err := t.getMessageHandler().Send(ctx, &Request{
// 		Method: &method,
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	response := &GenericResponse{}
// 	err = json.Unmarshal(respBytes, response)
// 	if err != nil {
// 		return err
// 	}

// 	if response.Error != nil {
// 		return response.Error
// 	}

// 	return nil
// }

// func (t *Client) SetConfig(ctx context.Context, config *Config) error {

// 	config = config.Clone()
// 	config.Sanatize()

// 	err := t.setConfig(ctx, config)
// 	if err != nil {
// 		return err
// 	}

// 	currentConfig, err := t.GetConfig(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	for _, v := range currentConfig.Webhooks {
// 		if config.GetWebhook(*v.ID) == nil {
// 			err = t.Delete(ctx, *v.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func (t *Client) UpdateConfig(ctx context.Context, config *Config) error {

// 	config = config.Clone()
// 	config.Sanatize()

// 	return t.setConfig(ctx, config)
// }

// func (t *Client) setConfig(ctx context.Context, config *Config) error {

// 	for _, v := range config.Webhooks {
// 		if v.ID == nil {
// 			return fmt.Errorf("All webhooks must have an ID")
// 		}
// 		_, err := t.createOrUpdate(ctx, v)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (t *Client) Close() {
// 	if t._messageHandler != nil {
// 		t._messageHandler.Close()
// 	}
// }
