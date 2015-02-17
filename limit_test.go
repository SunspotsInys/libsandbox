package sandbox

import (
	"os"
	"os/exec"
	"syscall"
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
	c := exec.Command("test/test")
	c.Stdout = os.Stdout
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := c.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = setMemLimit(c.Process.Pid, 1024*512)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Wait()

	if err != nil {
		t.Fatal(err)
	}

}
