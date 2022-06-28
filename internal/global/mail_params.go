package global

var VarProcname string

const VAR_MULTI_NAME string = "multi_instance_name"
const DEF_MULTI_NAME string = ""

var VarMultiName string

const VAR_SYSLOG_NAME string = "syslog_name"
const DEF_SYSLOG_NAME string = "${" + VAR_MULTI_NAME + "?{$" + VAR_MULTI_NAME + "}:{postfix}}"

var VarSyslogName string
