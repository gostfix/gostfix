package global

import (
	"os"
	"path/filepath"

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

func MailConfCheckDir(config_dir string) {
	// TODO(alf): we may skip this check and allow any user
	// to load the config file
}

func MailConfRead() {
	MailConfSuck()
	MailParamsInit()
}

func MailConfSuck() {
	var config_dir string
	if config_dir = os.Getenv(CONF_ENV_PATH); config_dir == "" {
		config_dir = DEF_CONFIG_DIR
	}
	VarConfigDir = config_dir
	SetMailConfStr(VAR_CONFIG_DIR, VarConfigDir)

	// TODO(alf): This causes our config option to pass on the command
	// line to fail unless we run as root.  Therefore, just disable
	// for now.
	// if VarConfigDir != DEF_CONFIG_DIR && util.Unsafe() {
	// 	MailConfCheckDir(VarConfigDir)
	// }
	file := filepath.Join(VarConfigDir, "main.cf")
	if err := util.DictLoadFileXt(CONFIG_DICT, file); err != nil {
		util.MsgFatal("DictLoadFileXt failed", "file", file, "error", err)
	}
}

func MailConfFlush() {
	if dict := util.DictHandle(CONFIG_DICT); dict != nil {
		util.DictUnregister(CONFIG_DICT)
	}
}
func MailConfEval(str string) string {
	return util.DictEval(CONFIG_DICT, str, true)
}

func MailConfEvalOnce(str string) string {
	return util.DictEval(CONFIG_DICT, str, false)
}

func MailConfLookup(name string) string {
	return util.DictLookup(CONFIG_DICT, name)
}

func MailConfLookupEval(name string) string {
	var value string = ""
	if value = util.DictLookup(CONFIG_DICT, name); value != "" {
		value = util.DictEval(CONFIG_DICT, value, true)
	}
	return value
}

func MailConfUpdate(key string, value string) {
	util.DictUpdate(CONFIG_DICT, key, value)
}
