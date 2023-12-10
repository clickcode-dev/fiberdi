package flog

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(log.DebugLevel)

	assert.NotEqual(t, logger, nil)
}
