package sandbox

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/ggaaooppeenngg/ptrace.go/ptrace"
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

type RunningObject struct {
	Proc        *os.Process
	TimeLimit   int64
	MemoryLimit int64
	Memory      int64 //KB
	Time        int64 //MS
	Status      uint64
}

func (r *RunningObject) RunTick(dur time.Duration) {
	ticker := time.NewTicker(dur)
	for _ = range ticker.C {
		r.Proc.Signal(os.Signal(syscall.SIGALRM))
	}
}

func Complie(src string, des string, lan uint64) error {
	return compile(src, des, lan)
}

func Run(bin string, reader io.Reader, writer io.Writer, args []string, timeLimit int64, memoryLimit int64) *RunningObject {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var rusage syscall.Rusage
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
	go runningObject.RunTick(time.Nanosecond)
	var rlimit syscall.Rlimit
	rlimit.Cur = uint64(timeLimit)
	rlimit.Max = uint64(timeLimit)
	err = prLimit(proc.Pid, syscall.RLIMIT_CPU, &rlimit)
	if err != nil {
		fmt.Println(err)
		return &runningObject
	}
	/*
		get "no such process" error when add AS limit
		rlimit.Cur = uint64(memoryLimit) * 1024
		rlimit.Max = uint64(memoryLimit) * 1024
		err = prLimit(proc.Pid, syscall.RLIMIT_AS, &rlimit)
		if err != nil {
			fmt.Println(err)
			return &runningObject
		}

		/*
		err = prLimit(proc.Pid, syscall.RLIMIT_DATA, &rlimit)
		if err != nil {
			fmt.Println(err)
			return &runningObject
		}
		err = prLimit(proc.Pid, syscall.RLIMIT_STACK, &rlimit)
		if err != nil {
			fmt.Println(err)
			return &runningObject
		}
	*/
	for {
		status := syscall.WaitStatus(0)
		_, err := syscall.Wait4(proc.Pid, &status, syscall.WSTOPPED, &rusage)
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
		if status.Stopped() && status.StopSignal() != syscall.SIGTRAP {
			switch status.StopSignal() {
			case syscall.SIGALRM:
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / 1000
				runningObject.Time = rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000
				if runningObject.Time > runningObject.TimeLimit {
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
						panic(err)
					}
					return &runningObject
				}
				if vs/1000 > runningObject.MemoryLimit {
					runningObject.Status = MLE
					err := runningObject.Proc.Kill()
					if err != nil {
						panic(err)
					}
					return &runningObject
				}
			case syscall.SIGXCPU:
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / 1000
				runningObject.Time = rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000
				runningObject.Status = TLE
				err := runningObject.Proc.Kill()
				if err != nil {
					panic(err)
				}
				return &runningObject
			case syscall.SIGSEGV:
				vs := virtualMemory(runningObject.Proc.Pid)
				runningObject.Memory = vs / 1000
				runningObject.Time = rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000
				runningObject.Status = RE
				err := runningObject.Proc.Kill()
				if err != nil {
					panic(err)
				}
				return &runningObject
			default:
			}
		}
		//0表示不发出信号
		tracer.Syscall(syscall.Signal(0))
	}
}
