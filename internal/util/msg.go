package util

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MsgVerbose int = 0
var zapLog *zap.Logger = nil

func init() {
	// config := zap.NewProductionConfig()
	config := zap.NewDevelopmentConfig()

	// set output path to stdout so we can pipe results through jq or greap
	config.OutputPaths = []string{"stdout"}
	// log time in a readable format (e.g. 2022-07-07T00:52:25.135303Z), this sets to UTC timezone
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339Nano))
	})

	// only enable color encoding if the output is not full json
	// NOTE(alf): this should be done just before config.Build to make
	// sure the encode level isn't changed after this check.
	// zapcore.CapitalColorLevelEncoder with config.Encoding == "json" causes
	// console escape characters in the json string.
	if config.Encoding == "console" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	zapLog, _ = config.Build()
	// Since we wrap our logger with custom methods, we want to skip 1 level up the
	// call chain to report the proper logging code location.
	zapLog = zapLog.WithOptions(zap.AddCallerSkip(1))
}

func SetLogger(logger *zap.Logger) {
	zapLog = logger
}

func GetLogger() *zap.Logger {
	return zapLog
}

func LoggerConfig(decorators map[string]string) {
	var logger *zap.Logger = zapLog
	for k, v := range decorators {
		logger = logger.With(zap.String(k, v))
	}
	zapLog = logger
}

func MsgDPanic(msg string, kv ...any) {
	zapLog.Sugar().DPanicw(msg, kv...)
}

func MsgDPanicf(format string, a ...any) {
	zapLog.Sugar().DPanicf(format, a...)
}

func MsgPanic(msg string, kv ...any) {
	zapLog.Sugar().Panicw(msg, kv...)
}

func MsgPanicf(format string, a ...any) {
	zapLog.Sugar().Panicf(format, a...)
}

func MsgFatal(msg string, kv ...any) {
	zapLog.Sugar().Fatalw(msg, kv...)
}

func MsgFatalf(format string, a ...any) {
	zapLog.Sugar().Fatalf(format, a...)
}

func MsgError(msg string, kv ...any) {
	zapLog.Sugar().Errorw(msg, kv...)
}

func MsgErrorf(format string, a ...any) {
	zapLog.Sugar().Errorf(format, a...)
}

func MsgWarn(msg string, kv ...any) {
	zapLog.Sugar().Warnw(msg, kv...)
}

func MsgWarnf(format string, a ...any) {
	zapLog.Sugar().Warnf(format, a...)
}

func MsgInfo(msg string, kv ...any) {
	zapLog.Sugar().Infow(msg, kv...)
}

func MsgInfof(format string, a ...any) {
	zapLog.Sugar().Infof(format, a...)
}

func MsgDebug(msg string, kv ...any) {
	zapLog.Sugar().Debugw(msg, kv...)
}

func MsgDebugf(format string, a ...any) {
	zapLog.Sugar().Debugf(format, a...)
}
