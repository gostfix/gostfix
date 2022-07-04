package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMacParse(t *testing.T) {
	MacParse("hello", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_LITERAL, typ)
		assert.Equal(t, "hello", s)
		return 0
	}, nil)

	MacParse("$hello", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_EXPR, typ)
		assert.Equal(t, "hello", s)
		return 0
	}, nil)

	MacParse("$(hello)", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_EXPR, typ)
		assert.Equal(t, "hello", s)
		return 0
	}, nil)

	MacParse("${hello}", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_EXPR, typ)
		assert.Equal(t, "hello", s)
		return 0
	}, nil)

	MacParse("${${hello}}", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_EXPR, typ)
		assert.Equal(t, "${hello}", s)
		return 0
	}, nil)

	MacParse("$$hello", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_LITERAL, typ)
		assert.Equal(t, "$hello", s)
		return 0
	}, nil)

	MacParse("$", func(typ int, s string, x interface{}) int {
		return 0
	}, nil)

	MacParse("$$", func(typ int, s string, x interface{}) int {
		assert.Equal(t, MAC_PARSE_LITERAL, typ)
		assert.Equal(t, "$", s)
		return 0
	}, nil)

	MacParse("$$$", func(typ int, s string, x interface{}) int {
		return 0
	}, nil)
}
