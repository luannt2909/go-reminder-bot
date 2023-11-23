package enum

type WebhookType int

const (
	WTGoogleChat WebhookType = iota + 1
	WTDiscord
	WTSlack
)

var webhookTypeMap = map[WebhookType]string{
	WTGoogleChat: "google_chat",
	WTDiscord:    "discord",
	WTSlack:      "slack",
}

var webhookTypeStrMap = map[string]WebhookType{
	"google_chat": WTGoogleChat,
	"discord":     WTDiscord,
	"slack":       WTSlack,
}

func (t WebhookType) String() string {
	return webhookTypeMap[t]
}

func ParseWebhookType(str string) WebhookType {
	return webhookTypeStrMap[str]
}
