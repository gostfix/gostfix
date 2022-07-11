package global

import (
	"testing"

	"github.com/gostfix/gostfix/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestMailParamSplit(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	util.SetLogger(observedLogger)

	t.Run("env with equals", func(t *testing.T) {
		params := MailParamSplit("test", "TZ PATH=/bin:/usr/bin XAUTHORITY")
		assert.Equal(t, []string{"TZ", "PATH=/bin:/usr/bin", "XAUTHORITY"}, params)
	})

	t.Run("env with quote beginning", func(t *testing.T) {
		params := MailParamSplit("test", "{ LESS=-m -C -s -f -e } TZ XAUTHORITY")
		assert.Equal(t, []string{"LESS=-m -C -s -f -e", "TZ", "XAUTHORITY"}, params)
	})

	t.Run("env with quote end", func(t *testing.T) {
		params := MailParamSplit("test", "TZ XAUTHORITY { LESS=-m -C -s -f -e }")
		assert.Equal(t, []string{"TZ", "XAUTHORITY", "LESS=-m -C -s -f -e"}, params)
	})

	t.Run("env with no end quote", func(t *testing.T) {
		params := MailParamSplit("test", "TZ { LESS=-m -C -s -f -e  XAUTHORITY")
		assert.Equal(t, []string{"TZ", "LESS=-m -C -s -f -e  XAUTHORITY"}, params)
		logs := observedLogs.TakeAll()
		assert.Equal(t, 1, len(logs))
		assert.Equal(t, zapcore.Level(1), logs[0].Level)
		assert.Equal(t, "missing closing '}' in text \"{ LESS=-m -C -s -f -e  XAUTHORITY\"", logs[0].Message)
	})

	t.Run("env with equals", func(t *testing.T) {
		params := MailParamSplit("test", "TZ { LESS=-m -C -s -f -e }x XAUTHORITY")
		assert.Equal(t, []string{"TZ", "LESS=-m -C -s -f -e", "XAUTHORITY"}, params)
		logs := observedLogs.TakeAll()
		assert.Equal(t, 1, len(logs))
		assert.Equal(t, zapcore.Level(1), logs[0].Level)
		assert.Equal(t, "syntax error after '}' in \"{ LESS=-m -C -s -f -e }x\"", logs[0].Message)
	})
}
