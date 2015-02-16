package sandbox

import (
	"os"
	"testing"
)

func TestIO(t *testing.T) {
	r := Run("test/open", os.Stdin, os.Stdout, []string{""}, 1000, 2000)
	if r.Status != IOE {
		t.Fatal("IO test failed")
	}
}

func TestTime(t *testing.T) {
	obj := Run("/bin/sleep", os.Stdin, os.Stdout, []string{"5"}, 1000, 20000)
	if obj.Status != TLE {
		t.Log(status[obj.Status])
		t.Log(obj.Time)
		t.Fatal("time exceed test failed.")
	}
}

func TestCPUTime(t *testing.T) {
	obj := Run("test/time", os.Stdin, os.Stdout, []string{""}, 1000, 10000)
	if obj.Status != TLE {
		t.Log(status[obj.Status])
		t.Fatal("time exceed test failed")
	}
}

func TestMemory(t *testing.T) {
	obj := Run("test/memo", os.Stdin, os.Stdout, []string{""}, 1000, 10000)
	if obj.Status != MLE {
		t.Log(status[obj.Status])
		t.Log(obj.Memory)
		t.Fatal("memory exceed test failed")
	}
}
