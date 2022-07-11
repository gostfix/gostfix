package util

import (
	"os"
	"strings"
)

func UpdateEnv(preserve_list []string) {
	for _, env := range preserve_list {
		if strings.ContainsRune(env, '=') {
			nv := strings.SplitN(env, "=", 2)
			os.Setenv(nv[0], nv[1])
		}
	}
}
