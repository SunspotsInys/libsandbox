package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestJudgeAPlusB(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "-c", "judge/src/A1/a+b.c", "judge/binary/A1/a+b", "judge/src/A1/input", "judge/src/A1/output")
	out, _ := cmd.CombinedOutput()
	//t.Logf("%s", out)
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	if status != "AC" {
		t.Fatal("wrong answer")
	}

}
func TestJudgeWithoutCompiling(t *testing.T) {
	//run without compiling
	cmd := exec.Command("sandbox", "--lang=c", "judge/binary/A1/a+b", "judge/src/A1/input", "judge/src/A1/output")
	out, _ := cmd.CombinedOutput()
	t.Logf("%s", out)
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	time := results[1]
	memory := results[2]
	if status != "AC" || time != "0" || memory != "0" {
		t.Fatal("wrong answer")
	}
}

func TestTimeLimit(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "-c", "judge/src/A2/main.c", "judge/binary/A2/main")
	out, _ := cmd.CombinedOutput()
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	time := results[1]
	memory := results[2]
	t.Logf("%s", out)
	if status != "TLE" || time == "0" || memory == "0" {
		t.Fatal("wrong answer")
	}
}
