package messaging

import (
	"fmt"

	"github.com/rs/zerolog"
)

type RealSender struct {
	logger          zerolog.Logger
	unisenderClient *UnisenderClient
}

func NewSender(logger zerolog.Logger) *RealSender {
	return &RealSender{logger: logger, unisenderClient: NewUnisenderClient(logger)}
}

func (s *RealSender) SendCodeViaEmail(email string, code string) {
	body := fmt.Sprintf("Ваш код: %s. Никому не сообщайте.", code)
	s.unisenderClient.SendEmail(email, "Wasted: Авторизационный код", body)
}

func (s *RealSender) SendCodeViaSms(phone string, code string) {
	smsText := fmt.Sprintf("Авторизационный код: %s", code)
	s.unisenderClient.SendSMS(phone, smsText, "Wasted")
}
