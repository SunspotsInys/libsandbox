package sandbox

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
#include <unistd.h>
*/
import "C"

var (
	sc_clk_tck int64
	frequency  time.Duration
)

func init() {
	//timer click number per second
	sc_clk_tck = int64(C.sysconf(C._SC_CLK_TCK))
	frequency = time.Second / time.Duration(sc_clk_tck)
}

// get process virtual memory usage
func virtualMemory(pid int) int64 {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(stat)
	if err != nil {
		panic(err)
	}
	//virtual memory size is 23nd paramater in the stat file,in bytes
	vmSize, err := strconv.ParseInt(strings.Split(string(bs), " ")[22], 10, 64)

	if err != nil {
		panic(err)
	}
	return vmSize
}

// return process running time from the start
func realTime(pid int) int64 {
	upTimeFile, err := os.Open("/proc/uptime")
	if err != nil {
		panic(err)
	}
	defer upTimeFile.Close()
	bs, err := ioutil.ReadAll(upTimeFile)
	if err != nil {
		panic(err)
	}
	//uptime is first paramater in uptime file
	upTime, err := strconv.ParseFloat(strings.Split(string(bs), " ")[0], 64)
	//reserve milliensecond for further usage
	upTimeM := int64(upTime * 1000)
	if err != nil {
		panic(err)
	}
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		panic(err)
	}
	defer stat.Close()
	bs, err = ioutil.ReadAll(stat)
	if err != nil {
		panic(err)
	}
	//startTime is 22nd paramater in the stat file
	startTime, err := strconv.ParseInt(strings.Split(string(bs), " ")[21], 10, 64)
	//reserve milliensecond for further usage
	startTimeM := int64(float64(startTime) * 1000 / float64(sc_clk_tck))
	if err != nil {
		panic(err)
	}
	return upTimeM - startTimeM
}
