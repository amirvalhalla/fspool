package writer

import (
	mockfile "github.com/amirvalhalla/fspool/mocks/file"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewFileWriter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	assert.NotNil(t, fWriter)
}

func TestFileWriter_Write(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Write([]byte{}).Return(0, nil).Times(1)

	err := fWriter.Write([]byte{}, 0, io.SeekStart)

	assert.Nil(t, err)
}

func TestFileWriter_Write_CouldNotSeek(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), ErrFileWriterCouldNotSeek).Times(1)

	err := fWriter.Write([]byte{}, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFileWriterCouldNotSeek.Error())
}

func TestFileWriter_Write_CouldNotWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	mockFile.EXPECT().Seek(int64(0), 0).Return(int64(0), nil).Times(1)
	mockFile.EXPECT().Write([]byte{}).Return(0, ErrFileWriterCouldNotWrite).Times(1)

	err := fWriter.Write([]byte{}, 0, io.SeekStart)

	assert.EqualError(t, err, ErrFileWriterCouldNotWrite.Error())
}

func TestFileWriter_Close(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	mockFile.EXPECT().Close().Return(nil).Times(1)

	err := fWriter.Close()

	assert.Nil(t, err)
}

func TestFileWriter_Close_CouldNotClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFile := mockfile.NewMockFile(mockCtrl)
	fWriter := NewFileWriter(mockFile)

	mockFile.EXPECT().Close().Return(ErrFileWriterCouldNotClose).Times(1)

	err := fWriter.Close()

	assert.EqualError(t, err, ErrFileWriterCouldNotClose.Error())
}
