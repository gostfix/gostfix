package util

import (
	"os"
	"strings"

	"github.com/gostfix/gostfix/internal/ascii"
)

func DirForest(path string, depth int) string {
	if path == "" {
		MsgPanic("empty path")
	}
	if depth < 1 {
		MsgPanic("depth to low", "depth", depth)
	}

	sb := strings.Builder{}
	sb.Grow(depth * 2)

	for cp, n := []rune(path), 0; n < depth; n++ {
		ch := '_'
		if n < len(cp) {
			ch = cp[n]
			if !ascii.IsPrint(ch) || ch == '.' || ch == '/' || ch == os.PathSeparator {
				MsgPanic("invalid pathname", "path", path)
			}
		}

		sb.WriteRune(ch)
		sb.WriteRune(os.PathSeparator)
	}

	if MsgVerbose > 1 {
		MsgInfof("%s -> %s", path, sb.String())
	}

	return sb.String()
}
