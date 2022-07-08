package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/gostfix/gostfix/internal/global"
	"github.com/gostfix/gostfix/internal/postcat"
	"github.com/gostfix/gostfix/internal/util"
)

func main() {
	var cli struct {
		Body      bool     `kong:"short='b',help='print body records'"`
		ConfigDir string   `kong:"short='c',placeholder='DIR',type='existingdir',env='${CONF_ENV_PATH}',help='configuration directory'"`
		Decimal   bool     `kong:"short='d',help='print decimal record type'"`
		Envelope  bool     `kong:"short='e',help='print envelope records'"`
		Header    bool     `kong:"short='h',help='print header records'"`
		Offset    bool     `kong:"short='o',help='print record offsets'"`
		Queue     bool     `kong:"short='q',help='search queue directory'"`
		Raw       bool     `kong:"short='r',help='do not follow pointers'"`
		Skip      int      `kong:"short='s',placeholder='BYTES',help='skip number of bytes before processing'"`
		Verbose   int      `kong:"short='v',placeholder='LEVEL',env='${CONF_ENV_VERB}',type='counter',help='enable more verbose logging'"`
		Files     []string `kong:"arg"`
	}
	ctx := kong.Parse(&cli,
		kong.Vars{
			"CONF_ENV_PATH": global.CONF_ENV_PATH,
			"CONF_ENV_VERB": global.CONF_ENV_VERB,
		})

	defer util.GetLogger().Sync()
	_ = ctx

	var queue_names = []string{
		global.MAIL_QUEUE_MAILDROP,
		global.MAIL_QUEUE_INCOMING,
		global.MAIL_QUEUE_ACTIVE,
		global.MAIL_QUEUE_DEFERRED,
		global.MAIL_QUEUE_HOLD,
		global.MAIL_QUEUE_SAVED,
	}
	var flags = 0

	if cli.Body {
		flags |= postcat.PC_FLAG_PRINT_BODY
	}
	if cli.Decimal {
		flags |= postcat.PC_FLAG_PRINT_RTYPE_DEC
	} else {
		flags |= postcat.PC_FLAG_PRINT_RTYPE_SYM
	}
	if cli.Envelope {
		flags |= postcat.PC_FLAG_PRINT_ENV
	}
	if cli.Header {
		flags |= postcat.PC_FLAG_PRINT_HEADER
	}
	if cli.Offset {
		flags |= postcat.PC_FLAG_PRINT_OFFSET
	}
	if cli.Queue {
		flags |= postcat.PC_FLAG_SEARCH_QUEUE
	}
	if cli.Raw {
		flags |= postcat.PC_FLAG_RAW
	}

	if flags&postcat.PC_MASK_PRINT_ALL == 0 {
		flags |= postcat.PC_MASK_PRINT_ALL
	}

	if cli.ConfigDir != "" {
		global.VarConfigDir = cli.ConfigDir
		os.Setenv(global.CONF_ENV_PATH, cli.ConfigDir)
	}

	global.MailConfRead()

	dict := util.DictFindForUpdate(global.CONFIG_DICT)
	var nam string
	var val string
	for dict.Sequence(util.DICT_SEQ_FUN_FIRST, &nam, &val); nam != ""; dict.Sequence(util.DICT_SEQ_FUN_NEXT, &nam, &val) {
		util.MsgInfo("sequence", "name", nam, "value", val)
	}
	_ = queue_names

	util.MsgInfo("postcat", "cli", cli)
}
