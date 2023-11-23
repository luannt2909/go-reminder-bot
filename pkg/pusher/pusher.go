package pusher

import (
	"context"
	"taskbot/pkg/enum"
	"taskbot/pkg/xservice/ggchat"
)

type Pusher interface {
	PushMessage(ctx context.Context, whType enum.WebhookType, webhook, message string) error
}

type pusher struct {
	ggChatSvc ggchat.Service
}

func (p pusher) PushMessage(ctx context.Context, whType enum.WebhookType, webhook, message string) error {
	switch whType {
	case enum.WTGoogleChat:
		return p.ggChatSvc.PushMessage(ctx, webhook, message)
	}
	return nil
}

func NewPusher(ggChatSvc ggchat.Service) Pusher {
	return &pusher{ggChatSvc: ggChatSvc}
}
