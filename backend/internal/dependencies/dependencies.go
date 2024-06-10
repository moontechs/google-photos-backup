package dependencies

import (
	"fmt"
	"log"

	"google-backup/internal/account"
	"google-backup/internal/auth"
	"google-backup/internal/db"
	"google-backup/internal/downloader"
	"google-backup/internal/files"
	"google-backup/internal/google_client"
	"google-backup/internal/photos"
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
	PhotosUpdatesScanner   scanner.UpdatesScanner
	DownloadScheduler      downloader.Scheduler
	ClienLimiter           google_client.Limiter
	SettingsRepository     settings.Repository
	SettingsInitializer    settings.SettingsInitializer
	DownloaderRepository   downloader.Repository
	Downloader             downloader.Downloader
	FilesRepository        files.Repository
	FilesManager           files.FilesManager
	GoogleClientRepository google_client.Repository
	ReaderCreater          photos.ReaderCreater
	PhotosScanner          scanner.AccountScanner
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
		SettingsRepository:     settings.NewRepository(connection.DB),
		DownloaderRepository:   downloader.NewRepository(connection.DB),
		FilesRepository:        files.NewRepository(connection.DB),
		GoogleClientRepository: google_client.NewRepository(connection.DB),
	}

	deps.ClienLimiter = google_client.NewLimiter(deps.GoogleClientRepository)

	deps.SettingsInitializer = settings.NewSettings(deps.SettingsRepository)

	deps.Account = account.NewAccount(deps.AccountRepository, deps.GoogleClientRepository)

	deps.GoogleAuth = auth.NewGoogleAuth(deps.AuthRepository, deps.GoogleClientRepository)

	deps.FilesManager = files.NewFilesManager(deps.FilesRepository)

	deps.ReaderCreater = photos.NewReaderCreater(
		deps.Account,
		deps.GoogleAuth,
	)

	deps.PhotosScanner = scanner.NewPhotosScanner(
		deps.ScannerRepository,
		deps.ReaderCreater,
		deps.ClienLimiter,
		deps.DownloadScheduler,
	)

	deps.PhotosUpdatesScanner = scanner.NewUpdatesScanner(
		deps.AccountRepository,
		deps.SettingsRepository,
		deps.PhotosScanner,
	)

	// deps.Downloader = downloader.NewDownloader(
	// 	deps.DownloaderRepository,
	// 	&http.Client{},
	// 	deps.MediaReader,
	// 	deps.AccountLimiter,
	// 	deps.FilesManager,
	// )

	return deps, nil
}
