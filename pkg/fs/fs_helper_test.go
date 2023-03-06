package fs

import (
	mockfile "github.com/amirvalhalla/fspool/mocks/file"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "test.txt")

	mockFile := mockfile.NewMockFileHelper(mockCtrl)
	mockFile.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	err := IsFileExists(someFilePath, mockFile.Stat)

	assert.Nil(t, err)
}

func TestIsFileExists_File_IsNotExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "test.txt")

	mockFile := mockfile.NewMockFileHelper(mockCtrl)
	mockFile.EXPECT().Stat(someFilePath).Return(nil, ErrFileIsNotExists).Times(1)

	err := IsFileExists(someFilePath, mockFile.Stat)

	assert.EqualError(t, err, ErrFileIsNotExists.Error())
}

func TestIsDirectoryExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someDirPath := filepath.Join("/test")
	mockFile := mockfile.NewMockFileHelper(mockCtrl)
	mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFile.EXPECT().Stat(someDirPath).Return(mockFileInfo, nil).Times(1)
	mockFile.EXPECT().IsNotExist(nil).Return(true).Times(1)

	err := IsDirectoryExists(someDirPath, mockFile.Stat, mockFile.IsNotExist)

	assert.Nil(t, err)
}

func TestIsDirectoryExists_DirectoryIsNotExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someDirPath := filepath.Join("/test")
	mockFile := mockfile.NewMockFileHelper(mockCtrl)

	mockFile.EXPECT().Stat(someDirPath).Return(nil, ErrDirectoryIsNotExists).Times(1)
	mockFile.EXPECT().IsNotExist(ErrDirectoryIsNotExists).Return(true).Times(1)

	err := IsDirectoryExists(someDirPath, mockFile.Stat, mockFile.IsNotExist)

	assert.EqualError(t, err, ErrDirectoryIsNotExists.Error())
}

func TestCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someDirPath := filepath.Join("/test")
	mockFile := mockfile.NewMockFileHelper(mockCtrl)
	mockFile.EXPECT().MkdirAll(someDirPath, os.ModePerm).Return(nil).Times(1)

	err := CreateDirectory(someDirPath, mockFile.MkdirAll)

	assert.Nil(t, err)
}

func TestCreateDirectory_CouldNotCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someDirPath := filepath.Join("/test")
	mockFile := mockfile.NewMockFileHelper(mockCtrl)
	mockFile.EXPECT().MkdirAll(someDirPath, os.ModePerm).Return(ErrCouldNotCreateDirectory).Times(1)

	err := CreateDirectory(someDirPath, mockFile.MkdirAll)

	assert.EqualError(t, err, ErrCouldNotCreateDirectory.Error())
}
