package reminder

import (
	"go-reminder-bot/pkg/enum"
	"gorm.io/gorm"
	"time"
)

type Reminder struct {
	gorm.Model
	Name        string           `json:"name" gorm:"name"`
	Schedule    string           `json:"schedule" gorm:"schedule"`
	Message     string           `json:"message" gorm:"message"`
	Webhook     string           `json:"webhook" gorm:"webhook"`
	WebhookType enum.WebhookType `json:"webhook_type" gorm:"webhook_type"`
	IsActive    bool             `json:"is_active" gorm:"is_active"`
	CreatedBy   string           `json:"created_by" gorm:"created_by;size:256;index"`
}

var DefaultReminder = Reminder{
	Name:        "Daily reminder report",
	Schedule:    "* * * * *",
	Message:     "report reminder please",
	Webhook:     "this is webhook URL",
	WebhookType: enum.WTGoogleChat,
	IsActive:    false,
	Model: gorm.Model{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
