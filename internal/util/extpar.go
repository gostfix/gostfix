package util

import (
	"fmt"
	"strings"
)

const (
	EXTPAR_FLAG_NONE    = (0)
	EXTPAR_FLAG_STRIP   = (1 << 0) /* "{ text }" -> "text" */
	EXTPAR_FLAG_EXTRACT = (1 << 1) /* hint from caller's caller */
)

func ExtPar(val string, parens string, flags int) (string, error) {
	var text string = ""
	var err error = nil

	if len(parens) != 2 {
		MsgFatal("parens must contain two characaters as a delimter", "parens", parens)
	}

	if val[0] != parens[0] {
		return "", fmt.Errorf("no '%c' at start of text in \"%s\"", parens[0], val)
	} else if end := strings.IndexByte(val, parens[1]); end < 0 {
		text = val[1:]
		err = fmt.Errorf("missing closing '%c' in text \"%s\"", parens[1], val)
	} else if end != len(val)-1 {
		text = val[1:end]
		err = fmt.Errorf("syntax error after '%c' in \"%s\"", parens[1], val)
	} else {
		text = val[1:end]
	}

	if flags&EXTPAR_FLAG_STRIP == EXTPAR_FLAG_STRIP {
		text = strings.TrimSpace(text)
	}

	return text, err
}
