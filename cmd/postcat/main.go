package main

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/gostfix/gostfix/internal/global"
	"github.com/gostfix/gostfix/internal/postcat"
	"github.com/gostfix/gostfix/internal/util"
)

type flagConfig string

func (c flagConfig) AfterApply() error {
	if c != "" && c != flagConfig(global.VarConfigDir) {
		global.VarConfigDir = string(c)
		return os.Setenv(global.CONF_ENV_PATH, string(c))
	}
	return nil
}

type flagVerbose int

func (v flagVerbose) AfterApply() error {
	if int(v) != util.MsgVerbose {
		util.MsgVerbose = int(v)
	}
	return nil
}

type flagBodBit bool
type flagDecBit bool
type flagEnvBit bool
type flagHdrBit bool
type flagOffBit bool
type flagQueBit bool
type flagRawBit bool

func (f flagBodBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_PRINT_BODY
	}
	return nil
}

func (f flagDecBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_PRINT_RTYPE_DEC
		*flags &= ^postcat.PC_FLAG_PRINT_RTYPE_SYM
	} else {
		*flags &= ^postcat.PC_FLAG_PRINT_RTYPE_DEC
		*flags |= postcat.PC_FLAG_PRINT_RTYPE_SYM
	}
	return nil
}

func (f flagEnvBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_PRINT_ENV
	}
	return nil
}

func (f flagHdrBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_PRINT_HEADER
	}
	return nil
}

func (f flagOffBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_PRINT_OFFSET
	}
	return nil
}

func (f flagQueBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_SEARCH_QUEUE
	}
	return nil
}

func (f flagRawBit) AfterApply(flags *int) error {
	if f {
		*flags |= postcat.PC_FLAG_RAW
	}
	return nil
}

func main() {
	var flags int = postcat.PC_FLAG_PRINT_RTYPE_SYM

	var cli struct {
		Body      flagBodBit  `kong:"short='b',help='print body records'"`
		ConfigDir flagConfig  `kong:"short='c',placeholder='DIR',type='existingdir',env='${CONF_ENV_PATH}',help='configuration directory'"`
		Decimal   flagDecBit  `kong:"short='d',help='print decimal record type'"`
		Envelope  flagEnvBit  `kong:"short='e',help='print envelope records'"`
		Header    flagHdrBit  `kong:"short='h',help='print header records'"`
		Offset    flagOffBit  `kong:"short='o',help='print record offsets'"`
		Queue     flagQueBit  `kong:"short='q',help='search queue directory'"`
		Raw       flagRawBit  `kong:"short='r',help='do not follow pointers'"`
		Skip      int         `kong:"short='s',placeholder='BYTES',help='skip number of bytes before processing'"`
		Verbose   flagVerbose `kong:"short='v',placeholder='LEVEL',env='${CONF_ENV_VERB}',type='counter',help='enable more verbose logging'"`
		Files     []string    `kong:"arg,optional"`
	}
	ctx := kong.Parse(&cli,
		kong.Vars{
			"CONF_ENV_PATH": global.CONF_ENV_PATH,
			"CONF_ENV_VERB": global.CONF_ENV_VERB,
		}, kong.Bind(&flags))

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

	// print all fields if we don't have a specific field set.
	if flags&postcat.PC_MASK_PRINT_ALL == 0 {
		flags |= postcat.PC_MASK_PRINT_ALL
	}

	global.MailConfRead()
	import_env := global.MailParamSplit(global.VAR_IMPORT_ENVIRON, global.VarImportEnviron)
	util.UpdateEnv(import_env)

	if len(cli.Files) == 0 {
		postcat.Postcat(os.Stdin, flags)
	} else if flags&postcat.PC_FLAG_SEARCH_QUEUE == postcat.PC_FLAG_SEARCH_QUEUE {
		if err := os.Chdir(global.VarQueueDir); err != nil {
			util.MsgFatal("unable to change directory", "err", err, "dir", global.VarQueueDir)
		}
		for _, fname := range cli.Files {
			if !global.MailQueueIdOk(fname) {
				util.MsgFatal("bad mail queue id", "queue_id", fname)
			}
			var fp io.Reader = nil
			var err error
			// NOTE(alf):  I'm getting rid of the 2 retries
			// for tries := 0; fp == nil && tries < 2; tries++ {
			for _, queue := range queue_names {
				fp, err = global.MailQueueOpen(queue, fname, os.O_RDONLY, 0)
				if err == nil {
					break
				}
			}
			//}

			if err != nil {
				util.MsgFatal("failed to open queue file", "queue_id", fname, "err", err)
			}
			postcat.Postcat(fp, flags)
		}
	} else {
		for _, fname := range cli.Files {
			fd, err := os.Open(fname)
			if err != nil {
				util.MsgFatal("postcat open file failed", "err", err, "file", fname)
			}
			postcat.Postcat(fd, flags)
			if err = fd.Close(); err != nil {
				util.MsgWarn("postcat close file failed", "err", err, "file", fname)
			}
		}
	}

	_ = queue_names
}
