// Package fs contains functionalities of reader and writer for a file
package fs

import (
	fsCfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	fsPoolCfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fspool"
	"github.com/amirvalhalla/fspool/pkg/reader"
	"github.com/amirvalhalla/fspool/pkg/writer"
	"sync"
)

type ConfigurationOverrideType uint8

const (
	FsPoolConfiguration ConfigurationOverrideType = 0
	FsConfiguration     ConfigurationOverrideType = 1
)

// configurationOverride: for detecting which Configuration Filesystem instance using , if client enter FSConfiguration then it will override on FSPoolConfiguration
// by default if client don't fill FSConfiguration struct and just pass it , the system will use FSPoolConfiguration as general config
// tip : entering FsPerm & FilePath in FSConfiguration is required!
type filesystem struct {
	fsPoolCfgs            fsPoolCfgs.FSPoolConfiguration
	fsCfgs                fsCfgs.FSConfiguration
	configurationOverride ConfigurationOverrideType
	mu                    sync.RWMutex
	fileReader            reader.FileReader
	fileWriter            writer.FileWriter
}

// Filesystem interface is top layer of FileReader and FileWriter to handle reading or writing easier
type Filesystem interface {
}

// NewFilesystem will return a new instance from Filesystem interface
func NewFilesystem(fsPoolConfig fsPoolCfgs.FSPoolConfiguration, fsConfig fsCfgs.FSConfiguration) (Filesystem, error) {
	var err error
	var configurationOverride ConfigurationOverrideType
	var fileReader reader.FileReader
	var fileWriter writer.FileWriter

	switch fsConfig.FsPerm {
	case fsCfgs.ROnly:
		if fileReader, err = reader.NewFileReader(fsConfig.FilePath); err != nil {
			return &filesystem{}, err
		}
	case fsCfgs.WOnly:
		if fileWriter, err = writer.NewFileWriter(fsConfig.FilePath); err != nil {
			return &filesystem{}, err
		}
	case fsCfgs.RW:
		if fileWriter, err = writer.NewFileWriter(fsConfig.FilePath); err != nil {
			return &filesystem{}, err
		}
		if fileReader, err = reader.NewFileReader(fsConfig.FilePath); err != nil {
			return &filesystem{}, err
		}
	}

	if fsConfig.FlushSize == 0 && fsConfig.FlushDuration == 0 && fsConfig.MemoryRent == 0 && fsConfig.ReaderLimit == 0 {
		configurationOverride = FsPoolConfiguration
	} else {
		configurationOverride = FsConfiguration
	}

	return &filesystem{
		fsPoolCfgs:            fsPoolConfig,
		fsCfgs:                fsConfig,
		configurationOverride: configurationOverride,
		mu:                    sync.RWMutex{},
		fileWriter:            fileWriter,
		fileReader:            fileReader,
	}, nil
}
