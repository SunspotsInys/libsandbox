package sandbox

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
#include <unistd.h>
*/
import "C"

var (
	sc_clk_tck int64
	//sc_page_size int64
)

func init() {
	sc_clk_tck = int64(C.sysconf(C._SC_CLK_TCK))
	//sc_page_size = int64(C.sysconf(C._SC_PAGE_SIZE))
}

func virtualMemory(pid int) int64 {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(stat)
	if err != nil {
		panic(err)
	}
	//virtual memory size is 23nd paramater in the stat file
	vmSize, err := strconv.ParseInt(strings.Split(string(bs), " ")[22], 10, 64)
	if err != nil {
		panic(err)
	}
	return vmSize
}

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
