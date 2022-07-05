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

/*
 * See dict_open.c embedded manpage for flag definitions.
 */
const (
	DICT_FLAG_NONE        int = (0)
	DICT_FLAG_DUP_WARN    int = (1 << 0) /* warn about dups if not supported */
	DICT_FLAG_DUP_IGNORE  int = (1 << 1) /* ignore dups if not supported */
	DICT_FLAG_TRY0NULL    int = (1 << 2) /* do not append 0 to key/value */
	DICT_FLAG_TRY1NULL    int = (1 << 3) /* append 0 to key/value */
	DICT_FLAG_FIXED       int = (1 << 4) /* fixed key map */
	DICT_FLAG_PATTERN     int = (1 << 5) /* keys are patterns */
	DICT_FLAG_LOCK        int = (1 << 6) /* use temp lock before access */
	DICT_FLAG_DUP_REPLACE int = (1 << 7) /* replace dups if supported */
	DICT_FLAG_SYNC_UPDATE int = (1 << 8) /* sync updates if supported */
	DICT_FLAG_DEBUG       int = (1 << 9) /* log access */
	/*DICT_FLAG_FOLD_KEY int = (1<<10) /* lowercase the lookup key */
	DICT_FLAG_NO_REGSUB       int = (1 << 11) /* disallow regexp substitution */
	DICT_FLAG_NO_PROXY        int = (1 << 12) /* disallow proxy mapping */
	DICT_FLAG_NO_UNAUTH       int = (1 << 13) /* disallow unauthenticated data */
	DICT_FLAG_FOLD_FIX        int = (1 << 14) /* case-fold key with fixed-case map */
	DICT_FLAG_FOLD_MUL        int = (1 << 15) /* case-fold key with multi-case map */
	DICT_FLAG_FOLD_ANY        int = (DICT_FLAG_FOLD_FIX | DICT_FLAG_FOLD_MUL)
	DICT_FLAG_OPEN_LOCK       int = (1 << 16) /* perm lock if not multi-writer safe */
	DICT_FLAG_BULK_UPDATE     int = (1 << 17) /* optimize for bulk updates */
	DICT_FLAG_MULTI_WRITER    int = (1 << 18) /* multi-writer safe map */
	DICT_FLAG_UTF8_REQUEST    int = (1 << 19) /* activate UTF-8 if possible */
	DICT_FLAG_UTF8_ACTIVE     int = (1 << 20) /* UTF-8 proxy layer is present */
	DICT_FLAG_SRC_RHS_IS_FILE int = (1 << 21) /* Map source RHS is a file */
)

const DICT_FLAG_UTF8_MASK int = (DICT_FLAG_UTF8_REQUEST)

/* IMPORTANT: Update the dict_mask[] table when the above changes */

/*
 * The subsets of flags that control how a map is used. These are relevant
 * mainly for proxymap support. Note: some categories overlap.
 *
 * DICT_FLAG_IMPL_MASK - flags that are set by the map implementation itself.
 *
 * DICT_FLAG_PARANOID - requestor flags that forbid the use of insecure map
 * types for security-sensitive operations. These flags are checked by the
 * map implementation itself upon open, lookup etc. requests.
 *
 * DICT_FLAG_RQST_MASK - all requestor flags, including paranoid flags, that
 * the requestor may change between open, lookup etc. requests. These
 * specify requestor properties, not map properties.
 *
 * DICT_FLAG_INST_MASK - none of the above flags. The requestor may not change
 * these flags between open, lookup, etc. requests (although a map may make
 * changes to its copy of some of these flags). The proxymap server opens
 * only one map instance for all client requests with the same values of
 * these flags, and the proxymap client uses its own saved copy of these
 * flags. DICT_FLAG_SRC_RHS_IS_FILE is an example of such a flag.
 */
const DICT_FLAG_PARANOID int = (DICT_FLAG_NO_REGSUB | DICT_FLAG_NO_PROXY | DICT_FLAG_NO_UNAUTH)
const DICT_FLAG_IMPL_MASK int = (DICT_FLAG_FIXED | DICT_FLAG_PATTERN | DICT_FLAG_MULTI_WRITER)
const DICT_FLAG_RQST_MASK int = (DICT_FLAG_FOLD_ANY | DICT_FLAG_LOCK |
	DICT_FLAG_DUP_REPLACE | DICT_FLAG_DUP_WARN |
	DICT_FLAG_DUP_IGNORE | DICT_FLAG_SYNC_UPDATE |
	DICT_FLAG_PARANOID | DICT_FLAG_UTF8_MASK)

//const DICT_FLAG_INST_MASK int = ~(DICT_FLAG_IMPL_MASK | DICT_FLAG_RQST_MASK)

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
	// DICT_JMP_BUF *jbuf;   /* exception handling */
	// struct DICT_UTF8_BACKUP *utf8_backup; /* see below */
	// struct VSTRING *file_buf;  /* dict_file_to_buf() */
	// struct VSTRING *file_b64;  /* dict_file_to_b64() */
}

/*
 * dict->error values. Errors must be negative; smtpd_check depends on this.
 */
const (
	DICT_ERR_NONE   int = 0    /* no error */
	DICT_ERR_RETRY  int = (-1) /* soft error */
	DICT_ERR_CONFIG int = (-2) /* configuration error */
)

/*
 * Sequence function types.
 */
const (
	DICT_SEQ_FUN_FIRST = 0 /* set cursor to first record */
	DICT_SEQ_FUN_NEXT  = 1 /* set cursor to next record */
)

var dict_table map[string]*Dict

func dict_eval_lookup(key string, _ int, context any) (string, error) {
	dict_name := context.(string)

	if dict := dict_table[dict_name]; dict != nil {
		val, err := dict.Lookup(key)
		if err != nil {
			MsgFatal("operation failed", "dictionary", dict_name, "lookup", key, "error", err)
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
		MsgFatal("macro processing error", "dictionary", dict_name)
	}

	if MsgVerbose > 1 {
		if value != buf.String() {
			MsgInfo("expand", "function", myname, "from", value, "to", buf.String())
		} else {
			MsgInfo("const", "function", myname, "value", value)
		}
	}
	return buf.String()
}
