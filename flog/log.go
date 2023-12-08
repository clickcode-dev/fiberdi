package flog

import (
	"os"

	"github.com/charmbracelet/log"
)

func NewLogger(level log.Level) *log.Logger {
	return log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		Level:           level,
	})
}
