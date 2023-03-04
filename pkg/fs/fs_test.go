package fs

import (
	cfgs2 "github.com/amirvalhalla/fspool/pkg/cfgs"
	cfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFilesystem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath(someFilePath).Return(mockFile, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.Nil(t, err)
}

func TestNewFilesystem_With_ROnly_Perm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath(someFilePath).Return(mockFile, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.Perm = cfgs2.ROnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.Nil(t, err)
}

func TestNewFilesystem_With_WOnly_Perm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath(someFilePath).Return(mockFile, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.Perm = cfgs2.WOnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.Nil(t, err)
}

func TestNewFilesystem_EmptyFilePath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem("", fsConfig, mockFile)

	assert.EqualError(t, err, ErrFilesystemFilepathIsEmpty.Error())
}

func TestNewFilesystem_ConflictInMemoryRentSizeWithFlushSize(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.FlushSize = 60 * 1024 * 1024

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.EqualError(t, err, ErrFilesystemMemoryRentConflictWithFlushSize.Error())
}

func TestNewFilesystem_FilePathIsNotExists_With_ReadOnly_Permission(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.ROnly

	mockFile.EXPECT().StatWithFilePath(someFilePath).Return(mockFile, ErrFileIsNotExists).Times(1)

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.EqualError(t, err, ErrFileIsNotExists.Error())
}

func TestNewFilesystem_FilePathIsNotExists_CouldNotCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")
	someDirPath := filepath.Join("/test")

	mockFile := NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	mockFile.EXPECT().StatWithFilePath(someFilePath).Return(nil, ErrFileIsNotExists).MaxTimes(2)
	mockFile.EXPECT().IsNotExist(ErrFileIsNotExists).Return(true).Times(1)
	mockFile.EXPECT().MkdirAll(someDirPath, os.ModePerm).Return(ErrCouldNotCreateDirectory).Times(1)

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.EqualError(t, err, ErrCouldNotCreateDirectory.Error())
}
