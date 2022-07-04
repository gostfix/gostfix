package ascii

func IsDigit(r rune) bool {
	return r <= MaxASCII && asciitable[r]&digit == digit
}

func IsXdigit(r rune) bool {
	return r <= MaxASCII && asciitable[r]&xdigit == xdigit
}
