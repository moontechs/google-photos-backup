package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moontechs/photos-backup/internal/cron"
	"github.com/moontechs/photos-backup/internal/dependencies"
	"github.com/moontechs/photos-backup/internal/downloader"
	"github.com/moontechs/photos-backup/internal/handlers"
	"github.com/moontechs/photos-backup/internal/scanner"
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

	err = dependencies.Config.Init()
	if err != nil {
		log.Fatal(fmt.Errorf("init config: %w", err))
	}

	cronRunner, err := createCron(dependencies)
	if err != nil {
		log.Fatal(fmt.Errorf("create cron runner: %w", err))
	}

	ginServer := createGinServer(dependencies)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go cronRunner.Start(ctx)

	err = ginServer.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("server run: %v", err)
	}
}

func createGinServer(dependencies dependencies.Dependencies) *gin.Engine {
	ginEngine := gin.Default()

	ginEngine.GET("/auth/google/:clientName", func(c *gin.Context) {
		var client struct {
			ClientName string `uri:"clientName" binding:"required"`
		}

		if err := c.ShouldBindUri(&client); err != nil {
			c.JSON(400, gin.H{"msg": err})

			return
		}

		c.Redirect(http.StatusFound, dependencies.GoogleAuth.GetRedirectUrl(client.ClientName))
	})

	ginEngine.GET("/auth/google/callback/:clientName", handlers.NewGoogleCallbackHandler(
		dependencies.AccountRepository,
		dependencies.GoogleAuth,
	).Handle)

	ginEngine.POST("/api/v1/rescan", handlers.NewRescanHandler(
		dependencies.AccountRepository,
		scanner.NewScheduler(dependencies.ScannerRepository),
	).Handle)

	return ginEngine
}

func createCron(dependencies dependencies.Dependencies) (cron.Runner, error) {
	return cron.NewCron(
		[]cron.Job{
			scanner.NewScannerJob(
				dependencies.UpdatesScanner,
				dependencies.Config,
			),
			downloader.NewDownloaderJob(
				dependencies.Downloader,
				dependencies.Config,
			),
		},
	), nil
}
