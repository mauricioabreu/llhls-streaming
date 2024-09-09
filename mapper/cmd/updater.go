package cmd

import (
	"context"
	"mapper/config"
	"mapper/log"
	"mapper/originapi"
	"mapper/store"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var updaterCmd = &cobra.Command{
	Use:   "updater",
	Short: "Start the stream updater service",
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(
				config.New,
				store.NewRedisStore,
				originapi.New,
				log.New,
			),
			fx.Invoke(runUpdater),
		)
		app.Run()
	},
}

func runUpdater(lc fx.Lifecycle, apiClient *originapi.Client, redis *store.RedisStore, cfg *config.Config, logger *zap.SugaredLogger) {
	var cancelFunc context.CancelFunc

	var ctx context.Context

	lc.Append(fx.Hook{
		OnStart: func(startContext context.Context) error {
			logger.Info("Starting stream updater...")
			ctx, cancelFunc = context.WithCancel(context.Background())

			go func() {
				ticker := time.NewTicker(time.Duration(cfg.UpdateInterval) * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						streams, err := apiClient.ListStreams(cfg.APIHosts)
						if err != nil {
							logger.Errorw("Error listing streams", "error", err)
							continue
						}

						logger.Infow("Updating streams in Redis", "streams", streams)

						err = redis.UpdateStreams(ctx, streams, cfg.StreamTTL)
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

func init() {
	rootCmd.AddCommand(updaterCmd)
}
