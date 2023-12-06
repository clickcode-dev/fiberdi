package fiberdi

import (
	"os"

	"github.com/charmbracelet/log"
)

func newLogger() *log.Logger {
	return log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})
}
