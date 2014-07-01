package sandbox

import (
	"os"
	"testing"
	"time"
)

func TestRealTime(t *testing.T) {
	proc, err := os.StartProcess("/bin/sleep", []string{"sleep", "5"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)
	realTime := realTime(proc.Pid)
	if realTime < 1000 {
		t.Fatal("real Time measure error")
	}
}

func TestVmSize(t *testing.T) {
	proc, err := os.StartProcess("test/main", []string{"test"}, &os.ProcAttr{})
	if err != nil {
		panic(err)
	}
	time.Sleep(2)
	vs := virtualMemory(proc.Pid)
	if vs < 2048*100+40*1024 {
		t.Fatal("real vmSize error")
	}

}
