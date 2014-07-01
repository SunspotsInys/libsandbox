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
	t.Log(realTime)
	if realTime != 1 {
		t.Fatal("real Time measure error")
	}
}

func TestVmSize(t *testing.T) {
	/*
		proc, err := os.StartProcess("/bin/sleep", []string{"sleep", "5"}, &os.ProcAttr{})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 1)
		vs := virtualMemory(proc.Pid)
		t.Log(vs)
		if vs > 10000*1024 {
			t.Fatal("real vmSize error")
		}
	*/
}
