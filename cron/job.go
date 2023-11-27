package cron

import (
	"context"
	"fmt"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
)

type job struct {
	reminder reminder.Reminder
	pusher   pusher.Pusher
}

func NewJob(reminder reminder.Reminder, pusher pusher.Pusher) *job {
	return &job{reminder: reminder, pusher: pusher}
}

func (j job) Schedule() string {
	return j.reminder.Schedule
}

func (j job) Run() {
	fmt.Println("message: ", j.reminder.Message)
	err := j.pusher.PushMessage(context.Background(), j.reminder.WebhookType, j.reminder.Webhook, j.reminder.Message)
	if err != nil {
		fmt.Println(err)
	}
}

func (j job) Name() string {
	return j.reminder.Name
}
