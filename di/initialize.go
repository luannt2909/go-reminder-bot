package di

import (
	"github.com/asaskevich/EventBus"
	"go-reminder-bot/admin/server"
	"go-reminder-bot/cron"
	"go-reminder-bot/pkg/config"
	"go-reminder-bot/pkg/db"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/token"
	"go-reminder-bot/pkg/user"
	"go-reminder-bot/pkg/xservice/discordsvc"
	"go-reminder-bot/pkg/xservice/ggchat"
	"go-reminder-bot/pkg/xservice/msteams"
	"go-reminder-bot/pkg/xservice/slacksvc"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Provide(
	provideConfig,
	provideDB,
	provideReminderStorage,
	provideServer,
	provideGGChatService,
	provideMSTeamsService,
	provideSlackService,
	provideDiscordService,
	providePusher,
	provideEventBus,
	provideTokenizer,
	provideUserStorage,
	provideUserReminderCronJob,
	provideReminderPusher,
)

func provideConfig() (config.Config, error) {
	return config.LoadEnv()
}
func provideDB(cfg config.Config) (*gorm.DB, error) {
	return db.InitDatabase(cfg.DBConfig)
}

func provideReminderStorage(db *gorm.DB) reminder.Storage {
	return reminder.NewStorage(db)
}

func provideGGChatService() ggchat.Service {
	return ggchat.NewService()
}

func provideDiscordService() discordsvc.Service {
	return discordsvc.NewService()
}

func provideMSTeamsService() msteams.Service {
	return msteams.NewService()
}

func provideSlackService() slacksvc.Service {
	return slacksvc.NewService()
}

func providePusher(ggChatSvc ggchat.Service, msTeamsSvc msteams.Service,
	slackSvc slacksvc.Service, discordSvc discordsvc.Service) pusher.Pusher {
	return pusher.NewPusher(ggChatSvc, msTeamsSvc, slackSvc, discordSvc)
}

func provideReminderPusher(p pusher.Pusher) pusher.ReminderPusher {
	return pusher.NewReminderPusher(p)
}

func provideUserStorage(db *gorm.DB) user.Storage {
	return user.NewStorage(db)
}

func provideServer(cfg config.Config, storage reminder.Storage, userStorage user.Storage, eventBus EventBus.Bus,
	pusher pusher.Pusher, tokenizer token.Tokenizer) server.Server {
	handler := server.NewHandler(cfg, storage, userStorage, eventBus, pusher, tokenizer)
	return server.NewServer(*handler, server.AuthenticateUserHandler(tokenizer, userStorage))
}

func provideUserReminderCronJob(userStorage user.Storage, reminderStorage reminder.Storage, pusher pusher.ReminderPusher, eventBus EventBus.Bus) cron.UserReminderJob {
	return cron.NewUserReminderCronJob(userStorage, reminderStorage, pusher, eventBus)
}

func provideEventBus() EventBus.Bus {
	return EventBus.New()
}

func provideTokenizer(cfg config.Config) token.Tokenizer {
	return token.NewJwtTokenizer([]byte(cfg.JwtSigningKey))
}
