// Package fs contains functionalities of reader and writer for a file
package fs

import (
	"github.com/amirvalhalla/fspool/pkg/reader"
	"github.com/amirvalhalla/fspool/pkg/writer"
	"sync"
)

type ConfigType uint8

const (
	FsPoolConfigurationType ConfigType = 0
	FsConfigurationType     ConfigType = 1
)

// configurationOverride: for detecting which Configuration Filesystem instance using , if client enter FSConfiguration then it will override on FSPoolConfiguration
// by default if client don't fill FSConfiguration struct and just pass it , the system will use FSPoolConfiguration as general config
// tip : entering FsPerm & FilePath in FSConfiguration is required!
type filesystem struct {
	filePath   string
	fsConfig   interface{}
	configType ConfigType
	mu         sync.RWMutex
	fileReader reader.FileReader
	fileWriter writer.FileWriter
}

//// Filesystem interface is top layer of FileReader and FileWriter to handle reading or writing easier
//type Filesystem interface {
//}
//
//// NewFilesystem will return a new instance from Filesystem interface
//func NewFilesystem(filePath string, fsPoolConfig fsPoolCfgs.FSPoolConfiguration, fsConfig fsCfgs.FSConfiguration) (Filesystem, error) {
//	var err error
//	var configType ConfigType
//	var config interface{}
//	var fileReader reader.FileReader
//	var fileWriter writer.FileWriter
//
//	if fsConfig.FlushSize == 0 && fsConfig.FlushDuration == 0 && fsConfig.MemoryRent == 0 && fsConfig.ReaderLimit == 0 {
//		configType = FsPoolConfigurationType
//		config = fsPoolConfig
//	} else {
//		configType = FsConfigurationType
//		config = fsConfig
//	}
//
//	if configType == FsPoolConfigurationType {
//		if fileWriter, fileReader, err = newRwBasedOnCfgs(filePath, fsPoolConfig.Perm); err != nil {
//			return nil, err
//		}
//	} else if configType == FsConfigurationType {
//		if fileWriter, fileReader, err = newRwBasedOnCfgs(filePath, fsConfig.Perm); err != nil {
//			return nil, err
//		}
//	}
//
//	return &filesystem{
//		fsConfig:   config,
//		configType: configType,
//		mu:         sync.RWMutex{},
//		fileWriter: fileWriter,
//		fileReader: fileReader,
//	}, nil
//}
//
//func newRwBasedOnCfgs(filePath string, perm cfgs.FSPerm) (writer.FileWriter, reader.FileReader, error) {
//	var err error
//	var fileReader reader.FileReader
//	var fileWriter writer.FileWriter
//
//	switch perm {
//	case cfgs.ROnly:
//		if fileReader, err = reader.NewFileReader(filePath); err != nil {
//			return nil, nil, err
//		}
//	case cfgs.WOnly:
//		if fileWriter, err = writer.NewFileWriter(filePath); err != nil {
//			return nil, nil, err
//		}
//	case cfgs.RW:
//		if fileWriter, err = writer.NewFileWriter(filePath); err != nil {
//			return nil, nil, err
//		}
//		if fileReader, err = reader.NewFileReader(filePath); err != nil {
//			return nil, nil, err
//		}
//	}
//
//	return fileWriter, fileReader, nil
//}
