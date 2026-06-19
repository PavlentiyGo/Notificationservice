package sender

import (
	"fmt"
	"github.com/PavlentiyGo/notification-service/services/notification/config"
	"net/http"
)

type Sender struct {
	botToken string
}

func NewSender(cfg config.Config) *Sender {
	return &Sender{
		botToken: cfg.BotToken,
	}
}

func (s *Sender) SendTgMessage(
	text string,
	userId int64,
) error {
	urlString := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		s.botToken,
		userId,
		text,
	)
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return fmt.Errorf("failed to crete new request: %w", err)
	}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return nil
}
