package main

import (
	"context"
	"log"
	"mapper/api"
	"mapper/config"
	"mapper/store"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func UpdateStreams(lc fx.Lifecycle, apiClient *api.Client, store *store.RedisStore, config *config.Config, logger *zap.SugaredLogger) {
	var cancelFunc context.CancelFunc
	var ctx context.Context

	lc.Append(fx.Hook{
		OnStart: func(startContext context.Context) error {
			logger.Info("Starting stream updater...")
			ctx, cancelFunc = context.WithCancel(context.Background())

			go func() {
				ticker := time.NewTicker(time.Duration(config.UpdateInterval) * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						streams, err := apiClient.ListStreams(config.Hosts)
						if err != nil {
							logger.Errorw("Error listing streams", "error", err)
							continue
						}

						logger.Infow("Updating streams in Redis", "streams", streams)

						err = store.UpdateStreams(ctx, streams, config.StreamTTL)
						if err != nil {
							logger.Errorw("Error updating streams in Redis", "error", err)
						}
					case <-ctx.Done():
						logger.Info("Stream updater stopped")
						return
					}
				}
			}()

			return nil
		},
		OnStop: func(stopContext context.Context) error {
			logger.Info("Stopping stream updater...")
			if cancelFunc != nil {
				cancelFunc()
			}
			return nil
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(
			config.New,
			store.NewRedisStore,
			api.New,
			log.New,
		),
		fx.Invoke(UpdateStreams),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}
