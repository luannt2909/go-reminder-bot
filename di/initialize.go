package di

import (
	"github.com/asaskevich/EventBus"
	"go-reminder-bot/admin/server"
	"go-reminder-bot/cron"
	"go-reminder-bot/pkg/db"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/xservice/ggchat"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Provide(
	provideSqlDB,
	provideReminderStorage,
	provideServer,
	provideGGChatService,
	providePusher,
	provideCronjob,
	provideEventBus,
)

func provideSqlDB() (*gorm.DB, error) {
	return db.InitSQLiteDB()
}

func provideReminderStorage(db *gorm.DB) reminder.Storage {
	return reminder.NewStorage(db)
}

func provideGGChatService() ggchat.Service {
	return ggchat.NewService()
}

func providePusher(ggChatSvc ggchat.Service) pusher.Pusher {
	return pusher.NewPusher(ggChatSvc)
}

func provideServer(storage reminder.Storage, eventBus EventBus.Bus) server.Server {
	handler := server.NewHandler(storage, eventBus)
	return server.NewServer(*handler)
}

func provideCronjob(storage reminder.Storage, pusher pusher.Pusher, eventBus EventBus.Bus) cron.CronJob {
	return cron.NewCron(storage, pusher, eventBus)
}

func provideEventBus() EventBus.Bus {
	return EventBus.New()
}
