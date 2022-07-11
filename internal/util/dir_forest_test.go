package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirForest(t *testing.T) {
	t.Run("length 1 path, depth 1", func(t *testing.T) {
		path := DirForest("a", 1)
		expected := filepath.Join("a") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})

	t.Run("length 5 path, depth 5", func(t *testing.T) {
		path := DirForest("abcde", 5)
		expected := filepath.Join("a", "b", "c", "d", "e") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})

	t.Run("length 4 path, depth 1", func(t *testing.T) {
		path := DirForest("abcd", 1)
		expected := filepath.Join("a") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})

	t.Run("length 4 path, depth 2", func(t *testing.T) {
		path := DirForest("abcd", 2)
		expected := filepath.Join("a", "b") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})

	t.Run("length 1 path, depth 2", func(t *testing.T) {
		path := DirForest("a", 2)
		expected := filepath.Join("a", "_") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})

	t.Run("length 1 path, depth 4", func(t *testing.T) {
		path := DirForest("a", 4)
		expected := filepath.Join("a", "_", "_", "_") + string(os.PathSeparator)
		assert.Equal(t, expected, path)
	})
}
