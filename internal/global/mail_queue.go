package global

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/gostfix/gostfix/internal/ascii"
	"github.com/gostfix/gostfix/internal/util"
)

/*
 * Mail queue names.
 */
const (
	MAIL_QUEUE_MAILDROP = "maildrop"
	MAIL_QUEUE_HOLD     = "hold"
	MAIL_QUEUE_INCOMING = "incoming"
	MAIL_QUEUE_ACTIVE   = "active"
	MAIL_QUEUE_DEFERRED = "deferred"
	MAIL_QUEUE_TRACE    = "trace"
	MAIL_QUEUE_DEFER    = "defer"
	MAIL_QUEUE_BOUNCE   = "bounce"
	MAIL_QUEUE_CORRUPT  = "corrupt"
	MAIL_QUEUE_FLUSH    = "flush"
	MAIL_QUEUE_SAVED    = "saved"
)

/*
 * The long non-repeating queue ID is encoded in an alphabet of 10 digits,
 * 21 upper-case characters, and 21 or fewer lower-case characters. The
 * alphabet is made "safe" by removing all the vowels (AEIOUaeiou). The ID
 * is the concatenation of:
 *
 * - the time in seconds (base 52 encoded, six or more chars),
 *
 * - the time in microseconds (base 52 encoded, exactly four chars),
 *
 * - the 'z' character to separate the time and inode information,
 *
 * - the inode number (base 51 encoded so that it contains no 'z').
 */
const (
	MQID_LG_SEC_BASE  = 52 /* seconds safe alphabet base */
	MQID_LG_SEC_PAD   = 6  /* seconds minimum field width */
	MQID_LG_USEC_BASE = 52 /* microseconds safe alphabet base */
	MQID_LG_USEC_PAD  = 4  /* microseconds exact field width */
	MQID_LG_TIME_PAD  = (MQID_LG_SEC_PAD + MQID_LG_USEC_PAD)
	MQID_LG_INUM_SEP  = 'z' /* time-inode separator */
	MQID_LG_INUM_BASE = 51  /* inode safe alphabet base */
	MQID_LG_INUM_PAD  = 0   /* no padding needed */
)

func MQID_FIND_LG_INUM_SEPARATOR(queue_id string) int {
	return strings.LastIndexByte(queue_id, MQID_LG_INUM_SEP)
}

func MQID_GET_INUM(path string) (uint, error) {
	// check for long queue_id
	if pos := MQID_FIND_LG_INUM_SEPARATOR(path); pos < 0 {
		return MQID_LG_DECODE_INUM(path[pos+1:])
	}
	return MQID_SH_DECODE_INUM(path[MQID_SH_USEC_PAD:])
}

func MQID_LG_ENCODE_SEC(val uint) string {
	return MQID_LG_ENCODE(val, MQID_LG_SEC_BASE, MQID_LG_SEC_PAD)
}

func MQID_LG_ENCODE_USEC(val uint) string {
	return MQID_LG_ENCODE(val, MQID_LG_USEC_BASE, MQID_LG_USEC_PAD)
}

func MQID_LG_ENCODE_INUM(val uint) string {
	return MQID_LG_ENCODE(val, MQID_LG_INUM_BASE, MQID_LG_INUM_PAD)
}

func MQID_LG_DECODE_USEC(str string) (uint, error) {
	return MQID_LG_DECODE(str, MQID_LG_USEC_BASE)
}

func MQID_LG_DECODE_INUM(str string) (uint, error) {
	return MQID_LG_DECODE(str, MQID_LG_INUM_BASE)
}

func MQID_LG_ENCODE(val uint, base uint, padlen int) string {
	return SafeUlToStr(val, base, padlen, '0')
}

func MQID_LG_DECODE(str string, base uint) (uint, error) {
	var end string
	ulval, err := SafeStrToUl(str, &end, base)
	if err == nil && end != "" {
		err = fmt.Errorf("str has remaining chars")
	}
	return ulval, err
}

func MQID_LG_GET_HEX_USEC(queue_id string) string {
	usec_str := queue_id[len(queue_id)-MQID_LG_USEC_PAD:]
	us_val, err := MQID_LG_DECODE_USEC(usec_str)
	if err != nil {
		us_val = 0
	}
	return MQID_SH_ENCODE_USEC(us_val)
}

/*
 * The short repeating queue ID is encoded in upper-case hexadecimal, and is
 * the concatenation of:
 *
 * - the time in microseconds (exactly five chars),
 *
 * - the inode number.
 */
const MQID_SH_USEC_PAD = 5 /* microseconds exact field width */

func MQID_SH_ENCODE_USEC(usec uint) string {
	return fmt.Sprintf("%05X", usec)
}

func MQID_SH_ENCODE_INUM(inum uint) string {
	return fmt.Sprint("%lX", inum)
}

func MQID_SH_DECODE_INUM(str string) (uint, error) {
	ulval, err := strconv.ParseUint(str, 16, 0)
	return uint(ulval), err
}

func MailQueueDir(queue_name string, queue_id string) string {
	if !MailQueueNameOk(queue_name) {
		util.MsgPanic("bad queue name", "queue_name", queue_name)
	}
	if !MailQueueIdOk(queue_id) {
		util.MsgPanic("bad queue id", "queue_id", queue_id)
	}

	path := queue_name
	hash_queue_names := util.MyStrTok(VarHashQueueNames, util.CHARS_COMMA_SP)

	for _, hash := range hash_queue_names {
		if queue_name == hash {
			// handle long queue_id
			if qidx := MQID_FIND_LG_INUM_SEPARATOR(queue_id); qidx != -1 {
				usec_buf := queue_id[:qidx]
				queue_id = MQID_LG_GET_HEX_USEC(usec_buf)
			}
			path = filepath.Join(path, util.DirForest(queue_id, VarHashQueueDepth))
			break
		}
	}

	return path
}

func MailQueuePath(queue_name string, queue_id string) string {
	queue_dir := MailQueueDir(queue_name, queue_id)
	return filepath.Join(queue_dir, queue_id)
}

/* MailQueueNameOk - validate mail queue name */
func MailQueueNameOk(queue_name string) bool {

	if queue_name == "" || len(queue_name) > 100 {
		return false
	}
	scan := util.NewRuneScanner(queue_name)
	scan.Skip(func(r rune) bool { return ascii.IsAlnum(r) })

	return scan.Peek() == scanner.EOF
}

func MailQueueIdOk(queue_id string) bool {
	if queue_id == "" || len(queue_id) > util.VALID_HOSTNAME_LEN {
		return false
	}

	scan := util.NewRuneScanner(queue_id)
	scan.Skip(func(r rune) bool { return ascii.IsAlnum(r) || r == '_' })

	// if we didn't skip all the characters in the queue_id, then
	// it must contain an invalid character.
	return scan.Peek() == scanner.EOF
}

func MailQueueOpen(queue_name string, queue_id string, flags int, mode fs.FileMode) (io.Reader, error) {
	path := MailQueuePath(queue_name, queue_id)
	fd, err := os.OpenFile(path, flags, mode)
	return fd, err
}
