package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// TelegramNotifier sends alert events to a Telegram chat via the Bot API.
type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
	next   alert.Notifier
}

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramNotifier creates a notifier that posts to Telegram and forwards to next.
func NewTelegramNotifier(token, chatID string, client *http.Client, next alert.Notifier) *TelegramNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &TelegramNotifier{token: token, chatID: chatID, client: client, next: next}
}

func (t *TelegramNotifier) Send(e alert.Event) error {
	text := fmt.Sprintf("*portwatch alert*\nPort: %d/%s\nAddr: %s\nChange: %s",
		e.Port, e.Proto, e.Addr, e.Change)

	payload := telegramPayload{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}

	if t.next != nil {
		return t.next.Send(e)
	}
	return nil
}
