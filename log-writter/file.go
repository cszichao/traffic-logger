package log_writter

import (
	"github.com/haxii/log/v2"
	"os"
)

// NewFileWriter make a log writer to save traffic into file
func NewFileWriter(logDir, name string) (*os.File, error) {
	return log.OpenLogFile(logDir, name)
}
