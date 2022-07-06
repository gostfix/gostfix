package util

import (
	"go.uber.org/zap"
)

var MsgVerbose int = 0
var _logger *zap.Logger = nil

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	_logger, _ = config.Build()
}

func SetLogger(logger *zap.Logger) {
	_logger = logger
}

func LoggerConfig(decorators map[string]string) {
	var logger *zap.Logger = _logger
	for k, v := range decorators {
		logger = logger.With(zap.String(k, v))
	}
	_logger = logger
}

func MsgPanic(msg string, kv ...any) {
	_logger.Sugar().Panicw(msg, kv...)
}

func MsgPanicf(format string, a ...any) {
	_logger.Sugar().Panicf(format, a...)
}

func MsgFatal(msg string, kv ...any) {
	_logger.Sugar().Fatalw(msg, kv...)
}

func MsgFatalf(format string, a ...any) {
	_logger.Sugar().Fatalf(format, a...)
}

func MsgError(msg string, kv ...any) {
	_logger.Sugar().Errorw(msg, kv...)
}

func MsgErrorf(format string, a ...any) {
	_logger.Sugar().Errorf(format, a...)
}

func MsgWarn(msg string, kv ...any) {
	_logger.Sugar().Warnw(msg, kv...)
}

func MsgWarnf(format string, a ...any) {
	_logger.Sugar().Warnf(format, a...)
}

func MsgInfo(msg string, kv ...any) {
	_logger.Sugar().Infow(msg, kv...)
}

func MsgInfof(format string, a ...any) {
	_logger.Sugar().Infof(format, a...)
}

func MsgDebug(msg string, kv ...any) {
	_logger.Sugar().Debugw(msg, kv...)
}

func MsgDebugf(format string, a ...any) {
	_logger.Sugar().Debugf(format, a...)
}
