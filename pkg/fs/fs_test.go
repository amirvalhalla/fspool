package fs

import (
	mockfile "github.com/amirvalhalla/fspool/mocks/file"
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

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.Nil(t, err)
}

func TestNewFilesystem_With_ROnly_Perm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.Perm = cfgs2.ROnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.Nil(t, err)
}

func TestNewFilesystem_With_WOnly_Perm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.Perm = cfgs2.WOnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.Nil(t, err)
}

func TestNewFilesystem_EmptyFilePath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem("", fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.EqualError(t, err, ErrFilesystemFilepathIsEmpty.Error())
}

func TestNewFilesystem_ConflictInMemoryRentSizeWithFlushSize(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.FlushSize = 60 * 1024 * 1024

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.EqualError(t, err, ErrFilesystemMemoryRentConflictWithFlushSize.Error())
}

func TestNewFilesystem_FilePathIsNotExists_With_ReadOnly_Permission(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, ErrFileIsNotExists).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.ROnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.EqualError(t, err, ErrFileIsNotExists.Error())
}

func TestNewFilesystem_FilePathIsNotExists_CouldNotCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someDirPath := filepath.Join("/test")
	someFilePath := filepath.Join(someDirPath, "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	//mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, ErrDirectoryIsNotExists).Times(2)
	mockFileHelper.EXPECT().IsNotExist(ErrDirectoryIsNotExists).Return(true).Times(1)
	mockFileHelper.EXPECT().MkdirAll(someDirPath, os.ModePerm).Return(ErrCouldNotCreateDirectory).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	assert.EqualError(t, err, ErrCouldNotCreateDirectory.Error())
}
