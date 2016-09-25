package sandbox

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRealTime(t *testing.T) {
	proc, err := os.StartProcess("/bin/sleep", []string{"sleep", "5"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)
	realTime := RunningTime(proc.Pid)
	if realTime < 1000 {
		t.Fatal("real Time measure error")
	}

}

func TestVmSize(t *testing.T) {
	cmd := exec.Command("test/memo")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	var pid = cmd.Process.Pid
	time.Sleep(time.Second)
	vs := VirtualMemory(pid)
	if vs < 10000*1024 {
		t.Fatalf("current virtual memory %d KB is smaller than 10000KB", vs/1024)
	}
}
