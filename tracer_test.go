package sandbox

import (
	"testing"
)

/*
func TestTime(t *testing.T) {
	obj := Run("/bin/sleep", []string{"sleep", "5"}, 1000, 10000)
	//t.Log(obj.Status)
	if obj.Status != TLE {
		t.Fatal("time exceed test failed.")
	}
}
*/

func TestCPUTime(t *testing.T) {
	obj := Run("test/time", []string{"time"}, 1000, 10000)
	//t.Log(obj.Status)
	if obj.Status != TLE {
		t.Fatal("time exceed test failed")
	}
}

func TestMemory(t *testing.T) {
	obj := Run("test/memo", []string{"memo"}, 1000, 10000)
	//t.Log(obj.Status)
	if obj.Status != MLE {
		t.Fatal("memory exceed test failed")
	}
}
