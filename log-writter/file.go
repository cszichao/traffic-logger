package log_writter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NewFileWriter make a log writer to save traffic into file
func NewFileWriter(logDir, name string) (*os.File, error) {
	timeNOW := func() string {
		return time.Now().Format("2006-01-02-15.04.05.999999999")
	}

	logFileName := name + ".log"
	logFilePath := filepath.Join(logDir, logFileName)

	if fileExists(logFilePath) {
		_ = os.Rename(logFilePath, filepath.Join(logDir, name+"."+timeNOW()+".log"))
	}
	return makefile(logDir, logFileName)
}

// log file helper
func makefile(dir string, filename string) (f *os.File, e error) {
	if err := createDirectories(dir); err != nil {
		return nil, err
	}
	filePath := filepath.Join(dir, filename)
	fileWriter, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return fileWriter, nil
}

func createDirectories(dir string) error {
	if fi, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("%v (checking directory)", err)
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("%v (creating directory)", err)
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("destination path is not directory")
	}
	return nil
}

func fileExists(path string) bool {
	if len(path) == 0 {
		return false
	}
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if fi.IsDir() {
		return false
	}
	return true
}
