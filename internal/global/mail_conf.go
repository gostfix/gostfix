package global

import "github.com/gostfix/gostfix/internal/util"

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

func MailConfEval(str string) string {
	return util.DictEval(CONFIG_DICT, str, true)
}
