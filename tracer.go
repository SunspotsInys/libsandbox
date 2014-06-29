package sandbox

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/hjr265/ptrace.go/ptrace"
)

type RunningObject struct {
	Time        syscall.Timeval
	TimeLimit   int64
	MemoryLimit int64
	Memory      int64
}

func (r *RunningObject) Millisecond() int64 {
	return r.Time.Sec*1000 + r.Time.Usec/1000
}

func Run(src string, args []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var rusage syscall.Rusage
	var incall = true
	var runningObject RunningObject
	proc, err := os.StartProcess(src, args, &os.ProcAttr{Sys: &syscall.SysProcAttr{
		Ptrace: true},
	})
	if err != nil {
		panic(err)
	}
	//set CPU time limit
	var rlimit syscall.Rlimit
	rlimit.Cur = 1
	rlimit.Max = 1 + 1
	err = prLimit(proc.Pid, syscall.RLIMIT_CPU, &rlimit)
	if err != nil {
		fmt.Println(err)
		return
	}
	rlimit.Cur = 100
	rlimit.Max = 1024 + 1024
	err = prLimit(proc.Pid, syscall.RLIMIT_STACK, &rlimit)
	if err != nil {
		fmt.Println(err)
		return
	}
	tracer, err := ptrace.Attach(proc)
	if err != nil {
		panic(err)
	}
	for {
		status := syscall.WaitStatus(0)
		_, err := syscall.Wait4(proc.Pid, &status, syscall.WSTOPPED, &rusage)
		if err != nil {
			panic(err)
		}
		if status.Exited() {
			fmt.Println("exit")
			fmt.Println(rusage.Stime)
			break
		}
		if status.CoreDump() {
			fmt.Println("CoreDump")
			return
		}
		if status.Continued() {
			fmt.Println("Continued")
			return
		}
		if status.Signaled() {
			return
		}
		if status.Stopped() && status.StopSignal() != syscall.SIGTRAP {
			switch status.StopSignal() {
			case syscall.SIGXCPU:
				fmt.Println("SIGXCPU")
				runningObject.Time = rusage.Utime
				fmt.Println(runningObject.Millisecond())
			case syscall.SIGSEGV:
				fmt.Println("SIGSEGV")
				runningObject.Memory = rusage.Minflt
				fmt.Println(runningObject.Memory)
			default:
				fmt.Println("default")
			}
			return
		} else {
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
		//0表示不发出信号
		err = tracer.Syscall(syscall.Signal(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}
