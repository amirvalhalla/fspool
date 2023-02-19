package reader

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	basePath      = "/test"
	fileName      = "testing.txt"
	validFilePath = filepath.Join(basePath, fileName)
)

func TestNewFileReader(t *testing.T) {
	//build file for testcase scenario , it automatically will delete after test ran
	os.Mkdir(basePath, os.ModePerm)
	f, _ := os.Create(validFilePath)

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
			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			fReader.Close()
		})
	}

	//delete created directory with files
	f.Close()
	os.RemoveAll(basePath)
}
