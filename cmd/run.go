package cmd

import (
	"context"
	"github.com/joho/godotenv"
	"go-reminder-bot/admin/server"
	"go-reminder-bot/cron"
	"go-reminder-bot/di"
	"go.uber.org/fx"
	"log"
)

func Execute() error {
	app := fx.New(
		fx.Invoke(loadEnv),
		di.Module,
		fx.Invoke(startCronJob),
		fx.Invoke(startAdminServer),
	)
	app.Run()
	return nil
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return err
}

func startCronJob(lc fx.Lifecycle, job cron.UserReminderJob) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			job.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			job.Stop(ctx)
			return nil
		},
	})
}

func startAdminServer(lc fx.Lifecycle, server server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Stop(ctx)
			return nil
		},
	})
}
