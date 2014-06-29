package sandbox

import (
	"testing"
)

func TestCPUTime(t *testing.T) {
	Run("test/main", []string{"main"})
}

func TestMemoryTime(t *testing.T) {
	Run("test/memo", []string{"memo"})
}
