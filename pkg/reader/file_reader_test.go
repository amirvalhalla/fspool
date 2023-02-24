package reader

import (
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

var (
	basePath      string
	fileName      string
	validFilePath string
)

func TestNewFileReader(t *testing.T) {

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
			expectedErr: ErrCouldNotOpenFile,
		},
		{
			test:        "valid file path",
			path:        validFilePath,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			//create new file reader
			fReader, err := NewFileReader(tc.path)
			defer fReader.Close()

			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

	//delete created directory with files
}

func TestFileReader_ReadData(t *testing.T) {

	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)
	defer fReader.Close()

	//build our needed testcase struct
	type testCase struct {
		test         string
		writeDataLen int
		buffLen      int
		offset       int64
		expectedErr  error
	}

	//create testcase scenarios
	tc := testCase{
		test:         "valid read data",
		writeDataLen: 5,
		buffLen:      5,
		offset:       0,
		expectedErr:  nil,
	}

	//initialize required things before run actual
	randomBuff := make([]byte, tc.writeDataLen)
	rand.Read(randomBuff)
	os.WriteFile(validFilePath, randomBuff, os.ModePerm)

	//scenario
	_, err := fReader.ReadData(tc.offset, tc.buffLen, io.SeekCurrent)

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}

}

func TestFileReader_ReadData_CouldNotSeek(t *testing.T) {

	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)
	defer fReader.Close()

	//build our needed testcase struct
	type testCase struct {
		test         string
		writeDataLen int
		buffLen      int
		offset       int64
		expectedErr  error
	}

	//create testcase scenarios
	tc := testCase{
		test:         "seek over file len",
		writeDataLen: 5,
		buffLen:      5,
		offset:       -1,
		expectedErr:  ErrCouldNotSeek,
	}

	//initialize required things before run actual
	randomBuff := make([]byte, tc.writeDataLen)
	rand.Read(randomBuff)
	os.WriteFile(validFilePath, randomBuff, os.ModePerm)

	//scenario
	_, err := fReader.ReadData(tc.offset, tc.buffLen, io.SeekCurrent)

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}

}

func TestFileReader_ReadData_CouldNotReadData(t *testing.T) {

	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)
	defer fReader.Close()

	//build our needed testcase struct
	type testCase struct {
		test         string
		writeDataLen int
		buffLen      int
		offset       int64
		expectedErr  error
	}

	//create testcase scenarios
	tc := testCase{
		test:         "could not read data",
		writeDataLen: 5,
		buffLen:      10,
		offset:       10,
		expectedErr:  ErrCouldNotRead,
	}

	//initialize required things before run actual
	randomBuff := make([]byte, tc.writeDataLen)
	rand.Read(randomBuff)
	os.WriteFile(validFilePath, randomBuff, os.ModePerm)

	//scenario
	_, err := fReader.ReadData(tc.offset, tc.buffLen, io.SeekCurrent)

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}

}

func TestFileReader_ReadAllData(t *testing.T) {

	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)
	defer fReader.Close()

	//build our needed testcase struct
	type testCase struct {
		test         string
		writeDataLen int
		expectedErr  error
	}

	//create testcase scenarios
	tc := testCase{
		test:         "valid read all data",
		writeDataLen: 5,
		expectedErr:  nil,
	}

	//initialize required things before run actual
	randomBuff := make([]byte, tc.writeDataLen)
	rand.Read(randomBuff)
	os.WriteFile(validFilePath, randomBuff, os.ModePerm)

	//scenario
	_, err := fReader.ReadAllData()

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}
}

func TestFileReader_Close_CouldNotClose(t *testing.T) {
	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)

	//build our needed testcase struct
	type testCase struct {
		test        string
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "valid close file reader",
		expectedErr: nil,
	}

	//scenario
	err := fReader.Close()

	// Check if the error matches the expected error
	if err != tc.expectedErr {
		t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
	}
}

func TestFileReader_Close(t *testing.T) {
	initializeRequiredSpaceForTest()
	fReader, _ := NewFileReader(validFilePath)
	defer os.RemoveAll(basePath)

	//build our needed testcase struct
	type testCase struct {
		test        string
		expectedErr error
	}

	//create testcase scenarios
	tc := testCase{
		test:        "could not close file reader",
		expectedErr: ErrCouldNotClose,
	}

	fReader.Close()

	//scenario
	err := fReader.Close()

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

	log.Println("just for test yml")

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
