package global

import (
	"fmt"
	"math/bits"
	"strings"
	"text/scanner"

	"github.com/gostfix/gostfix/internal/ascii"
	"github.com/gostfix/gostfix/internal/util"
)

const safe_chars string = "0123456789BCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz"
const SAFE_MAX_BASE uint = uint(len(safe_chars))
const SAFE_MIN_BASE uint = uint(2)

var char_map [256]uint8 = [256]uint8{}

func init() {
	var ch uint
	for ch = 0; ch < uint(len(char_map)); ch++ {
		char_map[ch] = uint8(SAFE_MAX_BASE)
	}
	for ch = 0; ch < SAFE_MAX_BASE; ch++ {
		char_map[safe_chars[ch]] = uint8(ch)
	}
}

func SafeUlToStr(val uint, base uint, padlen int, padchar rune) string {
	if base < SAFE_MIN_BASE || base > SAFE_MAX_BASE {
		util.MsgPanic("bad base", "base", base)
	}
	result := strings.Builder{}

	// accumulate the result, backwards
	for val > 0 {
		result.WriteByte(safe_chars[val%base])
		val /= base
	}
	for result.Len() < padlen {
		result.WriteRune(padchar)
	}

	// then, reverse the result
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func SafeStrToUl(str string, end *string, base uint) (uint, error) {
	var err error = nil

	if base < SAFE_MIN_BASE || base > SAFE_MAX_BASE {
		util.MsgPanic("bad base", "base", base)
	}

	var sum uint = 0
	var div_limit uint = (1<<bits.UintSize - 1) / base
	var mod_limit uint = (1<<bits.UintSize - 1) % base

	scan := util.NewRuneScanner(str)
	scan.Skip(ascii.IsSpace)

	for ch := scan.Scan(); ch != scanner.EOF; ch = scan.Scan() {
		if int(ch) < len(char_map) && uint(char_map[ch]) >= base {
			scan.Push()
			break
		}
		char_val := uint(char_map[ch])

		if sum > div_limit || (sum == div_limit && char_val > mod_limit) {
			// set sum to maximum uint value
			sum = (1<<bits.UintSize - 1)
			err = fmt.Errorf("out of range")
			// skip remaining valid numbers in base, per strtoul() spec
			scan.Skip(func(r rune) bool { return uint(char_map[r]) < base })
			break
		}
		sum = sum*base + char_val
	}
	if scan.Pos == 0 {
		err = fmt.Errorf("invalid value")
	} else if end != nil {
		*end = str[scan.Pos:]
	}
	return sum, err
}
