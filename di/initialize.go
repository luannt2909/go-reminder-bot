package di

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"taskbot/admin/server"
	"taskbot/cron"
	"taskbot/pkg/db"
	"taskbot/pkg/pusher"
	"taskbot/pkg/task"
	"taskbot/pkg/xservice/ggchat"
)

var Module = fx.Provide(
	provideSqlDB,
	provideTaskStorage,
	provideServer,
	provideGGChatService,
	providePusher,
	provideCronjob,
)

func provideSqlDB() (*gorm.DB, error) {
	return db.InitSQLiteDB()
}

func provideTaskStorage(db *gorm.DB) task.Storage {
	return task.NewStorage(db)
}

func provideGGChatService() ggchat.Service {
	return ggchat.NewService()
}

func providePusher(ggChatSvc ggchat.Service) pusher.Pusher {
	return pusher.NewPusher(ggChatSvc)
}

func provideServer(storage task.Storage) server.Server {
	handler := server.NewHandler(storage)
	return server.NewServer(*handler)
}

func provideCronjob(storage task.Storage, pusher pusher.Pusher) cron.CronJob {
	return cron.NewCron(storage, pusher)
}
