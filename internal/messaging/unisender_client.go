package messaging

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type UnisenderClient struct {
	ApiKey  string
	BaseURL string
	logger  zerolog.Logger
}

func NewUnisenderClient(logger zerolog.Logger) *UnisenderClient {
	return &UnisenderClient{
		ApiKey:  viper.GetString("UNISENDER_API_KEY"),
		BaseURL: viper.GetString("UNISENDER_BASE_URL"),
		logger:  logger,
	}
}

func (c *UnisenderClient) SendSMS(phone, message, sender string) error {
	sendSmsUrl := fmt.Sprintf("%s/sendSms", c.BaseURL)

	data := url.Values{}
	data.Set("api_key", c.ApiKey)
	data.Set("phone", phone)
	data.Set("sender", sender)
	data.Set("text", message)
	data.Set("format", "json")

	resp, err := http.PostForm(sendSmsUrl, data)

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	c.logger.Debug().Str("Response", resp.Status).Str("Body", bodyString).Msg("Результат отправки по sms")

	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse response: %v", err)
	}

	if status, ok := result["status"].(string); ok && status == "error" {
		return fmt.Errorf("unisender error: %v", result["error"])
	}

	return nil
}

func (c *UnisenderClient) SendEmail(to, subject, body string) error {
	sendEmailUrl := fmt.Sprintf("%s/sendEmail", c.BaseURL)

	data := url.Values{}

	data.Set("api_key", c.ApiKey)
	data.Set("email", to)
	data.Set("sender_name", "Wasted")
	data.Set("sender_email", "wasted@wasted.com")
	data.Set("subject", subject)
	data.Set("body", body)
	data.Set("list_id", "0")
	data.Set("format", "json")

	resp, err := http.PostForm(sendEmailUrl, data)

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	c.logger.Debug().Str("Response", resp.Status).Str("Body", bodyString).Msg("Результат отправки по email")

	if err != nil {
		c.logger.Err(err).Msg("Ошибка при отправке кода по email")
		return fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("parse response: %v", err)
	}

	if status, ok := result["status"].(string); ok && status == "error" {
		return fmt.Errorf("unisender error: %v", result["error"])
	}

	return nil
}
