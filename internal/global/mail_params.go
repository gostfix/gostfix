package global

/*
 * Name used when this mail system announces itself.
 */
const VAR_MAIL_NAME = "mail_name"
const DEF_MAIL_NAME = "Postfix"

var VarMailName string

/*
 * Location of configuration files.
 */
const VAR_CONFIG_DIR = "config_directory"
const DEF_CONFIG_DIR = "/etc/postfix"

var VarConfigDir string

/*
 * You want to be helped or not.
 */
const VAR_HELPFUL_WARNINGS = "helpful_warnings"
const DEF_HELPFUL_WARNINGS = 1

var VarHelpfulWarnings bool

/*
 * You want to be helped or not.
 */
const VAR_SHOW_UNK_RCPT_TABLE = "show_user_unknown_table_name"
const DEF_SHOW_UNK_RCPT_TABLE = 1

var VarShowUnkRcptTable bool

/*
 * Compatibility level and migration support. Update postconf(5),
 * COMPATIBILITY_README, global/mail_params.[hc] and conf/main.cf when
 * updating the current compatibility level.
 */
const (
	COMPAT_LEVEL_0    = "0"
	COMPAT_LEVEL_1    = "1"
	COMPAT_LEVEL_2    = "2"
	COMPAT_LEVEL_3_6  = "3.6"
	LAST_COMPAT_LEVEL = COMPAT_LEVEL_3_6
)

const VAR_COMPAT_LEVEL = "compatibility_level"
const DEF_COMPAT_LEVEL = COMPAT_LEVEL_0

var var_compatibility_level string

const VAR_PROCNAME string = "process_name"

var VarProcname string

const VAR_MULTI_NAME string = "multi_instance_name"
const DEF_MULTI_NAME string = ""

var VarMultiName string

const VAR_SYSLOG_NAME string = "syslog_name"
const DEF_SYSLOG_NAME string = "${" + VAR_MULTI_NAME + "?{$" + VAR_MULTI_NAME + "}:{postfix}}"

var VarSyslogName string

/*
 * Location of the mail queue directory tree.
 */
const VAR_QUEUE_DIR = "queue_directory"
const DEF_QUEUE_DIR = "/var/spool/postfix"

var VarQueueDir string

/*
 * Queue management: what queues are hashed behind a forest of
 * subdirectories, and how deep the forest is.
 */
const VAR_HASH_QUEUE_NAMES = "hash_queue_names"
const DEF_HASH_QUEUE_NAMES = "deferred, defer"

var VarHashQueueNames string

const VAR_HASH_QUEUE_DEPTH = "hash_queue_depth"
const DEF_HASH_QUEUE_DEPTH = 1

var VarHashQueueDepth int

/*
 * Short queue IDs contain the time in microseconds and file inode number.
 * Long queue IDs also contain the time in seconds.
 */
const VAR_LONG_QUEUE_IDS = "enable_long_queue_ids"
const DEF_LONG_QUEUE_IDS = false

var VarLongQueueIds bool

const VAR_SMTPD_SOFT_ERLIM string = "smtpd_soft_error_limit"
const DEF_SMTPD_SOFT_ERLIM string = "10"

var VarSmtpdSoftErlim int

const VAR_MASTER_DISABLE string = "master_service_disable"
const DEF_MASTER_DISABLE string = ""

/*
 * Environmental management - what Postfix imports from the external world,
 * and what Postfix exports to the external world.
 */
const VAR_IMPORT_ENVIRON string = "import_environment"
const DEF_IMPORT_ENVIRON string = `MAIL_CONFIG MAIL_DEBUG MAIL_LOGTAG 
	TZ XAUTHORITY DISPLAY LANG=C 
	POSTLOG_SERVICE POSTLOG_HOSTNAME`

var VarImportEnviron string

func MailParamsInit() {
	var compat_level_defaults = []CONFIG_STR_TABLE{
		{VAR_COMPAT_LEVEL, DEF_COMPAT_LEVEL, &var_compatibility_level, 0, 0},
	}

	_ = compat_level_defaults
}
