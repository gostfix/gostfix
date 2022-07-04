package util

import "github.com/gostfix/gostfix/internal/ascii"

func AllDig(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, d := range str {
		if !ascii.IsDigit(d) {
			return false
		}
	}
	return true
}

func AllAlnum(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ascii.IsAlnum(r) {
			return false
		}
	}
	return true
}
