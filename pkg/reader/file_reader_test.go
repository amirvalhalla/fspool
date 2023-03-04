package reader

import (
	mockreader "github.com/amirvalhalla/fspool/mocks/reader"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	assert.NotNil(t, fReader)
}

func TestFileReader_ReadData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(0, nil).Times(1)

	_, err := fReader.ReadData(0, 0, io.SeekStart)

	assert.Nil(t, err)
}

func TestFileReader_ReadData_CouldNotSeek(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), ErrFileReaderCouldNotSeek).Times(1)

	_, err := fReader.ReadData(0, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFileReaderCouldNotSeek.Error())
}

func TestFileReader_ReadData_CouldNotRead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(0, ErrFileReaderCouldNotRead).Times(1)

	_, err := fReader.ReadData(0, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFileReaderCouldNotRead.Error())
}

func TestFileReader_ReadAllData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Stat().Return(mockFile, nil).Times(1)
	mockFile.EXPECT().Size().Return(int64(0)).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(0, nil).Times(1)

	_, err := fReader.ReadAllData()

	assert.Nil(t, err)
}

func TestFileReader_ReadAllData_CouldNotGetFileStat(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Stat().Return(nil, ErrFileReaderCouldNotGetFileStat).Times(1)

	_, err := fReader.ReadAllData()

	assert.EqualError(t, err, ErrFileReaderCouldNotGetFileStat.Error())
}

func TestFileReader_ReadAllData_CouldNotReadAllData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Stat().Return(mockFile, nil).Times(1)
	mockFile.EXPECT().Size().Return(int64(0)).Times(1)
	mockFile.EXPECT().Read([]byte{}).Return(0, ErrFileReaderCouldNotReadAllData).Times(1)

	_, err := fReader.ReadAllData()

	assert.EqualError(t, err, ErrFileReaderCouldNotReadAllData.Error())
}

func TestFileReader_Close(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Close().Return(nil).Times(1)

	err := fReader.Close()

	assert.Nil(t, err)
}

func TestFileReader_Close_CouldNotClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockreader.NewMockFile(mockCtrl)
	fReader := NewFileReader(mockFile)

	mockFile.EXPECT().Close().Return(ErrFileReaderCouldNotClose).Times(1)

	err := fReader.Close()

	assert.EqualError(t, err, ErrFileReaderCouldNotClose.Error())
}
