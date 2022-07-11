package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeStrToUl(t *testing.T) {
	t.Run("4294967295 base 2", func(t *testing.T) {
		var end string
		val, err := SafeStrToUl("11111111111111111111111111111111", &end, 2)
		assert.NoError(t, err)
		assert.Equal(t, uint(4294967295), val)
		assert.Equal(t, "", end)
	})

	t.Run("4294967295 base 10", func(t *testing.T) {
		var end string
		val, err := SafeStrToUl("4294967295", &end, 10)
		assert.NoError(t, err)
		assert.Equal(t, uint(4294967295), val)
		assert.Equal(t, "", end)
	})

	t.Run("4294967295 base 16", func(t *testing.T) {
		var end string
		val, err := SafeStrToUl("HHHHHHHH", &end, 16)
		assert.NoError(t, err)
		assert.Equal(t, uint(4294967295), val)
		assert.Equal(t, "", end)
	})

	t.Run("4294967295 base 52", func(t *testing.T) {
		var end string
		val, err := SafeStrToUl("CHPgSv", &end, 52)
		assert.NoError(t, err)
		assert.Equal(t, uint(4294967295), val)
		assert.Equal(t, "", end)
	})
}

func TestSafeUlToStr(t *testing.T) {
	t.Run("4294967295 base 2", func(t *testing.T) {
		val := SafeUlToStr(4294967295, 2, 5, '0')
		assert.Equal(t, "11111111111111111111111111111111", val)
	})

	t.Run("4294967295 base 10", func(t *testing.T) {
		val := SafeUlToStr(4294967295, 10, 5, '0')
		assert.Equal(t, "4294967295", val)
	})

	t.Run("4294967295 base 16", func(t *testing.T) {
		val := SafeUlToStr(4294967295, 16, 5, '0')
		assert.Equal(t, "HHHHHHHH", val)
	})

	t.Run("4294967295 base 52", func(t *testing.T) {
		val := SafeUlToStr(4294967295, 52, 5, '0')
		assert.Equal(t, "CHPgSv", val)
	})
}
