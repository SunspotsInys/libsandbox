package sandbox

import (
	"testing"

	"github.com/ggaaooppeenngg/util"
)

func TestCompile(t *testing.T) {
	if !util.IsExist("test/main") {
		util.CreateFile("test/main")
	}
	err := compile("test/main.go", "test/main", GO)
	if err != nil || !util.IsExist("test/main") {
		t.Log(err)
		t.Fatal("golang compile failed")
	}
	compile("test/memo.go", "test/memo", GO)
	if !util.IsExist("test/test") {
		util.CreateFile("test/test")
	}
	err = compile("test/test.c", "test/test", C)
	if err != nil || !util.IsExist("test/test") {
		t.Log(err)
		t.Fatal("c compile failed")
	}
	compile("test/time.c", "test/time", C)
}
