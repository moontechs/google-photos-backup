package main

import (
	"context"
	"fmt"
	"os"

	"google-backup/internal/dependencies"
	"google-backup/internal/scanner"

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

	scannerJob := scanner.NewScannerJob(
		dependencies.UpdatesScanner,
		dependencies.Config,
	)

	err = scannerJob.Run(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("run scanner job: %w", err))
	}
}
