package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"taskbot/pkg/pusher"
	"taskbot/pkg/task"
)

type CronJob interface {
	Start(ctx context.Context)
	Reload(ctx context.Context)
	Stop(ctx context.Context)
}
type cronJob struct {
	*cron.Cron
	storage task.Storage
	pusher  pusher.Pusher
}

func (c cronJob) configure(ctx context.Context) {
	c.Cron = cron.New()
}

func (c cronJob) buildListJobs(ctx context.Context) {
	tasks, err := c.storage.GetActiveTasks(ctx)
	if err != nil {
		return
	}
	for _, t := range tasks {
		job := NewJob(t, c.pusher)
		_, err := c.AddJob(job.Schedule(), job)
		if err != nil {
			fmt.Println("failed to add job: ", job)
			continue
		}
	}
}

func (c cronJob) Start(ctx context.Context) {
	c.configure(ctx)
	c.buildListJobs(ctx)
	c.Cron.Start()
}

func (c cronJob) Reload(ctx context.Context) {
	c.Start(ctx)
}

func (c cronJob) Stop(ctx context.Context) {
	if c.Cron != nil {
		c.Cron.Stop()
	}
}

func NewCron(storage task.Storage, pusher pusher.Pusher) CronJob {
	return &cronJob{
		storage: storage,
		pusher:  pusher}
}
