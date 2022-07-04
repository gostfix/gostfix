package ascii

func ToLower(r rune) rune {
	lower := r
	if IsUpper(r) {
		lower -= 0x20
	}
	return lower
}

func ToUpper(r rune) rune {
	upper := r
	if IsLower(r) {
		upper += 0x20
	}
	return upper
}
