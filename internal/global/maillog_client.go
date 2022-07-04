package global

import "github.com/gostfix/gostfix/internal/util"

const (
	MAILLOG_CLIENT_FLAG_NONE               = 0
	MAILLOG_CLIENT_FLAG_LOGWRITER_FALLBACK = 1 << iota
)

func MailLogClientInit(progname string, flags int) {
	util.LoggerConfig(map[string]string{"program": progname})
}
