package dependencies

import (
	"fmt"
	"log"
	"net/http"

	"google-backup/internal/account"
	"google-backup/internal/auth"
	"google-backup/internal/db"
	"google-backup/internal/downloader"
	"google-backup/internal/files"
	"google-backup/internal/google_client"
	"google-backup/internal/media_reader"
	"google-backup/internal/scanner"
	"google-backup/internal/settings"
)

type Factory interface {
	Create() (Dependencies, error)
}

type Dependencies struct {
	DbConnection           *db.Connection
	GoogleAuth             auth.Auth
	AuthRepository         auth.Repository
	Account                account.Account
	AccountRepository      account.Repository
	ScannerRepository      scanner.Repository
	UpdatesScanner         scanner.UpdatesScanner
	DownloadScheduler      downloader.Scheduler
	AccountLimiter         account.Limiter
	SettingsRepository     settings.Repository
	SettingsInitializer    settings.SettingsInitializer
	DownloaderRepository   downloader.Repository
	Downloader             downloader.Downloader
	MediaReader            media_reader.Reader
	FilesRepository        files.Repository
	FilesManager           files.FilesManager
	GoogleClientRepository google_client.Repository
}

type factory struct{}

func NewFactory() factory {
	return factory{}
}

func (f factory) Create() (Dependencies, error) {
	connection, err := db.NewConnection()
	if err != nil {
		log.Print(err)

		return Dependencies{}, fmt.Errorf("new connection: %w", err)
	}

	deps := Dependencies{
		DbConnection:           connection,
		AuthRepository:         auth.NewRepository(connection.DB),
		AccountRepository:      account.NewRepository(connection.DB),
		ScannerRepository:      scanner.NewRepository(connection.DB),
		DownloadScheduler:      downloader.NewScheduler(downloader.NewRepository(connection.DB)),
		AccountLimiter:         account.NewLimiter(account.NewRepository(connection.DB)),
		SettingsRepository:     settings.NewRepository(connection.DB),
		DownloaderRepository:   downloader.NewRepository(connection.DB),
		FilesRepository:        files.NewRepository(connection.DB),
		GoogleClientRepository: google_client.NewRepository(connection.DB),
	}

	deps.SettingsInitializer = settings.NewSettings(deps.SettingsRepository)

	deps.Account = account.NewAccount(deps.AccountRepository)

	deps.GoogleAuth = auth.NewGoogleAuth(deps.AuthRepository, deps.GoogleClientRepository)

	deps.FilesManager = files.NewFilesManager(deps.FilesRepository)

	deps.MediaReader = media_reader.NewMediaReader(
		deps.Account,
		deps.GoogleAuth,
		deps.AccountLimiter,
	)

	deps.UpdatesScanner = scanner.NewUpdatesScanner(
		deps.ScannerRepository,
		deps.DownloadScheduler,
		deps.AccountLimiter,
		deps.MediaReader,
	)

	deps.Downloader = downloader.NewDownloader(
		deps.DownloaderRepository,
		&http.Client{},
		deps.MediaReader,
		deps.AccountLimiter,
		deps.FilesManager,
	)

	return deps, nil
}
