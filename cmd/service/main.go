package main

import (
	"context"
	"fmt"
	"log"

	handler "github.com/cbullard-dev/rose-reaper-notification-service/internal/processor"
	client "github.com/cbullard-dev/rose-reaper-notification-service/internal/websocket"
	"github.com/cbullard-dev/rose-reaper-notification-service/pkg/config"
)

func main() {
	fmt.Println("RoseReaper twitch notification service starting up...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Config loaded:\n LogLevel: %s\n", cfg.LogLevel)

	ctx := context.Background()
	client := &client.Client{
		MsgChannel: make(chan []byte, 10),
	}

	go handler.StartProcessor(ctx, cfg, client.MsgChannel)

	client.Start(ctx, cfg)
}
