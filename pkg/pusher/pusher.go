package pusher

import (
	"context"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/xservice/ggchat"
	"go-reminder-bot/pkg/xservice/msteams"
	"go-reminder-bot/pkg/xservice/slacksvc"
)

type Pusher interface {
	PushMessage(ctx context.Context, whType enum.WebhookType, webhook, message string) error
}

type pusher struct {
	ggChatSvc  ggchat.Service
	msTeamsSvc msteams.Service
	slackSvc   slacksvc.Service
}

func (p pusher) PushMessage(ctx context.Context, whType enum.WebhookType, webhook, message string) error {
	switch whType {
	case enum.WTGoogleChat:
		return p.ggChatSvc.PushMessage(ctx, webhook, message)
	case enum.WTMicrosoftTeams:
		return p.msTeamsSvc.PushMessage(ctx, webhook, message)
	case enum.WTSlack:
		return p.slackSvc.PushMessage(ctx, webhook, message)
	}
	return nil
}

func NewPusher(ggChatSvc ggchat.Service, msTeamsSvc msteams.Service, slackSvc slacksvc.Service) Pusher {
	return &pusher{ggChatSvc: ggChatSvc, msTeamsSvc: msTeamsSvc, slackSvc: slackSvc}
}
