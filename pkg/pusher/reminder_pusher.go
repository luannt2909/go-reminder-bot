package pusher

import (
	"context"
	"fmt"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/reminder"
)

type ReminderPusher interface {
	PushMessage(ctx context.Context, r reminder.Reminder) error
}

type reminderPusher struct {
	pusher Pusher
}

const ggChatMsgFooter = `
================
Sent by: *%s*
Powered by: <https://luciannguyen.blog|reminder-bot.com>
`

func (p reminderPusher) PushMessage(ctx context.Context, r reminder.Reminder) error {
	switch r.WebhookType {
	case enum.WTGoogleChat:
		body := r.Message
		footer := fmt.Sprintf(ggChatMsgFooter, r.CreatedBy)
		message := fmt.Sprint(body, footer)
		return p.pusher.PushMessage(ctx, enum.WTGoogleChat, r.Webhook, message)
	}
	return nil
}

func NewReminderPusher(pusher Pusher) ReminderPusher {
	return &reminderPusher{pusher: pusher}
}
