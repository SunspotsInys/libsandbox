package sandbox

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/ggaaooppeenngg/ptrace.go/ptrace"

	"golang.org/x/sys/unix"
)

const (
	AC uint64 = iota
	PE
	TLE
	MLE
	WA
	RE
	OLE
	CE
	SE
)

const (
	KB = 1024
)

//for debug
var status = map[uint64]string{
	AC:  "Accept",
	PE:  "Presentation Error",
	TLE: "Time Limit Error",
	MLE: "Memory Limit Error",
	WA:  "Wrong Answer",
	RE:  "Runtime Error",
	OLE: "Output Limit Error",
	CE:  "Compile Error",
	SE:  "Segmentfault Error",
}

type RunningObject struct {
	Proc        *os.Process
	TimeLimit   int64
	MemoryLimit int64
	Memory      int64 //KB
	Time        int64 //MS
	Status      uint64
}

func (r *RunningObject) RunTick() {
	ticker := time.NewTicker(frequency)
	//send alarm signal with time tick frequency
	for _ = range ticker.C {
		r.Proc.Signal(os.Signal(unix.SIGALRM))
	}
}

//wrap compile function
func Complie(src string, des string, lan uint64) error {
	return compile(src, des, lan)
}

func Run(bin string, reader io.Reader, writer io.Writer,
	args []string, timeLimit int64, memoryLimit int64) *RunningObject {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var rusage unix.Rusage
	var runningObject RunningObject
	runningObject.TimeLimit = timeLimit
	runningObject.MemoryLimit = memoryLimit
	cmd := exec.Command(bin, args...)
	cmd.Stdin = reader
	cmd.Stderr = writer
	cmd.Stdout = writer
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	proc := cmd.Process
	tracer, err := ptrace.Attach(proc)
	if err != nil {
		panic(err)
	}
	runningObject.Proc = proc
	go runningObject.RunTick()
	var rlimit unix.Rlimit
	rlimit.Cur = uint64(timeLimit)
	rlimit.Max = uint64(timeLimit)
	err = prLimit(proc.Pid, unix.RLIMIT_CPU, &rlimit)
	if err != nil {
		fmt.Println(err)
		return &runningObject
	}
	for {
		status := unix.WaitStatus(0)
		//ptrace stopped
		_, err := unix.Wait4(proc.Pid,
			&status,
			unix.WSTOPPED,
			&rusage)
		if err != nil {
			panic(err)
		}
		if status.Exited() {
			return &runningObject
		}
		if status.CoreDump() {
			fmt.Println("CoreDump")
			return &runningObject
		}
		if status.Continued() {
			fmt.Println("Continued")
			return &runningObject
		}
		if status.Signaled() {
			fmt.Println("signal")
			return &runningObject
		}
		if status.Stopped() &&
			status.StopSignal() != unix.SIGTRAP {
			switch status.StopSignal() {
			case unix.SIGALRM:
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / KB
				runningObject.Time = rusage.Utime.Sec*1000 +
					rusage.Utime.Usec/1000
				if runningObject.Time >
					runningObject.TimeLimit {
					runningObject.Status = TLE
					err := runningObject.Proc.Kill()
					if err != nil {
						panic(err)
					}
					return &runningObject
				}
				realTime := realTime(runningObject.Proc.Pid)
				if realTime > runningObject.TimeLimit {
					runningObject.Status = TLE
					runningObject.Time = realTime
					err := runningObject.Proc.Kill()
					if err != nil {
						log.Println(err)
					}
					return &runningObject
				}
				if vs/KB > runningObject.MemoryLimit {
					runningObject.Status = MLE
					err := runningObject.Proc.Kill()
					if err != nil {
						log.Println(err)
					}
					return &runningObject
				}
			case unix.SIGXCPU:
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / KB
				runningObject.Time = rusage.Utime.Sec*1000 +
					rusage.Utime.Usec/1000
				runningObject.Status = TLE
				err := runningObject.Proc.Kill()
				if err != nil {
					log.Println(err)
				}
				return &runningObject
			case unix.SIGSEGV:
				//if segmentfault
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / KB
				runningObject.Time = rusage.Utime.Sec*1000 +
					rusage.Utime.Usec/1000
				runningObject.Status = RE
				err := runningObject.Proc.Kill()
				if err != nil {
					log.Println(err)
				}
				return &runningObject
			default:
			}
		}
		//0表示不发出信号
		tracer.Syscall(syscall.Signal(0))
	}
}
