package global

import (
	"fmt"
	"strings"

	"github.com/gostfix/gostfix/internal/util"
)

var canon_name string

func MailTask(argv0 string) string {
	var tag string
	if argv0 == "" && canon_name == "" {
		argv0 = "unknown"
	}

	if argv0 != "" {
		slash := strings.Split(argv0, "/")
		argv0 = slash[len(slash)-1]

		if tag = util.SafeGetenv(CONF_ENV_LOGTAG); tag == "" {
			if VarSyslogName != "" {
				tag = VarSyslogName
			} else {
				tag = MailConfEval(DEF_SYSLOG_NAME)
			}
		}
		canon_name = fmt.Sprintf("%s/%s", tag, argv0)
	}
	return canon_name
}
