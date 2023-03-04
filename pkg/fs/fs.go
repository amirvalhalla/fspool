// Package fs contains functionalities of reader and writer for a file
package fs

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/cfgs"
	fsConfig "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/amirvalhalla/fspool/pkg/reader"
	"github.com/amirvalhalla/fspool/pkg/writer"
	"io"
	"io/fs"
	"os"
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

// File override os.File interface of golang with RW interfaces
type File interface {
	io.Writer
	io.WriterAt
	io.WriterTo
	io.WriteCloser
	io.WriteSeeker
	io.ByteWriter
	io.StringWriter
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	io.ByteReader
	io.ReaderFrom
	io.RuneReader
	io.RuneScanner
	io.ReadSeekCloser
	io.ReadSeeker
	fs.FileInfo
	Stat() (os.FileInfo, error)
	StatWithFilePath(filePath string) (os.FileInfo, error)
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
}

func NewFilesystem(fPath string, config fsConfig.FSConfiguration, file File) (Filesystem, error) {
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
		if err := IsFileExists(fPath, file); err != nil {
			return nil, err
		}
	} else if err := IsFileExists(fPath, file); err != nil {
		dirPath = filepath.Dir(fPath)
		if err := IsDirectoryExists(fPath, file); err != nil {
			if err := CreateDirectory(dirPath, file); err != nil {
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
