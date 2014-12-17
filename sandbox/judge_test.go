package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestPresentationError(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c",
		"-c", "-s", "judge/src/A7/main.c", "-b",
		"judge/binary/A7/main", "-i", "judge/src/A7/input",
		"-o", "judge/src/A7/output")
	out, _ := cmd.CombinedOutput()
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	if status != "FE" {
		t.Logf("%s", status)
		t.Logf("%s", out)
		t.Fatal("Test Presentation Error Failed")
	}
}

func TestSegmentfault(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c",
		"-c", "-s", "judge/src/A6/main.c", "-b",
		"judge/binary/A6/main", "-i", "judge/src/A6/input",
		"-o", "judge/src/A6/output")
	out, _ := cmd.CombinedOutput()
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	if status != "RE" {
		t.Logf("%s", out)
		t.Fatal("wrong answer")
	}
}

func TestJudgeAPlusB(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c",
		"-c", "-s", "judge/src/A1/a+b.c", "-b",
		"judge/binary/A1/a+b", "-i", "judge/src/A1/input",
		"-o", "judge/src/A1/output")
	out, _ := cmd.CombinedOutput()
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	if status != "AC" {
		t.Logf("%s", out)
		t.Fatal("wrong answer")
	}

}
func TestJudgeWithoutCompiling(t *testing.T) {
	//run without compiling
	cmd := exec.Command("sandbox", "--lang=c",
		"-b", "judge/binary/A1/a+b", "-i", "judge/src/A1/input",
		"-o", "judge/src/A1/output")
	out, _ := cmd.CombinedOutput()
	t.Logf("%s", out)
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	if status != "AC" {
		t.Fatal("wrong answer")
	}
}
func TestNTimesAPlusB(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "-c", "-s",
		"judge/src/A3/main.c", "-b", "judge/binary/A3/tmp",
		"-i", "judge/src/A3/input",
		"-o", "judge/src/A3/output")
	out, _ := cmd.CombinedOutput()
	t.Logf("%s", out)
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	nth := results[3]
	if status != "WA" || nth != "2" {
		t.Fatal("wrong answer")
	}

}

func TestTimeLimit(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "-c", "-s",
		"judge/src/A2/main.c", "-b",
		"judge/binary/A2/main")
	out, _ := cmd.CombinedOutput()
	result := fmt.Sprintf("%s", out)
	results := strings.Split(result, ":")
	status := results[0]
	time := results[1]
	memory := results[2]
	if status != "TL" || time == "0" || memory == "0" {
		t.Logf("%s", out)
		t.Fatal("wrong answer")
	}
	cmd = exec.Command("sandbox", "--lang=c", "-c", "-s",
		"judge/src/A8/main.c", "-b",
		"judge/binary/A8/main", "-i", "judge/src/A8/input",
		"-o", "judge/src/A8/output")
	out, _ = cmd.CombinedOutput()
	result = fmt.Sprintf("%s", out)
	results = strings.Split(result, ":")
	status = results[0]
	time = results[1]
	memory = results[2]
	if status != "TL" || time == "0" {
		t.Logf("%s", out)
		t.Fatal("wrong answer")
	}
}
