package jsonrpc

import "gotrading/socket"

type JsonRPC2Client struct {
	*socket.WebSocketClient
}

func NewJsonRPC2Client(url string) (*JsonRPC2Client, error) {
	websocketClient, err := socket.NewWebSocketClient(url, nil)
	if err != nil {
		return nil, err
	}

	return &JsonRPC2Client{
		WebSocketClient: websocketClient,
	}, nil
}

func (c *JsonRPC2Client) Send(jsonRPC *JsonRPC2) error {
	if err := c.WebSocketClient.Send(jsonRPC); err != nil {
		return err
	}

	return nil
}

func (c *JsonRPC2Client) Recieve(jsonRPC *JsonRPC2) error {
	if err := c.WebSocketClient.Recieve(jsonRPC); err != nil {
		return err
	}

	return nil
}
