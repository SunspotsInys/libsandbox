package sandbox

import (
	"os/exec"
)

const (
	C   uint64 = iota // C language
	CPP               // C Plus Plus langua
	GO                // Go language
)

// default comiple options
func compile(src string, des string, lan uint64) error {
	var cmd = new(exec.Cmd)
	switch lan {
	case C:
		cmd = exec.Command("gcc", "-o", des, src, "-lm") //-lm for gcc math link option
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
