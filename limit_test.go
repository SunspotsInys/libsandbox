package sandbox

import (
	"os"
	"os/exec"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestCPULimit(t *testing.T) {
	c := exec.Command("test/main")

	c.Start()

	err := setTimelimit(c.Process.Pid, 3)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Wait()
	if err == nil {
		t.Fatal("CPU time limit test failed")
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
	// sometime
	err = setMemLimit(c.Process.Pid, 1024*512)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Wait()

	// soemtimes got exit status 127,sometimes got segment fault,not know why.

	if err == nil {
		t.Fatal("memory limit test failed")
	}

}
