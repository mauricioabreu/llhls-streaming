package cmd

import (
	"mapper/api"
	"mapper/config"
	"mapper/log"
	"mapper/store"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API service",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			fx.Provide(
				config.New,
				store.NewRedisStore,
				api.New,
				log.New,
			),
			fx.Invoke(runAPI),
		)
		return app.Start(cmd.Context())
	},
}

func runAPI(ginAPI *api.API, logger *zap.SugaredLogger) {
	router := ginAPI.SetupRouter()
	logger.Info("Starting API server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Fatalw("Failed to start API server", "error", err)
	}
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
