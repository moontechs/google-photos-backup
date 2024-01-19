package main

import (
	"context"
	"fmt"
	"os"

	"google-backup/internal/dependencies"
	"google-backup/internal/downloader"

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

	downloaderJob := downloader.NewDownloaderJob(
		dependencies.Downloader,
		dependencies.SettingsRepository,
	)

	err = downloaderJob.Run(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("run downloader job: %w", err))
	}
}
