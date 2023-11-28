package cron

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/robfig/cron/v3"
	"go-reminder-bot/pkg/consts"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/user"
	"log"
	"sync"
)

type UserReminderJob interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}
type userReminderJob struct {
	cronJob         *cron.Cron
	userStorage     user.Storage
	reminderStorage reminder.Storage
	pusher          pusher.Pusher
	reminderMap     map[uint]cron.EntryID
	subscriber      EventBus.BusSubscriber
	syncLocker      sync.Locker
}

func (u *userReminderJob) Stop(ctx context.Context) {
	if u.cronJob != nil {
		u.cronJob.Stop()
	}
	if u.subscriber != nil {
		err := u.subscriber.Unsubscribe(consts.ReminderEventBusTopic, u.processReminderEvent)
		if err != nil {
			log.Println("failed to unsubscribe reminder event: ", err)
		}
	}
}

func NewUserReminderCronJob(
	userStorage user.Storage,
	reminderStorage reminder.Storage,
	pusher pusher.Pusher,
	subscriber EventBus.BusSubscriber) UserReminderJob {
	return &userReminderJob{
		userStorage:     userStorage,
		reminderStorage: reminderStorage,
		pusher:          pusher,
		subscriber:      subscriber,
	}
}
func (u *userReminderJob) Configure(ctx context.Context) error {
	u.cronJob = cron.New()
	u.reminderMap = make(map[uint]cron.EntryID)
	u.syncLocker = &sync.Mutex{}
	err := u.subscribeReminderUpdateEvent()
	return err
}
func (u *userReminderJob) Start(ctx context.Context) {
	err := u.Configure(ctx)
	if err != nil {
		return
	}
	users, err := u.userStorage.GetActiveUsers(ctx)
	if err != nil {
		return
	}
	emails := make([]string, 0, len(users))
	for _, us := range users {
		emails = append(emails, us.Email)
	}
	reminders, err := u.reminderStorage.GetActiveRemindersByUsers(ctx, emails)
	if err != nil {
		return
	}
	for _, r := range reminders {
		u.addReminderJob(ctx, r)
	}
	u.cronJob.Start()
}

func (u *userReminderJob) subscribeReminderUpdateEvent() error {
	err := u.subscriber.Subscribe(consts.ReminderEventBusTopic, u.processReminderEvent)
	return err
}

func (u *userReminderJob) processReminderEvent(reminderID uint, event enum.ReminderEvent) {
	log.Printf("process reminder event, reminderID: %d, event: %d", reminderID, event)
	ctx := context.Background()
	switch event {
	case enum.RECreate:
		u.addReminderJobByID(ctx, reminderID)
	case enum.REUpdate:
		u.updateReminderJobByID(ctx, reminderID)
	case enum.REDelete:
		u.deleteReminderJobByID(ctx, reminderID)
	}
}
func (u *userReminderJob) deleteReminderJobByID(ctx context.Context, reminderID uint) {
	if entryID, ok := u.reminderMap[reminderID]; ok {
		u.cronJob.Remove(entryID)
		u.removeReminderMap(reminderID)
	}
}

func (u *userReminderJob) updateReminderJobByID(ctx context.Context, reminderID uint) {
	r, err := u.reminderStorage.GetOne(ctx, int64(reminderID))
	if err != nil {
		return
	}
	if entryID, ok := u.reminderMap[reminderID]; ok {
		u.cronJob.Remove(entryID)
	}
	u.addReminderJob(ctx, r)
}
func (u *userReminderJob) addReminderJobByID(ctx context.Context, reminderID uint) {
	r, err := u.reminderStorage.GetOne(ctx, int64(reminderID))
	if err != nil {
		return
	}
	u.addReminderJob(ctx, r)
}
func (u *userReminderJob) addReminderJob(ctx context.Context, r reminder.Reminder) {
	if !r.IsActive {
		log.Println("reminder is not active, skipping... ", r.ID)
		return
	}
	entryID, err := u.cronJob.AddFunc(r.Schedule, func() {
		u.pushMessageCronFunc(r)
	})
	if err != nil {
		log.Println("[ERROR] failed to add reminder job: ", r, "error", err)
		return
	}
	u.setReminderMap(r.ID, entryID)
}

func (u *userReminderJob) setReminderMap(reminderID uint, entryID cron.EntryID) {
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()
	u.reminderMap[reminderID] = entryID
}

func (u *userReminderJob) removeReminderMap(reminderID uint) {
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()
	delete(u.reminderMap, reminderID)
}

func (u *userReminderJob) pushMessageCronFunc(r reminder.Reminder) {
	err := u.pusher.PushMessage(context.Background(), r.WebhookType, r.Webhook, r.Message)
	if err != nil {
		log.Printf("[ERROR] failed to push message to webhook: %d, url: %s, message: %s",
			r.WebhookType, r.Webhook, r.Message)
	}
}
