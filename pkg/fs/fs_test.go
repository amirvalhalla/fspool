package fs

import (
	mockfile "github.com/amirvalhalla/fspool/mocks/file"
	cfgs2 "github.com/amirvalhalla/fspool/pkg/cfgs"
	cfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewFilesystem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.Nil(t, err)
}

func TestNewFilesystem_With_ROnly_Perm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

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

	mockFile := mockfile.NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.Perm = cfgs2.WOnly

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.Nil(t, err)
}

func TestNewFilesystem_EmptyFilePath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	_, err := NewFilesystem("", fsConfig, mockFile)

	assert.EqualError(t, err, ErrFilesystemFilepathIsEmpty.Error())
}

func TestNewFilesystem_ConflictInMemoryRentSizeWithFlushSize(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	fsConfig.FlushSize = 60 * 1024 * 1024

	_, err := NewFilesystem(someFilePath, fsConfig, mockFile)

	assert.EqualError(t, err, ErrFilesystemMemoryRentConflictWithFlushSize.Error())
}
