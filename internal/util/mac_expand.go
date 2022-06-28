package util

const MAC_EXP_FLAG_NONE int = (0)
const MAC_EXP_FLAG_RECURSE int = (1 << 0)
const MAC_EXP_FLAG_APPEND int = (1 << 1)
const MAC_EXP_FLAG_SCAN int = (1 << 2)
const MAC_EXP_FLAG_PRINTABLE int = (1 << 3)

type MacExpandContext struct {
}

func MacExpand(pattern string, flags int, filter *string, lookup func(string, int, interface{}) string, context interface{}) (string, error) {
	return pattern, nil
}
