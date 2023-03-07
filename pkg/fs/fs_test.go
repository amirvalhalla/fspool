package fs

import (
	mockfile "github.com/amirvalhalla/fspool/mocks/file"
	cfgs2 "github.com/amirvalhalla/fspool/pkg/cfgs"
	cfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
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

func TestFilesystem_Write(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Write([]byte{2}).Return(int(0), nil).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.Write([]byte{2}, 0, io.SeekStart)

	assert.Nil(t, err)
}

func TestFilesystem_Write_Writer_nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.ROnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.Write([]byte{2}, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFilesystemWriterHasBeenClosed.Error())
}

func TestFilesystem_Write_Writer_CouldNotWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Write([]byte{2}).Return(int(0), ErrFilesystemCouldNotWrite).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.Write([]byte{2}, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFilesystemCouldNotWrite.Error())
}

func TestFilesystem_CloseWriter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Close().Return(nil).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseWriter()

	assert.Nil(t, err)
}

func TestFilesystem_CloseWriter_nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.ROnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseWriter()

	assert.EqualError(t, err, ErrFilesystemWriterHasBeenClosed.Error())
}

func TestFilesystem_CloseWriter_CouldNotClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Close().Return(ErrFilesystemCouldNotCloseWriter).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseWriter()

	assert.EqualError(t, err, ErrFilesystemCouldNotCloseWriter.Error())
}

func TestFilesystem_ReadData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(int(0), nil).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.ReadData(0, 0, io.SeekStart)

	assert.Nil(t, err)
}

func TestFilesystem_ReadData_Reader_nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.WOnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.ReadData(0, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFilesystemReaderEmpty.Error())
}

func TestFilesystem_ReadData_Reader_Occupying(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil).AnyTimes()
	mockFile.EXPECT().Read([]byte{}).Return(int(0), nil).AnyTimes()

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	go func() {
		for {
			f.ReadData(0, 0, io.SeekStart)
		}
	}()

	var err error

	for {
		_, err = f.ReadData(0, 0, io.SeekStart)
		if err != nil {
			break
		}
	}

	assert.EqualError(t, err, ErrFilesystemReaderOccupying.Error())
}

func TestFilesystem_ReadData_CouldNotReadData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(int(0), ErrFilesystemCouldNotReadData).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.ReadData(0, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFilesystemCouldNotReadData.Error())
}

func TestFilesystem_ReadAllData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFile.EXPECT().Read([]byte{}).Return(int(0), nil).Times(1)
	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).Times(1)

	mockFileInfo.EXPECT().Size().Return(int64(0)).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.ReadAllData()

	assert.Nil(t, err)
}

func TestFilesystem_ReadAllData_Reader_Nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.WOnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)
	_, err := f.ReadAllData()

	assert.EqualError(t, err, ErrFilesystemReaderEmpty.Error())
}

func TestFilesystem_ReadAllData_Reader_Occupying(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFile.EXPECT().Read([]byte{}).Return(int(0), nil).AnyTimes()
	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).AnyTimes()

	mockFileInfo.EXPECT().Size().Return(int64(0)).AnyTimes()

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	var err error
	go func() {
		for {
			f.ReadAllData()
		}
	}()

	for {
		_, err = f.ReadAllData()
		if err != nil {
			break
		}
	}

	assert.EqualError(t, err, ErrFilesystemReaderOccupying.Error())
}

func TestFilesystem_ReadAllData_CouldNotReadAllData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFile.EXPECT().Read([]byte{}).Return(int(0), ErrFilesystemCouldNotReadAllData).Times(1)
	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).Times(1)

	mockFileInfo.EXPECT().Size().Return(int64(0)).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.ReadAllData()

	assert.EqualError(t, err, ErrFilesystemCouldNotReadAllData.Error())
}

func TestFilesystem_CloseReader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Close().Return(nil).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseReader()

	assert.Nil(t, err)
}

func TestFilesystem_CloseReader_Nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.WOnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseReader()

	assert.EqualError(t, err, ErrFilesystemReaderEmpty.Error())
}

func TestFilesystem_CloseReader_Occupying(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFileInfo := mockfile.NewMockFileInfo(mockCtrl)

	mockFile.EXPECT().Close().Return(nil).AnyTimes()
	mockFile.EXPECT().Read([]byte{}).Return(int(0), nil).AnyTimes()
	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).AnyTimes()

	mockFileInfo.EXPECT().Size().Return(int64(0)).AnyTimes()

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	var err error
	go func() {
		for {
			f.ReadAllData()
		}
	}()

	for {
		err = f.CloseReader()
		if err != nil {
			break
		}
	}

	assert.EqualError(t, err, ErrFilesystemReaderOccupying.Error())
}

func TestFilesystem_CloseReader_CouldNotClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)
	mockFile.EXPECT().Close().Return(ErrFilesystemCouldNotCloseReader).Times(1)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	err := f.CloseReader()

	assert.EqualError(t, err, ErrFilesystemCouldNotCloseReader.Error())
}

func TestFilesystem_GetReaderState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.GetReaderState()

	assert.Nil(t, err)
}

func TestFilesystem_GetReaderState_Nil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	someFilePath := filepath.Join("/test", "/test.txt")

	mockFile := mockfile.NewMockFile(mockCtrl)

	mockFileHelper := mockfile.NewMockFileHelper(mockCtrl)
	mockFileHelper.EXPECT().Stat(someFilePath).Return(nil, nil).Times(1)

	fsConfig := cfgs.FSConfiguration{}
	fsConfig.New()
	fsConfig.Perm = cfgs2.WOnly

	f, _ := NewFilesystem(someFilePath, fsConfig, mockFile, mockFileHelper.Stat, mockFileHelper.IsNotExist, mockFileHelper.MkdirAll)

	_, err := f.GetReaderState()

	assert.EqualError(t, err, ErrFilesystemReaderEmpty.Error())
}
