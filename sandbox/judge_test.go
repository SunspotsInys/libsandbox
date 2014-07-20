package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestJudge(t *testing.T) {
	cmd := exec.Command("sandbox", "--lang=c", "../test/judge/a+b.c", "../test/judge/a+b", "..test/judge/input", "../test/judge/output")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		println(err)
	}
	t.Log(out.Bytes())
	if fmt.Sprintf("%s", out.Bytes()) != "AC" {
		t.Log(err)
		t.Fatal("wrong answer")
	}
}
