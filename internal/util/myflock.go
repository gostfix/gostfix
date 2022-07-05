package util

import "syscall"

/*
 * Lock styles.
 */
const (
	MYFLOCK_STYLE_FLOCK int = 1
	MYFLOCK_STYLE_FCNTL int = 2
)

/*
 * Lock request types.
 */
const (
	MYFLOCK_OP_NONE      int = 0
	MYFLOCK_OP_SHARED    int = 1 << 1
	MYFLOCK_OP_EXCLUSIVE int = 1 << 2
	MYFLOCK_OP_NOWAIT    int = 1 << 3
	MYFLOCK_OP_BITS      int = (MYFLOCK_OP_SHARED | MYFLOCK_OP_EXCLUSIVE | MYFLOCK_OP_NOWAIT)
)

func MyFlock(fd int, lock_style int, operation int) int {
	var status int = 0

	/*
	 * Sanity check.
	 */
	if operation&MYFLOCK_OP_BITS != operation {
		MsgPanic("improper operation type", "functon", "myflock", "operation", operation)
	}

	switch lock_style {
	case MYFLOCK_STYLE_FLOCK:
		var lock_ops = []int{
			syscall.LOCK_UN, syscall.LOCK_SH, syscall.LOCK_EX, -1, -1,
			syscall.LOCK_SH | syscall.LOCK_NB, syscall.LOCK_EX | syscall.LOCK_NB, -1,
		}
		if ret := syscall.Flock(fd, lock_ops[operation]); ret != nil {
			status = -1
		}
	case MYFLOCK_STYLE_FCNTL:
		var lock = syscall.Flock_t{}
		var lock_ops = []int16{syscall.F_UNLCK, syscall.F_RDLCK, syscall.F_WRLCK}
		var request = syscall.F_SETLKW

		lock.Type = lock_ops[operation & ^MYFLOCK_OP_NOWAIT]
		if operation&MYFLOCK_OP_NOWAIT == MYFLOCK_OP_NOWAIT {
			request = syscall.F_SETLK
		}

		if ret := syscall.FcntlFlock(uintptr(fd), request, &lock); ret != nil {
			status = -1
		}
	default:
		MsgPanic("unsupported lock style", "function", "myflock", "lock_style", lock_style)
	}

	return status
}
