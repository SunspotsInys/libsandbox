package sandbox

import (
	"testing"
)

func TestTime(t *testing.T) {
	Run("/bin/sleep", []string{"sleep", "5"})
}

/*
func TestCPUTime(t *testing.T) {
	Run("test/main", []string{"main"})
}

/*
func TestMemoryTime(t *testing.T) {
	Run("test/memo", []string{"memo"})
}
*/
