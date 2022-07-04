package util

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"unicode"

	"github.com/gostfix/gostfix/internal/ascii"
)

var runePool = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz")
var runeTable = unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: '0', Hi: '9', Stride: 1},
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: '_', Hi: '_', Stride: 1},
		{Lo: 'a', Hi: 'z', Stride: 1},
	},
}
var runeString = string(runePool)

var table = []struct {
	input string
}{
	{input: RandStringRunes(1)},
	{input: RandStringRunes(5)},
	{input: RandStringRunes(10)},
	{input: RandStringRunes(50)},
	{input: RandStringRunes(100)},
	{input: RandStringRunes(500)},
	{input: RandStringRunes(1000)},
	{input: RandStringRunes(2500)},
	{input: RandStringRunes(5000)},
	{input: RandStringRunes(10000)},
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runePool[rand.Intn(len(runePool))]
	}
	return string(b)
}

func (s *RuneScanner) SpanTable(rt *unicode.RangeTable) {
	for ; s.Pos < len(s.buf) && unicode.In(s.buf[s.Pos], rt); s.Pos++ {
		// skip runes until we match one of the stop runes
	}
}

func (s *RuneScanner) SpanString(str string) {
	for ; s.Pos < len(s.buf) && strings.ContainsRune(str, s.buf[s.Pos]); s.Pos++ {
		// skip runes until we match one of the stop runes
	}
}

func BenchmarkScanIfTable(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", len(v.input)), func(b *testing.B) {
			rs := NewRuneScanner(v.input)
			for i := 0; i < b.N; i++ {
				rs.Reset()
				rs.SpanTable(&runeTable)
			}
		})
	}
}

func BenchmarkScanIfString(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", len(v.input)), func(b *testing.B) {
			rs := NewRuneScanner(v.input)
			for i := 0; i < b.N; i++ {
				rs.Reset()
				rs.SpanString(runeString)
			}
		})
	}
}

func BenchmarkScanIfCmpLUT(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", len(v.input)), func(b *testing.B) {
			rs := NewRuneScanner(v.input)
			for i := 0; i < b.N; i++ {
				rs.Reset()
				rs.Span(func(r rune) bool { return ascii.IsAlnum(r) || r == '_' })
			}
		})
	}
}

func BenchmarkScanIfCmpRune(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", len(v.input)), func(b *testing.B) {
			rs := NewRuneScanner(v.input)
			for i := 0; i < b.N; i++ {
				rs.Reset()
				rs.Span(func(r rune) bool {
					return (r <= '9' && r >= '0') || (r <= 'Z' && r >= 'A') || (r <= 'z' || r >= 'a') || r == '_'
				})
			}
		})
	}
}

func BenchmarkScanIfSwitch(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", len(v.input)), func(b *testing.B) {
			rs := NewRuneScanner(v.input)
			for i := 0; i < b.N; i++ {
				rs.Reset()
				rs.Span(func(r rune) bool {
					switch {
					case r <= '9' && r >= '0':
						return true
					case r <= 'Z' && r >= 'A':
						return true
					case r <= 'z' && r >= 'a':
						return true
					case r == '_':
						return true
					default:
						return false
					}
				})
			}
		})
	}
}
