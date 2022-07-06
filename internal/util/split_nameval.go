package util

import (
	"errors"
	"strings"
	"text/scanner"

	"github.com/gostfix/gostfix/internal/ascii"
)

// SplitNameVal takes a logical line from Readlline and expects
// text of the form "name = value" or "name =".
//
// The buffer argument is broken up into name and value substrings.
func SplitNameVal(buf string) (string, string, error) {
	scan := NewRuneScanner(buf)

	scan.Skip(ascii.IsSpace) // find name begin
	if scan.Peek() == scanner.EOF {
		return "", "", errors.New("missing attribute name")
	}
	scan.Span(func(r rune) bool { return !ascii.IsSpace(r) && r != '=' }) // find name end
	name := scan.Emit()

	scan.Skip(ascii.IsSpace) // skip blanks before '='
	if scan.Peek() != '=' {
		return "", "", errors.New("missing '=' after attribute name")
	}
	scan.Next() // skip over '='

	scan.Skip(ascii.IsSpace) // skip leading blanks
	scan.End(false)
	value := scan.Emit()
	value = strings.TrimSpace(value)

	return name, value, nil
}
