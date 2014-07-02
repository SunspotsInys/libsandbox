package sandbox

import (
	"fmt"
	"os"
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
	Memory      int64
	Time        int64
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

func Run(src string, args []string, timeLimit int64, memoryLimit int64) *RunningObject {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var rusage syscall.Rusage
	//var incall = true
	var runningObject RunningObject
	runningObject.TimeLimit = timeLimit
	runningObject.MemoryLimit = memoryLimit
	proc, err := os.StartProcess(src, args, &os.ProcAttr{Sys: &syscall.SysProcAttr{
		Ptrace: true},
	})
	if err != nil {
		panic(err)
	}
	tracer, err := ptrace.Attach(proc)
	if err != nil {
		panic(err)
	}
	runningObject.Proc = proc
	go runningObject.RunTick(time.Millisecond)
	//set CPU time limit
	var rlimit syscall.Rlimit
	rlimit.Cur = uint64(timeLimit)
	rlimit.Max = uint64(timeLimit)
	err = prLimit(proc.Pid, syscall.RLIMIT_CPU, &rlimit)
	if err != nil {
		fmt.Println(err)
		return &runningObject
	}
	/*
		rlimit.Cur = 1024
		rlimit.Max = rlimit.Cur + 1024
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
			fmt.Println("exit")
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
			return &runningObject
		}
		if status.Stopped() && status.StopSignal() != syscall.SIGTRAP {
			switch status.StopSignal() {
			case syscall.SIGALRM:
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
					err := runningObject.Proc.Kill()
					if err != nil {
						panic(err)
					}
					return &runningObject
				}
				vs := virtualMemory(runningObject.Proc.Pid)
				if vs/1000 > runningObject.MemoryLimit {
					runningObject.Memory = vs / 1000
					runningObject.Status = MLE
					err := runningObject.Proc.Kill()
					if err != nil {
						panic(err)
					}
					return &runningObject
				}
			case syscall.SIGXCPU:
				runningObject.Time = rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000
				runningObject.Status = TLE
				err := runningObject.Proc.Kill()
				if err != nil {
					panic(err)
				}
				return &runningObject
			case syscall.SIGSEGV:
				runningObject.Status = MLE
				err := runningObject.Proc.Kill()
				if err != nil {
					panic(err)
				}
				return &runningObject
			default:
			}
		} /* else {
			regs, err := tracer.GetRegs()
			if err != nil {
				panic(err)
			}
			if regs.Orig_rax == syscall.SYS_WRITE {
				if incall {
					incall = false

					_, err = tracer.GetRegs()
					if err != nil {
						panic(err)
					}
					fmt.Printf("The child made a system call with, %d,%d,%d \n", regs.Rdi, regs.Rsi, regs.Rdx)
				} else {
					incall = true
					regs, err := tracer.GetRegs()
					if err != nil {
						panic(err)
					}
					fmt.Printf("write returned %v\n", regs.Rax)
					fmt.Printf("call %d\n", regs.Rdi)
				}
			}
		}
		*/
		//0表示不发出信号
		err = tracer.Syscall(syscall.Signal(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}
