package fiberdi

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateClientWithoutConfig(t *testing.T) {
	app := New(&Module{})

	assert.Equal(t, true, app.Config().ReduceMemoryUsage)
	assert.Equal(t, true, app.Config().StrictRouting)
	assert.Equal(t, true, app.Config().CaseSensitive)
	assert.Equal(t, true, app.Config().DisableStartupMessage)
}

func TestCreateClientWithConfig(t *testing.T) {
	app := New(&Module{}, fiber.Config{
		AppName:      "test-1234",
		ServerHeader: "test",
	})

	assert.Equal(t, true, app.Config().ReduceMemoryUsage)
	assert.Equal(t, true, app.Config().StrictRouting)
	assert.Equal(t, true, app.Config().CaseSensitive)
	assert.Equal(t, true, app.Config().DisableStartupMessage)

	assert.Equal(t, "test-1234", app.Config().AppName)
	assert.Equal(t, "test", app.Config().ServerHeader)
}
