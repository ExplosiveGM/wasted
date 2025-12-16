package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	t.Run("Generates codes", func(t *testing.T) {
		result := GenerateCode(100000, 999999)
		assert.GreaterOrEqual(t, result, 100000, "код должен быть >= 100000")
		assert.LessOrEqual(t, result, 999999, "код должен быть <= 999999")
	})
}
