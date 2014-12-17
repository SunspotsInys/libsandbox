package sandbox

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func prLimit(pid int, limit uintptr, rlimit *unix.Rlimit) error {
	_, _, errno := unix.RawSyscall6(unix.SYS_PRLIMIT64,
		uintptr(pid),
		limit,
		uintptr(unsafe.Pointer(rlimit)),
		0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
		return err
	} else {
		return nil
	}
}
