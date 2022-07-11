package util

import (
	"bufio"
	"errors"
	"strings"

	"github.com/gostfix/gostfix/internal/ascii"
)

// Readllines reads one logical line from the named stream.
//
// "blank lines and comments"
//    Empty lines and whitespace-only lines are ignored, as
//    are lines whose first non-whitespace character is a `#'.
//
// "multi-line text"
//    A logical line starts with non-whitespace text. A line that
//    starts with whitespace continues a logical line.
//
// The result value is the input buffer argument or a null pointer
// when no input is found.
//
// Arguments:
//    fp - io.Reader capable of reading a line of text
//
// .IP lineno
//	A null pointer, or a pointer to an integer that is incremented
//	after reading a physical line.
// .IP first_line
//	A null pointer, or a pointer to an integer that will contain
//	the line number of the first non-blank, non-comment line
//	in the result logical line.
// DIAGNOSTICS
//	Warning: a continuation line that does not continue preceding text.
//	The invalid input is ignored, to avoid complicating caller code.

type LogicalLine struct {
	scan      *bufio.Reader
	LineNo    int
	FirstLine int
}

func (ll *LogicalLine) Readllines() (string, error) {
	linebuf := strings.Builder{}
	for {
		line, err := ll.scan.ReadString('\n')
		if len(line) > 0 {
			ll.LineNo++
		}
		if trim := strings.TrimSpace(line); len(trim) > 0 && trim[0] != '#' {
			// non-commented line
			if linebuf.Len() == 0 {
				ll.FirstLine = ll.LineNo
			}
			linebuf.WriteString(strings.TrimSuffix(line, "\n"))
		}
		if err != nil {
			break
		}
		if linebuf.Len() > 0 {
			// peek at the next rune, unfortunately, Peek() gives us n bytes
			// and not a rune, so use the ReadRune/UnreadRune as a peek
			r, _, _ := ll.scan.ReadRune()
			ll.scan.UnreadRune()
			if r != '#' && !ascii.IsSpace(r) {
				break
			}
		}
	}

	/*
	 * Invalid input: continuing text without preceding text. Allowing this
	 * would complicate "postconf -e", which implements its own multi-line
	 * parsing routine. Do not abort, just warn, so that critical programs
	 * like postmap do not leave behind a truncated table.
	 */
	if linebuf.Len() > 0 && ascii.IsSpace(rune(linebuf.String()[0])) {
		MsgWarn("logical line must not start with whitespace", "line", linebuf.String()[:30])
		return ll.Readllines()
	}

	// if we didn't read anyting, return end of file error
	if linebuf.Len() == 0 {
		return "", errors.New("end of file")
	}

	return linebuf.String(), nil
}
