package fs

import (
	"errors"
	"os"
)

var (
	ErrFileIsNotExists         = errors.New("file path doesn't exist")
	ErrDirectoryIsNotExists    = errors.New("directory path doesn't exist")
	ErrCouldNotCreateDirectory = errors.New("could not create directory")
)

func IsFileExists(filePath string, file File) error {
	if _, err := file.StatWithFilePath(filePath); err != nil {
		return ErrFileIsNotExists
	}

	return nil
}

func IsDirectoryExists(dirPath string, file File) error {
	if _, err := file.StatWithFilePath(dirPath); file.IsNotExist(err) {
		if err != nil {
			return ErrDirectoryIsNotExists
		}
	}
	return nil
}

func CreateDirectory(dirPath string, file File) error {
	if err := file.MkdirAll(dirPath, os.ModePerm); err != nil {
		return ErrCouldNotCreateDirectory
	}
	return nil
}
