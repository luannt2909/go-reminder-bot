package server

import (
	"encoding/json"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/util"
)

type GetListReminderRequest struct {
	Filter map[string]interface{} `form:"filter"`
	Range  string                 `form:"range"`
	Sort   string                 `form:"sort"`
}

func (r GetListReminderRequest) toGetListParams() reminder.GetListParams {
	p := reminder.GetListParams{
		Filter:   r.Filter,
		Limit:    10,
		Offset:   0,
		SortBy:   "id",
		SortType: "ASC",
	}
	if r.Range != "" {
		var queryRange []int
		_ = json.Unmarshal([]byte(r.Range), &queryRange)
		if len(queryRange) == 2 {
			p.Offset, p.Limit = queryRange[0], queryRange[1]
		}
	}
	if r.Sort != "" {
		var querySort []string
		_ = json.Unmarshal([]byte(r.Sort), &querySort)
		if len(querySort) == 2 {
			p.SortBy, p.SortType = querySort[0], querySort[1]
		}
	}
	return p
}

type Reminder struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	Type        string `json:"type"`
	Schedule    string `json:"schedule"`
	NextTime    string `json:"next_time"`
	Message     string `json:"message"`
	Webhook     string `json:"webhook"`
	WebhookType string `json:"webhook_type"`
}

func transformRemindersFromRemindersDB(reminders []reminder.Reminder) []Reminder {
	result := make([]Reminder, 0, len(reminders))
	for _, reminder := range reminders {
		result = append(result, transformReminderFromReminderDB(reminder))
	}
	return result
}

func transformReminderFromReminderDB(t reminder.Reminder) Reminder {
	var nextTime string
	if nTime, err := util.Parse(t.Schedule); err == nil {
		nextTime = nTime.Format("15:04:05 02-01-2006")
	} else {
		nextTime = "invalid"
	}
	return Reminder{
		ID:          int64(t.Model.ID),
		Name:        t.Name,
		IsActive:    t.IsActive,
		Schedule:    t.Schedule,
		NextTime:    nextTime,
		Message:     t.Message,
		Webhook:     t.Webhook,
		WebhookType: t.WebhookType.String(),
	}
}
