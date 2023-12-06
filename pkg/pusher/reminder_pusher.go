package pusher

import (
	"context"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/reminder"
)

type ReminderPusher interface {
	PushMessage(ctx context.Context, r reminder.Reminder) error
}

type reminderPusher struct {
	pusher Pusher
}

func (p reminderPusher) PushMessage(ctx context.Context, r reminder.Reminder) error {
	message := InjectFooter(r.WebhookType, r.CreatedBy, r.Message)
	switch r.WebhookType {
	case enum.WTGoogleChat:
		return p.pusher.PushMessage(ctx, enum.WTGoogleChat, r.Webhook, message)
	}
	return nil
}

func NewReminderPusher(pusher Pusher) ReminderPusher {
	return &reminderPusher{pusher: pusher}
}
