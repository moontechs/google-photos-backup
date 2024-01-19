package main

import (
	"fmt"
	"os"

	"google-backup/internal/cron"
	"google-backup/internal/dependencies"
	"google-backup/internal/downloader"
	"google-backup/internal/handlers"
	"google-backup/internal/scanner"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if os.Getenv("PRODUCTION") == "true" {
		log.SetLevel(log.ErrorLevel)
	}

	dependencies, err := dependencies.NewFactory().Create()
	if err != nil {
		log.Fatal(fmt.Errorf("create depdendencies: %w", err))
	}
	defer dependencies.DbConnection.Close()

	err = dependencies.SettingsInitializer.Init()
	if err != nil {
		log.Fatal(fmt.Errorf("init config: %w", err))
	}

	// cronRunner, err := createCron(dependencies)
	if err != nil {
		log.Fatal(fmt.Errorf("create cron runner: %w", err))
	}

	ginServer := createGinServer(dependencies)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// go cronRunner.Start(ctx)

	err = ginServer.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("server run: %v", err)
	}
}

func createGinServer(dependencies dependencies.Dependencies) *gin.Engine {
	ginEngine := gin.Default()

	// ginEngine.GET("/auth/google/callback/:clientId", handlers.NewGoogleCallbackHandler(
	// 	dependencies.AccountRepository,
	// 	dependencies.GoogleAuth,
	// ).Handle)

	// ginEngine.POST("/api/v1/rescan", handlers.NewRescanHandler(
	// 	dependencies.AccountRepository,
	// 	scanner.NewScheduler(dependencies.ScannerRepository),
	// ).Handle)

	ginEngine.Any("/api/v1/clients/:clientId", handlers.NewClientsApiHandler(
		dependencies.GoogleClientRepository,
		dependencies.SettingsRepository,
	).Handle)

	ginEngine.Any("/api/v1/clients/:clientId/redirect-url", handlers.NewGoogleRedirectUrlHandler(
		dependencies.GoogleAuth,
	).Handle)

	return ginEngine
}

func createCron(dependencies dependencies.Dependencies) (cron.Runner, error) {
	return cron.NewCron(
		[]cron.Job{
			scanner.NewScannerJob(
				dependencies.UpdatesScanner,
				dependencies.SettingsRepository,
			),
			downloader.NewDownloaderJob(
				dependencies.Downloader,
				dependencies.SettingsRepository,
			),
		},
	), nil
}
