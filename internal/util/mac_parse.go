// This module recognizes macro expressions in null-terminated
// strings.  Macro expressions have the form $name, $(text) or
// ${text}. A macro name consists of alphanumerics and/or
// underscore. Text other than macro expressions is treated
// as literal text.

package util

import (
	"strings"
	"text/scanner"

	"github.com/gostfix/gostfix/internal/ascii"
)

const (
	MAC_PARSE_LITERAL int = 1              // The content of buf is literal text.
	MAC_PARSE_EXPR    int = 2              // The content of buf is a macro expression.
	MAC_PARSE_VARNAME int = MAC_PARSE_EXPR // 2.1 compatibility
)

const (
	MAC_PARSE_OK    int = 0        // No errors detected during macro parsing
	MAC_PARSE_ERROR int = (1 << 0) // A parsing error was detected.
	MAC_PARSE_UNDEF int = (1 << 1) // A macro was expanded but not defined.
	MAC_PARSE_USER  int = 2        // start user definitions
)

var paren map[rune]rune = map[rune]rune{
	'{': '}',
	'(': ')',
}

// MacParse locates macro references in string
//
// MacParse breaks up its string argument into macro references
// and other text, and invokes the action routine for each item
// found.  With each action routine call, the type argument
// indicates what was found, buf contains a copy of the text
// found, and context is passed on unmodified from the caller.
// The application is at liberty to clobber buf.
//
// Values for type argument in action:
//     MAC_PARSE_LITERAL  // We parsed non-macro literal text
//     MAC_PARSE_EXPR     // either a bare macro name without the preceding "$",
//                        // or all the text inside $() or ${}.
func MacParse(value string, action func(int, string, interface{}) int, context interface{}) int {
	myname := "MacParse"
	if MsgVerbose > 1 {
		MsgInfo("entry", "function", myname, "value", value)
	}

	// our return value/error
	var status int = 0

	// Temporary storage to buildup our text or macro variable
	var buf strings.Builder

	// Scanner gives us access to various methods that allow us to
	// parse the string one character at a time.  Specifically, we
	// want Next() and Peek() which return a rune or EOF if an error
	scan := NewRuneScanner(value)

	var ch rune
	for ch = scan.Next(); ch != scanner.EOF; ch = scan.Next() {
		if ch != '$' { // ordinary character
			buf.WriteRune(ch)
		} else if ch = scan.Peek(); ch == '$' { // $$ becomes $
			buf.WriteRune(scan.Next())
		} else { // found bare $
			if buf.Len() > 0 {
				status |= action(MAC_PARSE_LITERAL, buf.String(), context)
				buf.Reset()
			}

			if ch = scan.Next(); ch == scanner.EOF {
				status |= MAC_PARSE_ERROR
				break
			} else if close, hasParen := paren[ch]; hasParen { // ${x} or $(x)
				level := 1
				open := ch

				for level > 0 {
					ch = scan.Scan()
					if ch == scanner.EOF {
						MsgWarn("truncated macro reference", "value", value)
						status |= MAC_PARSE_ERROR
						break
					} else if ch == open {
						level++
					} else if ch == close {
						level--
					}
				}
				if status&MAC_PARSE_ERROR == MAC_PARSE_ERROR {
					break
				}
				if level == 0 {
					scan.Push() // push back the closing bracket
				}
				buf.WriteString(scan.Emit())
				scan.Next() // eat the closing bracket
			} else { // plain $x
				buf.WriteRune(ch)
				scan.Span(func(r rune) bool { return ascii.IsAlnum(r) || r == '_' })
				buf.WriteString(scan.Emit())
			}
			if buf.Len() == 0 {
				status |= MAC_PARSE_ERROR
				MsgWarn("empty macro name", "value", value)
				break
			}
			status |= action(MAC_PARSE_EXPR, buf.String(), context)
			buf.Reset()
		}
	}
	if buf.Len() > 0 && (status&MAC_PARSE_ERROR == 0) {
		status |= action(MAC_PARSE_LITERAL, buf.String(), context)
		buf.Reset()
	}

	return status
}
