package sandbox

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// set time limit for the process
func setTimelimit(pid int, timeLimit int64) error {
	var rlimit unix.Rlimit
	rlimit.Cur = uint64(timeLimit)
	rlimit.Max = uint64(timeLimit)
	return prLimit(pid, unix.RLIMIT_CPU, &rlimit)
}

// set memory limit for the process
func setMemLimit(pid int, memLimit int64) error {
	var rlimit unix.Rlimit
	rlimit.Cur = uint64(memLimit)
	rlimit.Max = uint64(memLimit)
	return prLimit(pid, unix.RLIMIT_AS, &rlimit)

}

// prLimit is the wrapper for the syscall prlimit.
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
