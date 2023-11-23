package ggchat

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Service interface {
	PushMessage(ctx context.Context, webhook, message string) error
}
type service struct{}

func (s service) PushMessage(ctx context.Context, webhook, message string) error {
	client := resty.New()
	body := map[string]interface{}{
		"text":          message,
		"formattedText": message,
	}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(webhook)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewService() Service {
	return &service{}
}
