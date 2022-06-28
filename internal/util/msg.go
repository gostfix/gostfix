package util

import "github.com/sirupsen/logrus"

var MsgVerbose int = 0

func MsgFatal(format string, a ...any) {
	logrus.Fatalf(format, a...)
}

func MsgInfo(format string, a ...any) {
	logrus.Infof(format, a...)
}
