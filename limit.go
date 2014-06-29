package sandbox

import (
	"syscall"
	"unsafe"
)

func prLimit(pid int, limit uintptr, rlimit *syscall.Rlimit) error {
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRLIMIT64, uintptr(pid), limit, uintptr(unsafe.Pointer(rlimit)), 0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
		return err
	} else {
		return nil
	}
}
