package global

import (
	"github.com/gostfix/gostfix/internal/util"
)

/*
 * Well known names. These are not configurable. One has to start somewhere.
 */
const CONFIG_DICT string = "mail_dict" // global Postfix dictionary

/*
 * Environment variables.
 */
const CONF_ENV_PATH string = "MAIL_CONFIG"   // config database
const CONF_ENV_VERB string = "MAIL_VERBOSE"  // verbose mode on
const CONF_ENV_DEBUG string = "MAIL_DEBUG"   // live debugging
const CONF_ENV_LOGTAG string = "MAIL_LOGTAG" // instance name

/*
 * External representation for booleans.
 */
const CONFIG_BOOL_YES string = "yes"
const CONFIG_BOOL_NO string = "no"

type CONFIG_BASE_TABLE struct {
	Min int
	Max int
}

type CONFIG_STR_TABLE struct {
	Name   string
	Defval string
	Target *string
	Min    int
	Max    int
}

type CONFIG_RAW_TABLE struct {
	Name   string
	Defval string
	Target *string
	Min    int
	Max    int
}

type CONFIG_INT_TABLE struct {
	Name   string
	Defval int
	Target *int
	Min    int
	Max    int
}

type CONFIG_LONG_TABLE struct {
	Name   string
	Defval int
	Target *int
	Min    int
	Max    int
}

type CONFIG_BOOL_TABLE struct {
	Name   string
	Defval bool
	Target *bool
}

type CONFIG_TIME_TABLE struct {
	Name   string
	Defval int
	Target *int
	Min    int
	Max    int
}

type CONFIG_NINT_TABLE struct {
	Name   string
	Defval string
	Target *int
	Min    int
	Max    int
}

type CONFIG_NBOOL_TABLE struct {
	Name   string
	Defval string
	Target *bool
}

func MailConfEval(str string) string {
	return util.DictEval(CONFIG_DICT, str, true)
}
