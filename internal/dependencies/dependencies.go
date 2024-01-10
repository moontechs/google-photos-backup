package dependencies

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moontechs/photos-backup/internal/account"
	"github.com/moontechs/photos-backup/internal/auth"
	"github.com/moontechs/photos-backup/internal/config"
	"github.com/moontechs/photos-backup/internal/db"
	"github.com/moontechs/photos-backup/internal/downloader"
	"github.com/moontechs/photos-backup/internal/files"
	"github.com/moontechs/photos-backup/internal/media_reader"
	"github.com/moontechs/photos-backup/internal/scanner"
)

type Factory interface {
	Create() (Dependencies, error)
}

type Dependencies struct {
	DbConnection         *db.Connection
	GoogleAuth           auth.Auth
	AuthRepository       auth.Repository
	Account              account.Account
	AccountRepository    account.Repository
	ScannerRepository    scanner.Repository
	UpdatesScanner       scanner.UpdatesScanner
	DownloadScheduler    downloader.Scheduler
	AccountLimiter       account.Limiter
	ConfigRepository     config.Repository
	Config               config.Config
	DownloaderRepository downloader.Repository
	Downloader           downloader.Downloader
	MediaReader          media_reader.Reader
	FilesRepository      files.Repository
	FilesManager         files.FilesManager
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
		DbConnection:         connection,
		AuthRepository:       auth.NewRepository(connection.DB),
		AccountRepository:    account.NewRepository(connection.DB),
		ScannerRepository:    scanner.NewRepository(connection.DB),
		DownloadScheduler:    downloader.NewScheduler(downloader.NewRepository(connection.DB)),
		AccountLimiter:       account.NewLimiter(account.NewRepository(connection.DB)),
		ConfigRepository:     config.NewRepository(connection.DB),
		DownloaderRepository: downloader.NewRepository(connection.DB),
		FilesRepository:      files.NewRepository(connection.DB),
	}

	deps.Config = config.NewConfig(deps.ConfigRepository)

	deps.Account = account.NewAccount(deps.AccountRepository)

	deps.GoogleAuth = auth.NewGoogleAuth(deps.AuthRepository)

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
