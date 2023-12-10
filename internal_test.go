package fiberdi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterArray(t *testing.T) {
	arr := []string{"applepen", "pineapple", "banana"}
	filteredArr := filter(arr, func(fruit string) bool {
		return strings.ContainsAny(fruit, "apple")
	})

	assert.Equal(t, 2, len(filteredArr))
}

func TestTernaryIf(t *testing.T) {
	result := ternary(1 == 1, true, false)

	assert.Equal(t, true, result)
}

func TestTernaryElse(t *testing.T) {
	result := ternary(1 != 1, true, false)

	assert.Equal(t, false, result)
}
