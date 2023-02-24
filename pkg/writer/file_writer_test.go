package writer

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

var (
	basePath      string
	fileName      string
	validFilePath string
)

func TestNewFileWriter(t *testing.T) {
	initializeRequiredSpaceForTest()
	defer os.RemoveAll(basePath)

	//build our needed testcase struct
	type testCase struct {
		test        string
		path        string
		expectedErr error
	}

	//create testcase scenarios
	testCases := []testCase{
		{
			test:        "empty file path",
			path:        "",
			expectedErr: ErrFileWriterCouldNotOpenFile,
		},
		{
			test:        "valid file path",
			path:        validFilePath,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			//create new file writer
			fWriter, err := NewFileWriter(tc.path)
			defer fWriter.Close()

			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

}

func TestFileWriter_AddOrUpdateData(t *testing.T) {
	initializeRequiredSpaceForTest()
	fWriter, _ := NewFileWriter(validFilePath)
	defer os.RemoveAll(basePath)
	defer fWriter.Close()

	//build our needed testcase struct
	type testCase struct {
		test        string
		data        []byte
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "valid write data",
		data:        []byte{14, 58, 74, 21, 36},
		expectedErr: nil,
	}

	//scenario
	err := fWriter.AddOrUpdateData(tc.data, 0, io.SeekStart)

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}

}

func TestFileWriter_AddOrUpdateData_CouldNotSeek(t *testing.T) {
	initializeRequiredSpaceForTest()
	fWriter, _ := NewFileWriter(validFilePath)
	defer os.RemoveAll(basePath)
	defer fWriter.Close()

	//build our needed testcase struct
	type testCase struct {
		test        string
		data        []byte
		offset      int64
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "could not seek",
		data:        []byte{14, 58, 74, 21, 36},
		offset:      -1,
		expectedErr: ErrFileWriterCouldNotSeek,
	}

	//scenario
	err := fWriter.AddOrUpdateData(tc.data, tc.offset, io.SeekStart)

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}

}

func TestFileWriter_Close(t *testing.T) {
	initializeRequiredSpaceForTest()
	fWriter, _ := NewFileWriter(validFilePath)
	defer os.RemoveAll(basePath)

	//build our needed testcase struct
	type testCase struct {
		test        string
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "valid close file writer",
		expectedErr: nil,
	}

	//scenario
	err := fWriter.Close()

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}
}

func TestFileWriter_Close_CouldNotClose(t *testing.T) {
	initializeRequiredSpaceForTest()
	fWriter, _ := NewFileWriter(validFilePath)
	defer os.RemoveAll(basePath)

	//build our needed testcase struct
	type testCase struct {
		test        string
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "could not close file writer",
		expectedErr: ErrFileWriterCouldNotClose,
	}

	fWriter.Close()

	//scenario
	err := fWriter.Close()

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}
}

func initializeRequiredSpaceForTest() {

	//get current dir
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//initialize required variables for test
	basePath = filepath.Join(dir, "testDir")
	fileName = "test.txt"
	validFilePath = filepath.Join(basePath, fileName)

	//build file and directory for testcase scenario , it automatically will delete after test ran
	os.Mkdir(basePath, os.ModePerm)
	if f, err := os.Create(validFilePath); err == nil {
		f.Close()
	}
}
