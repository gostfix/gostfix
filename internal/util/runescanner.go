package util

import (
	"text/scanner"
)

type RuneScanner struct {
	src string
	buf []rune
	Pos int
	Off int
}

func NewRuneScanner(str string) *RuneScanner {
	return &RuneScanner{
		src: str,
		buf: []rune(str),
		Pos: 0,
		Off: 0,
	}
}

func (s *RuneScanner) Next() rune {
	var r rune
	if r = s.Scan(); r != scanner.EOF {
		s.Off++
	}

	return r
}

func (s *RuneScanner) Scan() rune {
	var r rune
	if r = s.Peek(); r != scanner.EOF {
		s.Pos++
	}

	return r
}

func (s *RuneScanner) Peek() rune {
	if s.Pos >= len(s.buf) {
		return scanner.EOF
	}

	return s.buf[s.Pos]
}

func (s *RuneScanner) Push() {
	if s.Pos > 0 && s.Pos > s.Off {
		s.Pos--
	}
}

func (s *RuneScanner) Span(cmp func(rune) bool) {
	for ; s.Pos < len(s.buf) && cmp(s.buf[s.Pos]); s.Pos++ {
		// skip runes until we match one of the stop runes
	}
}

func (s *RuneScanner) Skip(cmp func(rune) bool) {
	s.Span(cmp)
	s.Sync()
}

func (s *RuneScanner) Emit() string {
	str := string(s.buf[s.Off:s.Pos])
	s.Off = s.Pos
	return str
}

func (s *RuneScanner) Sync() {
	s.Off = s.Pos
}

func (s *RuneScanner) String() string {
	return s.src
}

func (s *RuneScanner) Reset() {
	s.Pos = 0
	s.Off = 0
}

func (s *RuneScanner) End(sync bool) {
	s.Pos = len(s.buf)
	if sync {
		s.Off = s.Pos
	}
}

func (s *RuneScanner) PosString() string {
	return string(s.buf[s.Pos:])
}
