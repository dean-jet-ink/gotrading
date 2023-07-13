package socket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn      *websocket.Conn
	writeLock sync.Mutex
}

func NewWebSocketClient(url string, header http.Header) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, fmt.Errorf("NewWebSocketClient: %s", err.Error())
	}

	return &WebSocketClient{
		conn: conn,
	}, nil
}

func (c *WebSocketClient) Send(message interface{}) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	err := c.conn.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("Send: %s", err.Error())
	}

	return nil
}

func (c *WebSocketClient) Recieve(message interface{}) error {
	err := c.conn.ReadJSON(message)
	if err != nil {
		return fmt.Errorf("Recieve: %s", err.Error())
	}

	return nil
}

func (c *WebSocketClient) Close() error {
	return c.conn.Close()
}
