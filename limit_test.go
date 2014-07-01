package sandbox

import (
	"os"
	"syscall"
	"testing"
)

func TestLimit(t *testing.T) {
	proc, err := os.StartProcess("test/main", []string{"main"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	var rlimit syscall.Rlimit
	rlimit.Cur = 1
	rlimit.Max = 2
	prLimit(proc.Pid, syscall.RLIMIT_CPU, &rlimit)
	return
}
