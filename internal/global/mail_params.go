package global

const VAR_PROCNAME string = "process_name"

var VarProcname string

const VAR_MULTI_NAME string = "multi_instance_name"
const DEF_MULTI_NAME string = ""

var VarMultiName string

const VAR_SYSLOG_NAME string = "syslog_name"
const DEF_SYSLOG_NAME string = "${" + VAR_MULTI_NAME + "?{$" + VAR_MULTI_NAME + "}:{postfix}}"

var VarSyslogName string

const VAR_SMTPD_SOFT_ERLIM string = "smtpd_soft_error_limit"
const DEF_SMTPD_SOFT_ERLIM string = "10"

var VarSmtpdSoftErlim int

const VAR_MASTER_DISABLE string = "master_service_disable"
const DEF_MASTER_DISABLE string = ""
