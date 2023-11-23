package cmd

import (
	"context"
	"go.uber.org/fx"
	"taskbot/admin/server"
	"taskbot/cron"
	"taskbot/di"
)

func Execute() error {
	app := fx.New(
		di.Module,
		fx.Invoke(startCronJob),
		fx.Invoke(startAdminServer),
	)
	app.Run()
	return nil
}

func startCronJob(lc fx.Lifecycle, job cron.CronJob) {
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
