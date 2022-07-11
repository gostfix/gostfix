package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyStrTok(t *testing.T) {
	t.Run("empty sep, empty src", func(t *testing.T) {
		tokens := MyStrTok("", "")
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("empty sep, non-empty src", func(t *testing.T) {
		tokens := MyStrTok("abcd", "")
		assert.Equal(t, []string{"abcd"}, tokens)
	})

	t.Run("empty src", func(t *testing.T) {
		tokens := MyStrTok("", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single space src", func(t *testing.T) {
		tokens := MyStrTok(" ", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single tab src", func(t *testing.T) {
		tokens := MyStrTok("\t", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single carriage return src", func(t *testing.T) {
		tokens := MyStrTok("\r", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single newline src", func(t *testing.T) {
		tokens := MyStrTok("\n", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("multiple spaces src", func(t *testing.T) {
		tokens := MyStrTok("       ", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("mixed spaces src", func(t *testing.T) {
		tokens := MyStrTok(" \t \r \n \t\r\n ", CHARS_SPACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("one src, one rune", func(t *testing.T) {
		tokens := MyStrTok("a", CHARS_SPACE)
		assert.Equal(t, []string{"a"}, tokens)
	})

	t.Run("one src, multiple runes", func(t *testing.T) {
		tokens := MyStrTok("abc", CHARS_SPACE)
		assert.Equal(t, []string{"abc"}, tokens)
	})

	t.Run("multiple src, mixed runes, only spaces", func(t *testing.T) {
		tokens := MyStrTok("a bc def ghijk", CHARS_SPACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})

	t.Run("multiple src, mixed spaces", func(t *testing.T) {
		tokens := MyStrTok("\t  a \r\n\tbc     def \tghijk\r\n", CHARS_SPACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})

	t.Run("multiple src, mixed spaces", func(t *testing.T) {
		tokens := MyStrTok("\t  a \r\n\tbc     def \tghijk\r\n", CHARS_SPACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})
}

func TestMyStrTokQ(t *testing.T) {
	t.Run("empty sep, empty src", func(t *testing.T) {
		tokens := MyStrTokQ("", "", CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("empty sep, non-empty src", func(t *testing.T) {
		tokens := MyStrTokQ("abcd", "", CHARS_BRACE)
		assert.Equal(t, []string{"abcd"}, tokens)
	})

	t.Run("empty src", func(t *testing.T) {
		tokens := MyStrTokQ("", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single space src", func(t *testing.T) {
		tokens := MyStrTokQ(" ", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single tab src", func(t *testing.T) {
		tokens := MyStrTokQ("\t", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single carriage return src", func(t *testing.T) {
		tokens := MyStrTokQ("\r", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("single newline src", func(t *testing.T) {
		tokens := MyStrTokQ("\n", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("multiple spaces src", func(t *testing.T) {
		tokens := MyStrTokQ("       ", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("mixed spaces src", func(t *testing.T) {
		tokens := MyStrTokQ(" \t \r \n \t\r\n ", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{}, tokens)
	})

	t.Run("one src, one rune", func(t *testing.T) {
		tokens := MyStrTokQ("a", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"a"}, tokens)
	})

	t.Run("one src, multiple runes", func(t *testing.T) {
		tokens := MyStrTokQ("abc", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"abc"}, tokens)
	})

	t.Run("multiple src, mixed runes, only spaces", func(t *testing.T) {
		tokens := MyStrTokQ("a bc def ghijk", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})

	t.Run("multiple src, mixed spaces", func(t *testing.T) {
		tokens := MyStrTokQ("\t  a \r\n\tbc     def \tghijk\r\n", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})

	t.Run("multiple src, mixed spaces", func(t *testing.T) {
		tokens := MyStrTokQ("\t  a \r\n\tbc     def \tghijk\r\n", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"a", "bc", "def", "ghijk"}, tokens)
	})

	t.Run("two src", func(t *testing.T) {
		tokens := MyStrTokQ("foo bar", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"foo", "bar"}, tokens)
	})

	t.Run("single quote", func(t *testing.T) {
		tokens := MyStrTokQ("{ bar }  ", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"{ bar }"}, tokens)
	})

	t.Run("multiple tokens, single quote", func(t *testing.T) {
		tokens := MyStrTokQ("foo { bar } baz", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"foo", "{ bar }", "baz"}, tokens)
	})

	t.Run("multiple tokens, single quote attached left", func(t *testing.T) {
		tokens := MyStrTokQ("foo{ bar } baz", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"foo{ bar }", "baz"}, tokens)
	})

	t.Run("multiple tokens, single quote attached right", func(t *testing.T) {
		tokens := MyStrTokQ("foo { bar }baz", CHARS_SPACE, CHARS_BRACE)
		assert.Equal(t, []string{"foo", "{ bar }baz"}, tokens)
	})
}
