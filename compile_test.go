package sandbox

import (
	"testing"

	"github.com/ggaaooppeenngg/util"
)

func TestCompile(t *testing.T) {
	err := compile("test/main.go", "test/main", GO)
	if err != nil || !util.IsExist("test/main") {
		t.Log(err)
		t.Fatal("Go compile failed")
	}

	err = compile("test/memo.go", "test/memo", GO)
	if err != nil || !util.IsExist("test/test") {
		t.Log(err)
		t.Fatal("Go compile failed")
	}

	err = compile("test/test.c", "test/test", C)
	if err != nil || !util.IsExist("test/test") {
		t.Log(err)
		t.Fatal("C compile failed")
	}

	err = compile("test/time.c", "test/time", C)
	if err != nil || !util.IsExist("test/time") {
		t.Log(err)
		t.Fatal("C compile failed")
	}

	err = compile("test/memo_limit.cpp", "test/memo_limit", CPP)
	if err != nil || !util.IsExist("test/time") {
		t.Log(err)
		t.Fatal("CPP compile failed")
	}
}
