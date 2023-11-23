package cron

import (
	"context"
	"fmt"
	"taskbot/pkg/pusher"
	"taskbot/pkg/task"
)

type job struct {
	task   task.Task
	pusher pusher.Pusher
}

func NewJob(task task.Task, pusher pusher.Pusher) *job {
	return &job{task: task, pusher: pusher}
}

func (j job) Schedule() string {
	return j.task.Schedule
}

func (j job) Run() {
	fmt.Println("message: ", j.task.Message)
	err := j.pusher.PushMessage(context.Background(), j.task.WebhookType, j.task.Webhook, j.task.Message)
	if err != nil {
		fmt.Println(err)
	}
}

func (j job) Name() string {
	return j.task.Name
}
