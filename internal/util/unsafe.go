package util

import "os"

// Unsafe
// Can we trust user-provided environment, working directory, etc.
func Unsafe() bool {
	// the super-user is trusted
	if os.Getuid() == 0 && os.Geteuid() == 0 {
		return false
	}

	// Danger: don't trust inherited process attributes, and don't leak
	// privileged info that the parent has no access to.
	return (os.Getuid() != os.Geteuid() || os.Getgid() != os.Getegid())
}
