package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
func TestGen(t *testing.T) {
	cmd := exec.Command("./generate", "-i", "test/input", "-o", "test/output", "-s", "test/main")
	cmd.Run()
	isExist("test/output")
	f, _ := os.Open("test/output")
	b := readFile(f)
	if !bytes.Equal(b, []byte("3\n7\n!-_-\n7\n11\n15\n")) {
		b = bytes.Replace(b, []byte("\n"), []byte("\n\\n"), -1)
		t.Logf("%s != %s", b, []byte("3\n\\n7\n\\n!-_-\n\\n7\n\\n11\n\\n15\n\\n"))
		t.Fatal("test fail")
	}
}
