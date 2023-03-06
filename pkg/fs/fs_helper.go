package fs

import (
	"errors"
	"io/fs"
	"os"
)

type Stat func(path string) (fs.FileInfo, error)
type IsNotExist func(err error) bool
type MkdirAll func(path string, mode fs.FileMode) error

var (
	ErrFileIsNotExists         = errors.New("file path doesn't exist")
	ErrDirectoryIsNotExists    = errors.New("directory path doesn't exist")
	ErrCouldNotCreateDirectory = errors.New("could not create directory")
)

// IsFileExists checks file exist or not
func IsFileExists(path string, statFunc Stat) error {
	if _, err := statFunc(path); err != nil {
		return ErrFileIsNotExists
	}
	return nil
}

// IsDirectoryExists checks directory exist or not
func IsDirectoryExists(path string, statFunc Stat, isNotExistFunc IsNotExist) error {
	if _, err := statFunc(path); isNotExistFunc(err) {
		if err != nil {
			return ErrDirectoryIsNotExists
		}
	}
	return nil
}

// CreateDirectory will create directory recursively
func CreateDirectory(path string, mkdirAllFunc MkdirAll) error {
	if err := mkdirAllFunc(path, os.ModePerm); err != nil {
		return ErrCouldNotCreateDirectory
	}
	return nil
}
