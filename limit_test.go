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
	rlimit.Cur = 1024
	rlimit.Max = 1024 + 1024
	prLimit(proc.Pid, syscall.RLIMIT_AS, &rlimit)
	status, err := proc.Wait()
	if status.Success() {
		t.Fatal("memory test failed")
	}
	t.Log(status.String())
	/*
		proc, err = os.StartProcess("test/test", []string{"test"}, &os.ProcAttr{})
		if err != nil {
			panic(err)
		}
		prLimit(proc.Pid, syscall.RLIMIT_AS, &rlimit)
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
