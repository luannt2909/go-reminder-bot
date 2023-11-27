package cron

import (
	"context"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/robfig/cron/v3"
	"go-reminder-bot/pkg/consts"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
)

type CronJob interface {
	Start(ctx context.Context)
	Reload(ctx context.Context)
	Stop(ctx context.Context)
}
type cronJob struct {
	*cron.Cron
	storage  reminder.Storage
	pusher   pusher.Pusher
	eventBus EventBus.Bus
}

func (c *cronJob) configure(ctx context.Context) {
	c.Cron = cron.New()
	err := c.eventBus.Subscribe(consts.ReminderEventBusTopic, c.Reload)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *cronJob) buildListJobs(ctx context.Context) (jobs []*job, err error) {
	reminders, err := c.storage.GetActiveReminder(ctx)
	if err != nil {
		return
	}
	if len(reminders) == 0 {
		fmt.Println("skip build jobs, list reminder is empty")
		return
	}
	fmt.Println("reminders: ", reminders)
	jobs = make([]*job, 0, len(reminders))
	for _, t := range reminders {
		job := NewJob(t, c.pusher)
		jobs = append(jobs, job)
	}
	return
}

func (c *cronJob) Start(ctx context.Context) {
	c.configure(ctx)
	jobs, err := c.buildListJobs(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.runJobs(ctx, jobs)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *cronJob) Reload(ctx context.Context) {
	fmt.Println("reload cron job")
	if c.Cron != nil {
		c.Cron.Stop()
	}
	c.Cron = cron.New()
	jobs, err := c.buildListJobs(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.runJobs(ctx, jobs)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *cronJob) runJobs(ctx context.Context, jobs []*job) error {
	fmt.Println("run jobs")
	if c.Cron == nil {
		c.Cron = cron.New()
	}
	for _, job := range jobs {
		_, err := c.Cron.AddJob(job.Schedule(), job)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	c.Cron.Start()
	return nil
}

func (c *cronJob) Stop(ctx context.Context) {
	if c.Cron != nil {
		c.Cron.Stop()
	}
	err := c.eventBus.Unsubscribe(consts.ReminderEventBusTopic, c.Reload)
	if err != nil {
		fmt.Println(err)
	}
}

func NewCron(storage reminder.Storage, pusher pusher.Pusher, eventBus EventBus.Bus) CronJob {
	return &cronJob{
		storage:  storage,
		pusher:   pusher,
		eventBus: eventBus,
	}
}
