package sandbox

import (
	"os/exec"
)

const (
	C uint64 = iota
	CPP
	GO
)

func compile(src string, des string, lan uint64) error {
	var cmd = new(exec.Cmd)
	switch lan {
	case C:
		cmd = exec.Command("gcc", "-o", des, src)
	case CPP:
		cmd = exec.Command("g++", "-o", des, src)
	case GO:
		cmd = exec.Command("go", "build", "-o", des, src)
	}
	if err := cmd.Run(); err != nil {
		return err
	} else {
		return nil
	}

}
