package util

import (
	"fmt"
	"strings"
	"time"
)

/*
 * Note that trust levels are not in numerical order.
 */
const DICT_OWNER_UNKNOWN int = (-1)  // ex: unauthenticated tcp, proxy
const DICT_OWNER_TRUSTED int = (0)   // ex: root-owned config file
const DICT_OWNER_UNTRUSTED int = (1) // ex: non-root config file

type DictOwner struct {
	Status int // see below
	Uid    int // use only if status == UNTRUSTED
}

type DictUtf8Backup struct {
	Lookup func(*Dict, string) string
	Update func(*Dict, string, string) int
	Delete func(*Dict, string) int
}

type Dict struct {
	Type     string /* for diagnostics */
	Name     string /* for diagnostics */
	Flags    int    /* see below */
	Lookup   func(string) (string, error)
	Update   func(string, string) int
	Delete   func(string) int
	Sequence func(int, *string, *string)
	Lock     func(int) int
	Close    func()
	LockType int       /* for read/write lock */
	LockFd   int       /* for read/write lock */
	StatFd   int       /* change detection */
	ModTime  time.Time /* mod time at open */
	// foldBuf  *string   /* key folding buffer */
	Owner DictOwner /* provenance */
	Error int       /* last operation only */
	// DICT_JMP_BUF *jbuf;			/* exception handling */
	// struct DICT_UTF8_BACKUP *utf8_backup;	/* see below */
	// struct VSTRING *file_buf;		/* dict_file_to_buf() */
	// struct VSTRING *file_b64;		/* dict_file_to_b64() */
}

var dict_table map[string]*Dict

func dict_eval_lookup(key string, _ int, context any) (string, error) {
	dict_name := context.(string)

	if dict := dict_table[dict_name]; dict != nil {
		val, err := dict.Lookup(key)
		if err != nil {
			MsgFatal("dictionary %s: lookup %s: operation failed: %v", dict_name, key, err)
		}
		return val, nil
	}

	return "", fmt.Errorf("lookup failed: key %s not found", key)
}

func DictEval(dict_name string, value string, recursive bool) string {
	var myname string = "dict_eval"
	var flags int = MAC_EXP_FLAG_NONE

	if recursive {
		flags = MAC_EXP_FLAG_RECURSE
	}
	buf := strings.Builder{}
	status := MacExpand(&buf, value, flags, nil, dict_eval_lookup, dict_name)
	if status&MAC_PARSE_ERROR == MAC_PARSE_ERROR {
		MsgFatal("dictionary %s: macro processing error", dict_name)
	}

	if MsgVerbose > 1 {
		if value != buf.String() {
			MsgInfo("%s: expand %s -> %s", myname, value, buf.String())
		} else {
			MsgInfo("%s: const %s", myname, value)
		}
	}
	return buf.String()
}
