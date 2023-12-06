package pusher

import (
	"fmt"
	"go-reminder-bot/pkg/enum"
)

const GGChatMsgFooter = `
================
Sent by: *%s*
Powered by: <https://reminderbot.luciango.com/admin|reminder-bot>
`

const DiscordMsgFooter = `
================
Sent by: **%s**
Powered by: [reminder-bot](https://reminderbot.luciango.com/admin)
`

const defaultMsgFooter = `
================
Sent by: **reminder-bot**
Powered by: https://reminderbot.luciango.com/admin
`

func InjectFooter(whType enum.WebhookType, author, msg string) string {
	if author == "" {
		author = "reminder-bot"
	}
	footer := defaultMsgFooter
	switch whType {
	case enum.WTDiscord:
		footer = fmt.Sprintf(DiscordMsgFooter, author)
	default:
		footer = fmt.Sprintf(GGChatMsgFooter, author)
	}
	return fmt.Sprint(msg, footer)
}
