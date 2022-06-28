package util

import "strings"

const MAC_PARSE_LITERAL int = 1
const MAC_PARSE_EXPR int = 2
const MAC_PARSE_VARNAME int = MAC_PARSE_EXPR // 2.1 compatibility

const MAC_PARSE_OK int = 0
const MAC_PARSE_ERROR int = (1 << 0)
const MAC_PARSE_UNDEF int = (1 << 1)
const MAC_PARSE_USER int = 2 // start user definitions

func MacParse(value string, action func(int, string, interface{}) int, context interface{}) {
	myname := "MacParse"
	if MsgVerbose > 1 {
		MsgInfo("%s : %s", myname, value)
	}

	var buf strings.Builder
	var consumeDollar bool
	for i, c := range value {
		if c != '$' || consumeDollar {
			buf.WriteRune(c)
			consumeDollar = false
		} else if value[i+1] == '$' {
			consumeDollar = true
		} else {
		}
	}
}
