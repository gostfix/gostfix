package util

import "os"

// SafeGetenv - read environment variable with guard
func SafeGetenv(name string) string {
	if Unsafe() {
		return ""
	}
	return os.Getenv(name)
}
