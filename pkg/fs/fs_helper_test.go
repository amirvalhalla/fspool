package fs

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsFileExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath("/test/test.txt").Return(nil, nil).Times(1)

	err := IsFileExists("/test/test.txt", mockFile)

	assert.Nil(t, err)
}

func TestIsFileExists_File_IsNotExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath("/test/test.txt").Return(nil, ErrFileIsNotExists).Times(1)

	err := IsFileExists("/test/test.txt", mockFile)

	assert.EqualError(t, err, ErrFileIsNotExists.Error())
}

func TestIsDirectoryExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath("/test").Return(nil, nil).Times(1)
	mockFile.EXPECT().IsNotExist(nil).Return(false).Times(1)

	err := IsDirectoryExists("/test", mockFile)

	assert.Nil(t, err)
}

func TestIsDirectoryExists_DirectoryIsNotExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().StatWithFilePath("/test").Return(nil, ErrDirectoryIsNotExists).Times(1)
	mockFile.EXPECT().IsNotExist(ErrDirectoryIsNotExists).Return(true).Times(1)

	err := IsDirectoryExists("/test", mockFile)

	assert.EqualError(t, err, ErrDirectoryIsNotExists.Error())
}

func TestCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().MkdirAll("/test", os.ModePerm).Return(nil).Times(1)

	err := CreateDirectory("/test", mockFile)

	assert.Nil(t, err)
}

func TestCreateDirectory_CouldNotCreateDirectory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := NewMockFile(mockCtrl)
	mockFile.EXPECT().MkdirAll("/test", os.ModePerm).Return(ErrCouldNotCreateDirectory).Times(1)

	err := CreateDirectory("/test", mockFile)

	assert.EqualError(t, err, ErrCouldNotCreateDirectory.Error())
}
