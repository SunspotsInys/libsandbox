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
	rlimit.Cur = 1
	rlimit.Max = 2
	prLimit(proc.Pid, unix.RLIMIT_CPU, &rlimit)
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
	var rlimit unix.Rlimit
	rlimit.Cur = 1024 * 9999
	rlimit.Max = 1024 + 1024
	prLimit(proc.Pid, unix.RLIMIT_AS, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("memory test failed")
	}
	/*
		proc, err = os.StartProcess("test/test", []string{"test"}, &os.ProcAttr{})
		if err != nil {
			panic(err)
		}
		defer proc.Kill()
		prLimit(proc.Pid, unix.RLIMIT_AS, &rlimit)
		if status.Success() {
			t.Fatal("memory sest failed")
		}
	*/
}

/*
func TestMemoryDATALimit(t *testing.T) {
	proc, err := os.StartProcess("test/test", []string{"test"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	defer proc.Kill()
	var rlimit syscall.Rlimit
	rlimit.Cur = 1024
	rlimit.Max = 1024 + 1024
	prLimit(proc.Pid, syscall.RLIMIT_DATA, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("memory test failed")
	}
}
*/
