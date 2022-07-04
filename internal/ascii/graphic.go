package ascii

func IsAscii(r rune) bool {
	return r <= MaxASCII
}

func IsAlpha(r rune) bool {
	return r <= MaxASCII && asciitable[r]&alpha == alpha
}

func IsAlnum(r rune) bool {
	return r <= MaxASCII && asciitable[r]&(alpha|digit) != 0
}

func IsSpace(r rune) bool {
	return r <= MaxASCII && asciitable[r]&space == space
}

func IsBlank(r rune) bool {
	return r <= MaxASCII && asciitable[r]&blank == blank
}

func IsCntrl(r rune) bool {
	return r <= MaxASCII && asciitable[r]&cntrl == cntrl
}

func IsPrint(r rune) bool {
	return r <= MaxASCII && asciitable[r]&print == print
}

func IsGraph(r rune) bool {
	return r <= MaxASCII && asciitable[r]&graph == graph
}

func IsUpper(r rune) bool {
	return r <= MaxASCII && asciitable[r]&upper == upper
}

func IsLower(r rune) bool {
	return r <= MaxASCII && asciitable[r]&lower == lower
}

func IsPunct(r rune) bool {
	return r <= MaxASCII && asciitable[r]&punct == punct
}
