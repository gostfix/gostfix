package global

import (
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gostfix/gostfix/internal/util"
	"github.com/stretchr/testify/assert"
)

const alphaNum string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func randString(length int) string {
	rand.Seed(time.Now().UnixMicro())
	b := strings.Builder{}
	b.Grow(length)
	for i := 0; i < length; i++ {
		b.WriteByte(alphaNum[rand.Intn(len(alphaNum))])
	}
	return b.String()
}

func TestMailQueueIdIsOk(t *testing.T) {
	t.Run("empty queue id", func(t *testing.T) {
		b := MailQueueIdOk("")
		assert.False(t, b)
	})

	t.Run("max length + 1 queue id", func(t *testing.T) {
		b := MailQueueIdOk(randString(util.VALID_HOSTNAME_LEN + 1))
		assert.False(t, b)
	})

	t.Run("length 1 queue id", func(t *testing.T) {
		b := MailQueueIdOk("a")
		assert.True(t, b)
	})

	t.Run("max length queue id", func(t *testing.T) {
		b := MailQueueIdOk(randString(util.VALID_HOSTNAME_LEN))
		assert.True(t, b)
	})

	t.Run("valid queue id", func(t *testing.T) {
		b := MailQueueIdOk("abcd")
		assert.True(t, b)
	})

	t.Run("all characters valid queue id", func(t *testing.T) {
		b := MailQueueIdOk(alphaNum)
		assert.True(t, b)
	})

	// Check invalid characters on the boundaries of valid characters
	t.Run("invalid queue id, start of string /", func(t *testing.T) {
		b := MailQueueIdOk("/abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string /", func(t *testing.T) {
		b := MailQueueIdOk("abcd/")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string /", func(t *testing.T) {
		b := MailQueueIdOk("ab/cd")
		assert.False(t, b)
	})

	t.Run("invalid queue id, start of string :", func(t *testing.T) {
		b := MailQueueIdOk(":abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string :", func(t *testing.T) {
		b := MailQueueIdOk("abcd:")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string :", func(t *testing.T) {
		b := MailQueueIdOk("ab:cd")
		assert.False(t, b)
	})

	t.Run("invalid queue id, start of string @", func(t *testing.T) {
		b := MailQueueIdOk("@abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string @", func(t *testing.T) {
		b := MailQueueIdOk("abcd@")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string @", func(t *testing.T) {
		b := MailQueueIdOk("ab@cd")
		assert.False(t, b)
	})

	t.Run("invalid queue id, start of string [", func(t *testing.T) {
		b := MailQueueIdOk("[abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string [", func(t *testing.T) {
		b := MailQueueIdOk("abcd[")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string [", func(t *testing.T) {
		b := MailQueueIdOk("ab[cd")
		assert.False(t, b)
	})

	t.Run("invalid queue id, start of string `", func(t *testing.T) {
		b := MailQueueIdOk("`abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string `", func(t *testing.T) {
		b := MailQueueIdOk("abcd`")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string `", func(t *testing.T) {
		b := MailQueueIdOk("ab`cd")
		assert.False(t, b)
	})

	t.Run("invalid queue id, start of string {", func(t *testing.T) {
		b := MailQueueIdOk("{abcd")
		assert.False(t, b)
	})
	t.Run("invalid queue id, end of string {", func(t *testing.T) {
		b := MailQueueIdOk("abcd{")
		assert.False(t, b)
	})
	t.Run("invalid queue id, middle of string {", func(t *testing.T) {
		b := MailQueueIdOk("ab{cd")
		assert.False(t, b)
	})
}

// Long queue examples: /storage_array/mta/instances/mail180-7.suw31.mandrillapp.com/
//
// spool/deferred/0/4LhRXm0N74zHXYVZ1
// spool/deferred/1/4Lh7QT0ckvzHZBtdV
// spool/deferred/2/4LhY8l1KjxzHXYWlN
// spool/deferred/3/4LhYLp1SZWzHXYWlP
// spool/deferred/4/4LhBbb28Y9zHXYTmD
// spool/deferred/5/4Lgg1T2TjrzHXYR23
// spool/deferred/6/4LhN1Y2ntPzHXYV9K
// spool/deferred/7/4Lh9HG3XXpzHXYTMv
// spool/deferred/8/4LfzbW46tCzHXYR2M
// spool/deferred/9/4LftnJ4bz8zHXYR8x
// spool/deferred/A/4LhLzV4pLhzHZBtfC
// spool/deferred/B/4LhHTW5VCSzHXYTm0
// spool/deferred/C/4LgBlf5krJzHXYSb0
// spool/deferred/D/4Lg7hB6QFMzHXYQxL
// spool/deferred/E/4LhFXL6tkrzHXYVy0
// spool/deferred/F/4LhGrz75WTzHXYQyV

func TestMailQueueDirPanicQueueName(t *testing.T) {
	assert.Panics(t, func() {
		MailQueueDir("", "4LhFXL6tkrzHXYVy0")
	})

	assert.Panics(t, func() {
		MailQueueDir("ディレクトリ", "4LhFXL6tkrzHXYVy0")
	})

	assert.Panics(t, func() {
		MailQueueDir("bad directory", "4LhFXL6tkrzHXYVy0")
	})

	assert.Panics(t, func() {
		MailQueueDir(strings.Repeat("a", 101), "4LhFXL6tkrzHXYVy0")
	})

	assert.NotPanics(t, func() {
		MailQueueDir("a", "4LhFXL6tkrzHXYVy0")
	})

	assert.NotPanics(t, func() {
		MailQueueDir(strings.Repeat("a", 100), "4LhFXL6tkrzHXYVy0")
	})
}

func TestMailQueueDirPanicQueueId(t *testing.T) {
	assert.Panics(t, func() {
		MailQueueDir("deferred", "")
	})

	assert.Panics(t, func() {
		MailQueueDir("deferred", strings.Repeat("a", util.VALID_HOSTNAME_LEN+1))
	})

	assert.Panics(t, func() {
		MailQueueDir("deferred", "bad queue_id")
	})

	assert.Panics(t, func() {
		MailQueueDir("deferred", "id_列")
	})

	assert.NotPanics(t, func() {
		MailQueueDir("deferred", "a")
	})

	assert.NotPanics(t, func() {
		MailQueueDir("deferred", strings.Repeat("a", util.VALID_HOSTNAME_LEN))
	})
}

func TestMailQueueDirUnhashed(t *testing.T) {
	// The default hashed queue names is a reasonable test
	VarHashQueueNames = DEF_HASH_QUEUE_NAMES
	VarHashQueueDepth = 1

	t.Run("no hash directory path", func(t *testing.T) {
		queue_dir := MailQueueDir("active", "4LhFXL6tkrzHXYVy0")
		assert.Equal(t, "active", queue_dir)
	})

}

func TestMailQueueDirDepthOne(t *testing.T) {
	// The default hashed queue names is a reasonable test
	VarHashQueueNames = DEF_HASH_QUEUE_NAMES
	VarHashQueueDepth = 1

	testCases := []struct {
		Dir      string
		QID      string
		Expected string
	}{
		{"deferred", "4LhRXm0N74zHXYVZ1", filepath.Join("deferred", "0")},
		{"deferred", "4Lh7QT0ckvzHZBtdV", filepath.Join("deferred", "1")},
		{"deferred", "4LhY8l1KjxzHXYWlN", filepath.Join("deferred", "2")},
		{"deferred", "4LhYLp1SZWzHXYWlP", filepath.Join("deferred", "3")},
		{"deferred", "4LhBbb28Y9zHXYTmD", filepath.Join("deferred", "4")},
		{"deferred", "4Lgg1T2TjrzHXYR23", filepath.Join("deferred", "5")},
		{"deferred", "4LhN1Y2ntPzHXYV9K", filepath.Join("deferred", "6")},
		{"deferred", "4Lh9HG3XXpzHXYTMv", filepath.Join("deferred", "7")},
		{"deferred", "4LfzbW46tCzHXYR2M", filepath.Join("deferred", "8")},
		{"deferred", "4LftnJ4bz8zHXYR8x", filepath.Join("deferred", "9")},
		{"deferred", "4LhLzV4pLhzHZBtfC", filepath.Join("deferred", "A")},
		{"deferred", "4LhHTW5VCSzHXYTm0", filepath.Join("deferred", "B")},
		{"deferred", "4LgBlf5krJzHXYSb0", filepath.Join("deferred", "C")},
		{"deferred", "4Lg7hB6QFMzHXYQxL", filepath.Join("deferred", "D")},
		{"deferred", "4LhFXL6tkrzHXYVy0", filepath.Join("deferred", "E")},
		{"deferred", "4LhGrz75WTzHXYQyV", filepath.Join("deferred", "F")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Expected, func(t *testing.T) {
			queue_dir := MailQueueDir(testCase.Dir, testCase.QID)
			assert.Equal(t, testCase.Expected, queue_dir)
		})
	}
}

func TestMailQueueDirDepthTwo(t *testing.T) {
	// The default hashed queue names is a reasonable test
	VarHashQueueNames = DEF_HASH_QUEUE_NAMES
	VarHashQueueDepth = 2

	testCases := []struct {
		Dir      string
		QID      string
		Expected string
	}{
		{"deferred", "4LhRXm0N74zHXYVZ1", filepath.Join("deferred", "0", "D")},
		{"deferred", "4Lh7QT0ckvzHZBtdV", filepath.Join("deferred", "1", "5")},
		{"deferred", "4LhY8l1KjxzHXYWlN", filepath.Join("deferred", "2", "E")},
		{"deferred", "4LhYLp1SZWzHXYWlP", filepath.Join("deferred", "3", "2")},
		{"deferred", "4LhBbb28Y9zHXYTmD", filepath.Join("deferred", "4", "A")},
		{"deferred", "4Lgg1T2TjrzHXYR23", filepath.Join("deferred", "5", "5")},
		{"deferred", "4LhN1Y2ntPzHXYV9K", filepath.Join("deferred", "6", "0")},
		{"deferred", "4Lh9HG3XXpzHXYTMv", filepath.Join("deferred", "7", "9")},
		{"deferred", "4LfzbW46tCzHXYR2M", filepath.Join("deferred", "8", "D")},
		{"deferred", "4LftnJ4bz8zHXYR8x", filepath.Join("deferred", "9", "E")},
		{"deferred", "4LhLzV4pLhzHZBtfC", filepath.Join("deferred", "A", "5")},
		{"deferred", "4LhHTW5VCSzHXYTm0", filepath.Join("deferred", "B", "C")},
		{"deferred", "4LgBlf5krJzHXYSb0", filepath.Join("deferred", "C", "5")},
		{"deferred", "4Lg7hB6QFMzHXYQxL", filepath.Join("deferred", "D", "C")},
		{"deferred", "4LhFXL6tkrzHXYVy0", filepath.Join("deferred", "E", "C")},
		{"deferred", "4LhGrz75WTzHXYQyV", filepath.Join("deferred", "F", "3")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Expected, func(t *testing.T) {
			queue_dir := MailQueueDir(testCase.Dir, testCase.QID)
			assert.Equal(t, testCase.Expected, queue_dir)
		})
	}
}

func TestMailQueuePath(t *testing.T) {
	VarHashQueueNames = DEF_HASH_QUEUE_NAMES
	VarHashQueueDepth = 1

	t.Run("no hash directory path", func(t *testing.T) {
		queue_path := MailQueuePath("active", "4LhFXL6tkrzHXYVy0")
		assert.Equal(t, filepath.Join("active", "4LhFXL6tkrzHXYVy0"), queue_path)
	})

	t.Run("hash directory path", func(t *testing.T) {
		queue_path := MailQueuePath("deferred", "4LhFXL6tkrzHXYVy0")
		assert.Equal(t, filepath.Join("deferred", "E", "4LhFXL6tkrzHXYVy0"), queue_path)
	})
}
