package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

func SendMessage(discordWebhookURL string, content string) error {
	msg := DiscordMessage{Content: content}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}
	res, err := http.Post(discordWebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send discord message: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("discord returned status code %d", res.StatusCode)
	}

	log.Printf("Sent message to discord\n")
	return nil
}
