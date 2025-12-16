package messaging

import (
	"fmt"
)

type FakeSender struct{}

func NewFakeSender() *FakeSender {
	return &FakeSender{}
}

func (s *FakeSender) SendCodeViaEmail(email string, code string) {
	fmt.Printf("Ваш код: %s. Никому не сообщайте.", code)
}

func (s *FakeSender) SendCodeViaSms(phone string, code string) {
	fmt.Printf("Авторизационный код: %s", code)
}
