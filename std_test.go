package libsandbox

import (
	"testing"
)

func TestOutOfRunningTime(t *testing.T) {
	out, err := StdSandbox{
		Bin:         "sh",
		Args:        []string{"-c", "sleep 3"},
		Input:       nil,
		TimeLimit:   1000,
		MemoryLimit: 1,
	}.Run()
	if err == nil {
		t.Fatalf("no error get out '%s', want %s\n", out, OutOfTimeError)
	}
	if err != OutOfTimeError {
		t.Fatalf("unexpecged error %s, want %s\n", err, OutOfTimeError)
	}
}
func TestOutOfMemory(t *testing.T) {
	out, err := StdSandbox{
		Bin:         "sh",
		Args:        []string{"-c", "sleep 3"},
		Input:       nil,
		TimeLimit:   1000,
		MemoryLimit: 100,
	}.Run()
	if err == nil {
		t.Fatalf("no error get out '%s', want %s\n", out, OutOfMemoryError)
	}
	if err != OutOfMemoryError {
		t.Fatalf("unexpecged error %s, want %s\n", err, OutOfMemoryError)
	}

}
