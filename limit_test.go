package sandbox

import (
	"os"
	"syscall"
	"testing"
)

func TestCPULimit(t *testing.T) {
	proc, err := os.StartProcess("test/main", []string{"main"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	defer proc.Kill()
	var rlimit syscall.Rlimit
	rlimit.Cur = 1
	rlimit.Max = 2
	prLimit(proc.Pid, syscall.RLIMIT_CPU, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("cpu limit test failed")
	}
}

func TestMemoryLimit(t *testing.T) {
	proc, err := os.StartProcess("test/memo", []string{"memo"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	defer proc.Kill()
	var rlimit syscall.Rlimit
	rlimit.Cur = 10
	rlimit.Max = 10 + 1024
	prLimit(proc.Pid, syscall.RLIMIT_DATA, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("memory test failed")
	}
}
