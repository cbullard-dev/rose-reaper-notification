package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	discord "github.com/cbullard-dev/rose-reaper-notification-service/internal/notifier"
	"github.com/cbullard-dev/rose-reaper-notification-service/pkg/config"
)

const SubStreamOnline string = "stream.online"
const SubStreamOffline string = "stream.offline"

var cfg *config.Config

type SubscriptionNotificationType struct {
	Metadata struct {
		SubscriptionType string `json:"subscription_type"`
	} `json:"metadata"`
}

type StreamEndedNotification struct {
	Payload struct {
		Event struct {
			Id                 string `json:"id"`
			BroadcastUserLogin string `json:"broadcaster_user_login"`
			BroadcastUserName  string `json:"broadcaster_user_name"`
		} `json:"event"`
	} `json:"payload"`
}

type StreamStartedNotification struct {
	Payload struct {
		Event struct {
			Id                 string `json:"id"`
			BroadcastUserName  string `json:"broadcaster_user_name"`
			BroadcastUserLogin string `json:"broadcaster_user_login"`
			StartedAt          string `json:"started_at"`
		} `json:"event"`
	} `json:"payload"`
}

type BaseMessage struct {
	Metadata struct {
		MessageType string `json:"message_type"`
	} `json:"metadata"`
}

type WelcomeMessage struct {
	Payload struct {
		Session struct {
			ID string `json:"id"`
		} `json:"session"`
	} `json:"payload"`
}

type Subscription struct {
	Type      string `json:"type"`
	Version   string `json:"version"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
	} `json:"condition"`
	Transport struct {
		Method    string `json:"method"`
		SessionID string `json:"session_id"`
	} `json:"transport"`
}

type RefreshAuth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func StartProcessor(ctx context.Context, config *config.Config, msgChannel <-chan []byte) {
	cfg = config
	for {
		select {
		case data := <-msgChannel:
			processMessage(data)
		case <-ctx.Done():
			log.Println("Handler shutting down...")
			return
		}
	}
}

func processMessage(data []byte) {
	var base BaseMessage

	err := json.Unmarshal(data, &base)
	if err != nil {
		log.Printf("Error unmarshaling the metadata: %v", err)
		return
	}

	switch base.Metadata.MessageType {
	case "session_welcome":
		var msg WelcomeMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Failed to unmarshal the welcome message: %v", err)
			return
		}

		if err := sendSubscription(msg.Payload.Session.ID, SubStreamOnline); err != nil {
			log.Printf("Failed to send subscription: %v.", err)
			break
		}

		if err := sendSubscription(msg.Payload.Session.ID, SubStreamOffline); err != nil {
			log.Printf("Failed to send subscription: %v.", err)
			break
		}

	case "session_keepalive":
		log.Printf("Session was kept alive!")

	case "notification":
		processNotification(data)
	default:
		log.Printf("Unknown message type: %s", base.Metadata.MessageType)
		log.Printf("Received message: %s\n", string(data))
	}
}

func processNotification(data []byte) {
	var subType SubscriptionNotificationType
	if err := json.Unmarshal(data, &subType); err != nil {
		log.Printf("ERROR: Failed to unmarshal the notification type: %v", err)
	}
	switch subType.Metadata.SubscriptionType {
	case SubStreamOnline:
		var start StreamStartedNotification
		if err := json.Unmarshal(data, &start); err != nil {
			log.Printf("ERROR: Failed to unmarshal the notification type: %v", err)
		}
		log.Printf("The streamer %s has gone online at %s!", start.Payload.Event.BroadcastUserName, start.Payload.Event.StartedAt)
		go func() {
			err := discord.SendMessage(cfg.DiscordWebhookURL, fmt.Sprintf("%s has just gone live.\nHead over to https://twitch.tv/%s and say hello!", start.Payload.Event.BroadcastUserName, start.Payload.Event.BroadcastUserLogin))
			if err != nil {
				log.Printf("Failed to send discord message: %s", err)
			}
		}()

	case SubStreamOffline:
		var ended StreamEndedNotification
		if err := json.Unmarshal(data, &ended); err != nil {
			log.Printf("ERROR: Failed to unmarshal the notification type: %v", err)
		}
		log.Printf("The streamer %s has gone offline!", ended.Payload.Event.BroadcastUserName)
	}
}

func sendSubscription(sessionID string, subscriptionType string) error {
	sub := Subscription{
		Type:    subscriptionType,
		Version: "1",
	}
	sub.Condition.BroadcasterUserID = cfg.TwitchBroadcasterUserID
	sub.Transport.Method = "websocket"
	sub.Transport.SessionID = sessionID

	msg, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(msg)

	client := &http.Client{}

	req, err := http.NewRequest("POST", cfg.TwitchSubscriptionEndpoint, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.TwitchBearerToken))
	req.Header.Set("Client-Id", cfg.TwitchClientID)
	req.Header.Add("Content-Type", "application/json")

	log.Printf("Attempting to subscribe to the %s sub for broadcaster ID %s", subscriptionType, sub.Condition.BroadcasterUserID)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Printf("Response status from subscribe attempt to %s: %s Code: %v\n", subscriptionType, res.Status, res.StatusCode)

	if res.StatusCode == 401 {
		err := TryRefreshAuth()
		if err != nil {
			log.Printf("Error trying to refresh token: %v", err)
			return err
		}
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("response body: %s", string(body))
	}
	return nil
}

func TryRefreshAuth() error {
	refresh := RefreshAuth{
		ClientID:     cfg.TwitchClientID,
		ClientSecret: cfg.TwitchClientSecret,
		Code:         cfg.TwitchClientCode,
		GrantType:    "refresh_token",
		RefreshToken: cfg.TwitchRefreshToken,
	}

	var new RefreshResponse

	log.Println("Starting token refresh attempt...")

	msg, err := json.Marshal(refresh)
	if err != nil {
		return fmt.Errorf("error marshaling json: %v", err)
	}
	bodyReader := bytes.NewReader(msg)

	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/json", bodyReader)
	if err != nil {
		return fmt.Errorf("error trying to get new OAuth2 Token: %v", err)
	}
	log.Printf("Token refresh response status: %d", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error trying to read OAuth response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error in attempting to refresh OAuth token: %v", resp.Status)
	}

	if resp.StatusCode == 200 {
		err := json.Unmarshal(body, &new)
		if err != nil {
			return fmt.Errorf("error trying to unmarshal response body: %v", err)
		}
		cfg.TwitchRefreshToken = new.RefreshToken
		cfg.TwitchBearerToken = new.AccessToken
		log.Println("New token configuration has been set successfully")
	}

	return nil
}
