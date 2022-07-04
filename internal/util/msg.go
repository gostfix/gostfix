package util

import "github.com/sirupsen/logrus"

var MsgVerbose int = 0

func MsgPanic(format string, a ...any) {
	logrus.Panicf(format, a...)
}

func MsgFatal(format string, a ...any) {
	logrus.Fatalf(format, a...)
}

func MsgError(format string, a ...any) {
	logrus.Errorf(format, a...)
}

func MsgWarn(format string, a ...any) {
	logrus.Warnf(format, a...)
}

func MsgInfo(format string, a ...any) {
	logrus.Infof(format, a...)
}
