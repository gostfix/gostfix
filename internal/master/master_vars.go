package master

import "github.com/gostfix/gostfix/internal/global"

var VarMasterDisable string

func MasterVarsInit() {
	str_table := []global.CONFIG_STR_TABLE{
		{Name: global.VAR_MASTER_DISABLE, Defval: global.DEF_MASTER_DISABLE, Target: &VarMasterDisable, Max: 0, Min: 0},
	}

	// MailConfFlush()
	// SetMailConfStr(global.VAR_PROCNAME, global.VarProcname)
	// MailConfRead()
	// GetMailConfStrTable(str_table)
	// GetMailConfTimeTable(time_table)

	_ = str_table
}
