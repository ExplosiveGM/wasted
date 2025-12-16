package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEmail(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		result := IsEmail("some@email.com")
		assert.Equal(t, true, result)
	})
	t.Run("invalid email", func(t *testing.T) {
		result := IsEmail("some.email.com")
		assert.Equal(t, false, result)
	})
}

func TestIsPhoneNumber(t *testing.T) {
	t.Run("valid phone number", func(t *testing.T) {
		t.Run("number without formatting", func(t *testing.T) {
			result := IsPhoneNumber("+79161234567")
			assert.Equal(t, true, result)
		})
		t.Run("number with formatting", func(t *testing.T) {
			result := IsPhoneNumber("+7(916)-123-45-67")
			assert.Equal(t, true, result)
		})
		t.Run("number without plus", func(t *testing.T) {
			result := IsPhoneNumber("79161234567")
			assert.Equal(t, true, result)
		})
	})
}

func TestDetermineLoginType(t *testing.T) {
	t.Run("phone number passed", func(t *testing.T) {
		result, err := DetermineLoginType("+7(916)-123-45-67")
		require.NoError(t, err)
		assert.Equal(t, "phone", result)
	})
	t.Run("email passed", func(t *testing.T) {
		result, err := DetermineLoginType("some@email.com")
		require.NoError(t, err)
		assert.Equal(t, "email", result)
	})
	t.Run("invalid login passed", func(t *testing.T) {
		result, err := DetermineLoginType("some.email.com")
		require.ErrorContains(t, err, "unknown login type")
		assert.Equal(t, "", result)
	})
}
