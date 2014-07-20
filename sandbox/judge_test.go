package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestJudgeAPlusB(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "judge/src/A1/a+b.c", "judge/binary/A1/a+b", "judge/src/A1/input", "judge/src/A1/output")
	out, _ := cmd.CombinedOutput()
	t.Logf("%s", out)
	if fmt.Sprintf("%s", out) != "AC" {
		t.Fatal("wrong answer")
	}
}
func TestTimeLimit(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "judge/src/A2/main.c", "judge/binary/A2/main")
	out, _ := cmd.CombinedOutput()
	t.Logf("%s", out)
	if fmt.Sprintf("%s", out) != "TLE" {
		t.Fatal("wrong answer")
	}
}
