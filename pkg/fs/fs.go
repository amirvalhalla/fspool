// Package fs contains functionalities of reader and writer for a file
package fs

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/cfgs"
	fsConfig "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/amirvalhalla/fspool/pkg/file"
	"github.com/amirvalhalla/fspool/pkg/reader"
	"github.com/amirvalhalla/fspool/pkg/writer"
	"github.com/google/uuid"
	"log"
	"path/filepath"
)

var (
	ErrFilesystemFilepathIsEmpty                 = errors.New("could not create new file system instance - path is empty")
	ErrFilesystemMemoryRentConflictWithFlushSize = errors.New("filesystem memory rent size always should be greater than flush size")
	ErrFilesystemCouldNotWrite                   = errors.New("filesystem could not write")
	ErrFilesystemWriterHasBeenClosed             = errors.New("writer instance of filesystem has been closed")
	ErrFilesystemCouldNotCloseWriter             = errors.New("filesystem could not close writer instance")
	ErrFilesystemReaderOccupying                 = errors.New("filesystem doesn't have any free reader")
	ErrFilesystemCouldNotReadData                = errors.New("filesystem could not read data")
	ErrFilesystemCouldNotReadAllData             = errors.New("filesystem could not read all data")
	ErrFilesystemReaderEmpty                     = errors.New("filesystem doesn't have any reader")
	ErrFilesystemCouldNotCloseReader             = errors.New("filesystem could not close reader")
)

type Filesystem interface {
	// Write will write or update raw data into file
	Write(rawData []byte, offset int64, seek int) error
	// GetWriterId return id of writer instance
	GetWriterId() (uuid.UUID, error)
	// CloseWriter will close writer of filesystem instance
	CloseWriter() error
	// ReadData func provides reading data from file by defining custom pos & seek option
	ReadData(offset int64, length int, seek int) ([]byte, error)
	// ReadAllData func provides reading all data from file
	ReadAllData() ([]byte, error)
	// GetReaderId return id of reader instance
	GetReaderId() (uuid.UUID, error)
	// CloseReader func provides close reader of filesystem instance
	CloseReader() error
	// GetReaderState return state of reader instance
	GetReaderState() (bool, error)
}

type filesystem struct {
	buff        []byte
	filePath    string
	dirPath     string
	config      fsConfig.FSConfiguration
	readerState bool // false means free and true means occupying
	reader      reader.FileReader
	writer      writer.FileWriter
}

// NewFilesystem provide new instance of filesystem with readers and writer based on your configuration
func NewFilesystem(fPath string, config fsConfig.FSConfiguration, file file.File, statFunc Stat, isNotExistFunc IsNotExist, mkdirAllFunc MkdirAll) (Filesystem, error) {
	var dirPath string
	var fWriter writer.FileWriter
	var fReader reader.FileReader

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
		fReader, _ = reader.NewFileReader(file)
	case cfgs.WOnly:
		fWriter, _ = writer.NewFileWriter(file)
	case cfgs.RW:
		fReader, _ = reader.NewFileReader(file)
		fWriter, _ = writer.NewFileWriter(file)
	}

	return &filesystem{
		buff:     make([]byte, config.MemoryRent),
		filePath: fPath,
		dirPath:  dirPath,
		config:   config,
		reader:   fReader,
		writer:   fWriter,
	}, nil
}

func (f *filesystem) Write(rawData []byte, offset int64, seek int) error {

	if f.writer == nil {
		return ErrFilesystemWriterHasBeenClosed
	}

	if err := f.writer.Write(rawData, offset, seek); err != nil {
		return ErrFilesystemCouldNotWrite
	}

	return nil
}

func (f *filesystem) GetWriterId() (uuid.UUID, error) {

	if f.writer == nil {
		return uuid.Nil, ErrFilesystemWriterHasBeenClosed
	}

	return f.writer.GetId(), nil
}

func (f *filesystem) CloseWriter() error {

	if f.writer == nil {
		return ErrFilesystemWriterHasBeenClosed
	}

	if err := f.writer.Close(); err != nil {
		return ErrFilesystemCouldNotCloseWriter
	}

	return nil
}

func (f *filesystem) ReadData(offset int64, length int, seek int) ([]byte, error) {

	if f.reader == nil {
		return nil, ErrFilesystemReaderEmpty
	}

	if f.readerState {
		return nil, ErrFilesystemReaderOccupying
	}

	f.readerState = true
	rawData, err := f.reader.ReadData(offset, length, seek)
	f.readerState = false

	if err != nil {
		log.Println(ErrFilesystemCouldNotReadData.Error())
		return nil, ErrFilesystemCouldNotReadData
	}

	return rawData, nil
}

func (f *filesystem) ReadAllData() ([]byte, error) {

	if f.reader == nil {
		return nil, ErrFilesystemReaderEmpty
	}

	if f.readerState {
		return nil, ErrFilesystemReaderOccupying
	}

	f.readerState = true
	rawData, err := f.reader.ReadAllData()
	f.readerState = false

	if err != nil {
		log.Println(ErrFilesystemCouldNotReadAllData.Error())
		return nil, ErrFilesystemCouldNotReadAllData
	}

	return rawData, nil
}

func (f *filesystem) GetReaderId() (uuid.UUID, error) {

	if f.reader == nil {
		return uuid.Nil, ErrFilesystemReaderEmpty
	}

	return f.reader.GetId(), nil
}

func (f *filesystem) CloseReader() error {

	if f.reader == nil {
		return ErrFilesystemReaderEmpty
	}

	if f.readerState {
		return ErrFilesystemReaderOccupying
	}

	if err := f.reader.Close(); err != nil {
		return ErrFilesystemCouldNotCloseReader
	}

	return nil
}

func (f *filesystem) GetReaderState() (bool, error) {
	if f.reader == nil {
		return false, ErrFilesystemReaderEmpty
	}
	return f.readerState, nil
}
