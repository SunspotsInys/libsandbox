//Sandbox package is used for sandbox command.
//
//Sandbox uses /proc/{id}/stats to check virtual memory usage and,/proc/uptime for time run.
//Every timer tick,send a signal to check the running status,if any erros happend kill the
//process and report error,or check the standard input and output to report wrong answer or accept.
package sandbox

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	AC  uint64 = iota //Accept
	PE                //Present Erro
	TLE               //Time Limit Out Error
	MLE               //Memory Limit Out Error
	WA                //Wrong Answer
	RE                //Runtime Error
	OLE               //Output Limit Error
	CE                //Complie Error
	SE                //Segmenfault Error
)

const (
	KB = 1024 //KB==1024 bytes
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

// RunningObject is a process running information container.
type RunningObject struct {
	Proc        *os.Process
	TimeLimit   int64
	MemoryLimit int64
	Memory      int64  //KB
	Time        int64  //MS
	Status      uint64 //result status
}

// send signal SIGALRM to the process every tick.
func (r *RunningObject) runTick() {
	ticker := time.NewTicker(frequency)
	//send alarm signal with time tick frequency
	for _ = range ticker.C {
		r.Proc.Signal(os.Signal(unix.SIGALRM))
	}
}

// update virtual memory.
func (r *RunningObject) updateVirtualMemory() {
	r.Memory = virtualMemory(r.Proc.Pid) / KB
}

// update time used from the max of rusage and the /porc/{pid}/stat
func (r *RunningObject) updateTime(rusg *unix.Rusage) {
	r.Time = rusg.Utime.Sec*1000 +
		rusg.Utime.Usec/1000 // MS
	rt := realTime(r.Proc.Pid)
	if r.Time < rt {
		r.Time = rt
	}
}

// wether exceed resource limit
func (r *RunningObject) exceedLimit() uint64 {
	if r.Memory > r.MemoryLimit {
		return MLE
	}
	if r.Time > r.TimeLimit {
		return TLE
	}
	return 0
}

func wait(pid, options int, rusage *unix.Rusage) (int, *unix.WaitStatus, error) {
	var status unix.WaitStatus
	wpid, err := unix.Wait4(pid, &status, unix.WALL, rusage)
	return wpid, &status, err
}

// Compile compiles specific language source file
// and build into destination file.
func Complie(src string, des string, lan uint64) error {
	return compile(src, des, lan)
}

// Run runs the binary,and receive reader and writer for standard input and output,
// args are the binary arguments,timeLimit and memoryLimit are in MS and KB.
func Run(bin string, reader io.Reader, writer io.Writer,
	args []string, timeLimit int64, memoryLimit int64) *RunningObject {

	// We must ensure here that we are running on the same thread during
	// the execution of dbg. This is due to the fact that ptrace(2) expects
	// all commands after PTRACE_ATTACH to come from the same thread.

	runtime.LockOSThread()

	defer runtime.UnlockOSThread()

	var (
		rusage        unix.Rusage
		runningObject RunningObject
	)

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
	tracer, err := Attach(proc)
	if err != nil {
		panic(err)
	}
	runningObject.Proc = proc

	defer runningObject.Proc.Kill()

	go runningObject.runTick()
	setTimelimit(runningObject.Proc.Pid, timeLimit)
	if err != nil {
		fmt.Println(err)
		return &runningObject
	}

	for {
		_, status, err := wait(proc.Pid, unix.WSTOPPED, &rusage)
		if err != nil {
			panic(err)
		}
		// status exited
		if status.Exited() {
			return &runningObject
		}

		if status.CoreDump() {
			fmt.Println("CoreDump")
			return &runningObject
		}

		if status.Stopped() &&
			status.StopSignal() != unix.SIGTRAP {
			runningObject.updateVirtualMemory()
			runningObject.updateTime(&rusage)
			switch status.StopSignal() {
			case unix.SIGALRM:
				if typ := runningObject.exceedLimit(); typ != 0 {
					runningObject.Status = typ
					return &runningObject
				}
			case unix.SIGXCPU:
				runningObject.Status = TLE
				return &runningObject
			case unix.SIGSEGV:
				runningObject.Status = RE
				return &runningObject
			default:
			}
		}
		//0表示不发出信号
		tracer.Syscall(syscall.Signal(0))
	}
}
