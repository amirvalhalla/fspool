// Package fs contains functionalities of reader and writer for a file
package fs

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/cfgs"
	fsConfig "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/amirvalhalla/fspool/pkg/file"
	"github.com/amirvalhalla/fspool/pkg/reader"
	"github.com/amirvalhalla/fspool/pkg/writer"
	"path/filepath"
)

var (
	ErrFilesystemFilepathIsEmpty                 = errors.New("could not create new file system instance - path is empty")
	ErrFilesystemMemoryRentConflictWithFlushSize = errors.New("memory rent size always should be greater than flush size")
)

type Filesystem interface {
}

type filesystem struct {
	filePath   string
	dirPath    string
	config     fsConfig.FSConfiguration
	fileReader reader.FileReader
	fileWriter writer.FileWriter
}

func NewFilesystem(fPath string, config fsConfig.FSConfiguration, file file.File, statFunc Stat, isNotExistFunc IsNotExist, mkdirAllFunc MkdirAll) (Filesystem, error) {
	var dirPath string
	var fReader reader.FileReader
	var fWriter writer.FileWriter

	if fPath == "" || len(fPath) <= 0 {
		return nil, ErrFilesystemFilepathIsEmpty
	}

	if config.FlushType == cfgs.FlushBySize {
		if config.MemoryRent < config.FlushSize {
			return nil, ErrFilesystemMemoryRentConflictWithFlushSize
		}
	}

	if config.Perm == cfgs.ROnly {
		if err := IsFileExists(fPath, statFunc); err != nil {
			return nil, err
		}
	} else if err := IsFileExists(fPath, statFunc); err != nil {
		dirPath = filepath.Dir(fPath)
		if err := IsDirectoryExists(fPath, statFunc, isNotExistFunc); err != nil {
			if err := CreateDirectory(dirPath, mkdirAllFunc); err != nil {
				return nil, err
			}
		}
	}

	switch config.Perm {
	case cfgs.ROnly:
		fReader = reader.NewFileReader(file)
	case cfgs.WOnly:
		fWriter = writer.NewFileWriter(file)
	case cfgs.RW:
		fReader = reader.NewFileReader(file)
		fWriter = writer.NewFileWriter(file)
	}

	return &filesystem{
		filePath:   fPath,
		dirPath:    dirPath,
		config:     config,
		fileReader: fReader,
		fileWriter: fWriter,
	}, nil
}
