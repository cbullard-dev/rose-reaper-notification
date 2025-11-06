package config

import (
	"fmt"
	"os"
)

type Config struct {
	TwitchWebSocketURL         string
	TwitchSubscriptionEndpoint string
	TwitchClientID             string
	TwitchClientSecret         string
	TwitchRefreshToken         string
	TwitchBearerToken          string
	TwitchClientCode           string
	TwitchBroadcasterUserID    string
	DiscordWebhookURL          string
	LogLevel                   string
	DatabaseURL                string
}

func Load() (*Config, error) {
	cfg := &Config{
		TwitchSubscriptionEndpoint: "https://api.twitch.tv/helix/eventsub/subscriptions",
		TwitchWebSocketURL:         os.Getenv("TWITCH_WEBSOCKET_URL"),
		TwitchClientID:             os.Getenv("TWITCH_CLIENT_ID"),
		TwitchClientSecret:         os.Getenv("TWITCH_CLIENT_SECRET"),
		TwitchRefreshToken:         os.Getenv("TWITCH_REFRESH_TOKEN"),
		TwitchBearerToken:          os.Getenv("TWITCH_BEARER_TOKEN"),
		TwitchClientCode:           os.Getenv("TWITCH_CLIENT_CODE"),
		TwitchBroadcasterUserID:    os.Getenv("TWITCH_BROADCASTER_USER_ID"),
		DiscordWebhookURL:          os.Getenv("DISCORD_WEBHOOK_URL"),
		LogLevel:                   os.Getenv("LOG_LEVEL"),
		// DatabaseURL:                os.Getenv("DATABASE_URL"),
	}

	if cfg.TwitchWebSocketURL == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_WEBSOCKET_URL")
	}

	if cfg.TwitchClientID == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_CLIENT_ID")
	}

	if cfg.TwitchClientSecret == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_CLIENT_SECRET")
	}

	if cfg.TwitchRefreshToken == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_REFRESH_TOKEN")
	}

	if cfg.TwitchClientCode == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_CLIENT_CODE")
	}

	if cfg.TwitchBroadcasterUserID == "" {
		return nil, fmt.Errorf("missing the required environment variable: TWITCH_BROADCASTER_USER_ID")
	}

	if cfg.DiscordWebhookURL == "" {
		return nil, fmt.Errorf("missing the required environment variable: DISCORD_WEBHOOK_URL")
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}
