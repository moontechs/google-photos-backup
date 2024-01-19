package files

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"google-backup/internal/media"
)

type FilesManager interface {
	SaveDownloadError(email string, mediaItemId, message string) error
	SaveFileMeta(email string, fileMeta FileMeta) error
	FileExists(email string, mediaItem media.MediaItem) (bool, error)
	GenerateFilePathName(email string, mediaItem media.MediaItem) (string, error)
	EqualHash(filePathName string, reader io.Reader) (bool, error)
	AddRootFolderToPath(path string) string
	CreateFolderIfDoesNotExist(filePathName string) error
	UpdateCreationTime(filePathName string, creationTime string) error
}

type files struct {
	repository Repository
}

type FileMeta struct {
	FilePathName string          `json:"file_path_name"`
	MediaItem    media.MediaItem `json:"media_item"`
}

func NewFilesManager(repository Repository) files {
	return files{repository: repository}
}

func (f files) SaveDownloadError(email string, mediaItemId string, message string) error {
	return f.repository.SaveDownloadError(email, mediaItemId, message)
}

func (f files) SaveFileMeta(email string, fileMeta FileMeta) error {
	fileMetaJson, err := json.Marshal(fileMeta)
	if err != nil {
		return fmt.Errorf("marshal media item: %w", err)
	}

	return f.repository.SaveFileMeta(email, []byte(fileMeta.MediaItem.ID), fileMetaJson)
}

func (f files) FileExists(email string, mediaItem media.MediaItem) (bool, error) {
	filePathName, err := f.GenerateFilePathName(email, mediaItem)
	if err != nil {
		return false, fmt.Errorf("generate file path name: %w", err)
	}

	fileExists, err := f.fileExistsOnDisk(f.AddRootFolderToPath(filePathName))
	if err != nil {
		return false, fmt.Errorf("file exists on disk: %w", err)
	}

	if !fileExists {
		return false, nil
	}

	fileMetaJson, err := f.repository.GetFileMeta(email, []byte(mediaItem.ID))
	if err != nil {
		return false, fmt.Errorf("get file meta: %w", err)
	}

	if fileMetaJson == nil {
		return false, nil
	}

	var fileMeta FileMeta
	err = json.Unmarshal(fileMetaJson, &fileMeta)
	if err != nil {
		return false, fmt.Errorf("unmarshal file meta: %w", err)
	}

	return fileMeta.FilePathName == filePathName, nil
}

func (f files) GenerateFilePathName(email string, mediaItem media.MediaItem) (string, error) {
	creationTime, err := time.Parse(time.RFC3339, mediaItem.MediaMetadata.CreationTime)
	if err != nil {
		return "", fmt.Errorf("parse creation time: %w", err)
	}

	return email + "/" + strconv.Itoa(creationTime.Year()) + "/" + strconv.Itoa(int(creationTime.Month())) + "/" + mediaItem.Filename, nil
}

func (f files) EqualHash(filePathName string, reader io.Reader) (bool, error) {
	hash := md5.New()

	file, err := os.Open(f.AddRootFolderToPath(filePathName))
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("copy file: %w", err)
	}

	existingFileHash := string(hash.Sum(nil))

	if _, err := io.Copy(hash, reader); err != nil {
		return false, fmt.Errorf("copy file: %w", err)
	}

	newFileHash := string(hash.Sum(nil))

	return existingFileHash == newFileHash, nil
}

func (f files) AddRootFolderToPath(path string) string {
	// TODO get from config
	return "/Users/michael/github.com/moontechs/photos-backup/downloads" + "/" + path
}

func (f files) CreateFolderIfDoesNotExist(filePathName string) error {
	directory := filepath.Dir(filePathName)

	return os.MkdirAll(f.AddRootFolderToPath(directory), os.ModePerm)
}

func (f files) UpdateCreationTime(filePathName string, creationTime string) error {
	creationTimeParsed, err := time.Parse(time.RFC3339, creationTime)
	if err != nil {
		return fmt.Errorf("parse creation time: %w", err)
	}

	return os.Chtimes(f.AddRootFolderToPath(filePathName), creationTimeParsed, creationTimeParsed)
}

func (f files) fileExistsOnDisk(filePathName string) (bool, error) {
	_, err := os.Stat(filePathName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("stat file: %w", err)
	}

	return true, nil
}
