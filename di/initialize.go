package di

import (
	"github.com/asaskevich/EventBus"
	"go-reminder-bot/admin/server"
	"go-reminder-bot/cron"
	"go-reminder-bot/pkg/db"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/user"
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
	provideUserStorage,
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

func provideUserStorage(db *gorm.DB) user.Storage {
	return user.NewStorage(db)
}

func provideServer(storage reminder.Storage, userStorage user.Storage, eventBus EventBus.Bus, pusher pusher.Pusher) server.Server {
	handler := server.NewHandler(storage, userStorage, eventBus, pusher)
	return server.NewServer(*handler)
}

func provideCronjob(storage reminder.Storage, pusher pusher.Pusher, eventBus EventBus.Bus) cron.CronJob {
	return cron.NewCron(storage, pusher, eventBus)
}

func provideEventBus() EventBus.Bus {
	return EventBus.New()
}
