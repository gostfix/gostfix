package util

import (
	"strings"
	"text/scanner"
)

// MyStrTok - safe tokenizer
func MyStrTok(src string, sep string) []string {
	tokens := []string{}
	scan := NewRuneScanner(src)
	for {
		if scan.Peek() == scanner.EOF {
			break
		}
		scan.Skip(func(r rune) bool {
			return strings.ContainsRune(sep, r)
		})
		scan.Span(func(r rune) bool {
			return !strings.ContainsRune(sep, r)
		})
		tok := scan.Emit()
		if len(tok) > 0 {
			tokens = append(tokens, tok)
		}
	}

	return tokens
}

// MyStrTokQ - safe tokenizer with quoting support
func MyStrTokQ(src string, sep string, parens string) []string {
	tokens := []string{}
	rparens := []rune(parens)

	scan := NewRuneScanner(src)

	// skip over leading delimiters
	scan.Skip(func(r rune) bool {
		return strings.ContainsRune(sep, r)
	})

	var level int = 0
	for {
		ch := scan.Scan()

		if ch == rparens[0] {
			level++
		} else if level > 0 && ch == rparens[1] {
			level--
		} else if level == 0 && strings.ContainsRune(sep, ch) {
			scan.Push() // we don't want the seperator rune in our tok
			tok := scan.Emit()
			if len(tok) > 0 {
				tokens = append(tokens, tok)
			}
			// skip over leading delimiters
			scan.Skip(func(r rune) bool {
				return strings.ContainsRune(sep, r)
			})
		} else if ch == scanner.EOF {
			break
		}
	}

	// pickup the last token
	tok := scan.Emit()
	if len(tok) > 0 {
		tokens = append(tokens, tok)
	}
	return tokens
}

func MyStrTokDQ(src string, sep string) []string {
	tokens := []string{}
	MsgWarn("TODO(alf): Not Implemented", "src", src, "sep", sep)

	return tokens
}
