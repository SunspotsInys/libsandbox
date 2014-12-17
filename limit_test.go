package sandbox

import (
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestCPULimit(t *testing.T) {
	proc, err := os.StartProcess("test/main", []string{"main"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	defer proc.Kill()
	var rlimit unix.Rlimit
	rlimit.Cur = 1000
	rlimit.Max = 1000
	prLimit(proc.Pid, unix.RLIMIT_CPU, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("cpu limit test failed")
	}
}

func TestMemoryLimit(t *testing.T) {
	proc, err := os.StartProcess("test/test", []string{"test"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	defer proc.Kill()
	var rlimit unix.Rlimit
	rlimit.Cur = 1024 * 512
	rlimit.Max = 1024 * 512
	prLimit(proc.Pid, unix.RLIMIT_AS, &rlimit)
	status, err := proc.Wait()
	if err == nil || status.Success() {
		if err == nil {
			t.Log(err)
		}
		t.Fatal("memory sest failed")
	}
}
