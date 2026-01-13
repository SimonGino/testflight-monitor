package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Notifier interface for different notification channels
type Notifier interface {
	Send(title, message string) error
}

// TelegramNotifier sends notifications via Telegram bot
type TelegramNotifier struct {
	BotToken string
	ChatID   string
	client   *http.Client
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(botToken, chatID string, proxyURL string) *TelegramNotifier {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

	return &TelegramNotifier{
		BotToken: botToken,
		ChatID:   chatID,
		client:   client,
	}
}

// Send sends a message via Telegram
func (t *TelegramNotifier) Send(title, message string) error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("telegram not configured")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := map[string]interface{}{
		"chat_id":    t.ChatID,
		"text":       fmt.Sprintf("*%s*\n\n%s", title, message),
		"parse_mode": "Markdown",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned %d", resp.StatusCode)
	}

	return nil
}
