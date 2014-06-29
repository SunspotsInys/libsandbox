package sandbox

import (
	"syscall"
	"unsafe"
)

func setLimit(pid int) error {

	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRLIMIT64, uintptr(pid), syscall.RLIMIT_CPU, uintptr(unsafe.Pointer(&syscall.Rlimit{
		Cur: 1,
		Max: 1 + 1,
	})), 0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
		return err
	} else {
		return nil
	}
}
