package client

import (
	"context"
	"log"
	"time"

	"github.com/cbullard-dev/rose-reaper-notification-service/pkg/config"
	"github.com/coder/websocket"
)

type Client struct {
	conn       *websocket.Conn
	MsgChannel chan []byte
}

func (c *Client) ReadLoop(ctx context.Context) error {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return err
		}

		select {
		case c.MsgChannel <- data:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Client) Start(ctx context.Context, cfg *config.Config) {
	for {
		log.Println("Attempting to connect to websocket server...")

		// Connect
		conn, _, err := websocket.Dial(ctx, cfg.TwitchWebSocketURL, nil)
		if err != nil {
			log.Printf("Failed to connect: %v. Retrying in 1s...", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.conn = conn
		log.Println("Websocket server connected!")

		// sessionCh := make(chan string, 1)
		err = c.ReadLoop(ctx)
		if err != nil {
			log.Printf("Read loop exited: %v", err)
		}

		status := websocket.CloseStatus(err)
		if status == websocket.StatusNormalClosure {
			log.Println("Server closed the connection normally. Reconnecting in 1s...")
		} else {
			log.Printf("Server connection was closed unexpectedly with status: %d: %v. Reconnecting in 1s...\n", status, err)
		}

		conn.Close(websocket.StatusInternalError, "closing connection")
		time.Sleep(1 * time.Second)
	}
}
