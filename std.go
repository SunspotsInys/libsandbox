package sandbox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var (
	OutOfTimeError   = errors.New("out of time")
	OutOfMemoryError = errors.New("out of memory")
)

type StdSandbox struct {
	Bin         string    // binary path
	Args        []string  // arguments
	Input       io.Reader // standard input
	TimeLimit   int64     // time limit in ms
	MemoryLimit int64     // memory limit in kb
}

func (s StdSandbox) Run() ([]byte, error) {

	cmd := exec.Command(s.Bin, s.Args...)
	if cmd.Stdin != nil {
		return nil, errors.New("stdin is not nil")
	}
	if cmd.Stderr != nil {
		return nil, errors.New("stdout is not nil")
	}
	if cmd.Stdout != nil {
		return nil, errors.New("stdout is not nil")
	}
	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	cmd.Stdout = buf
	cmd.Stdin = s.Input

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	setTimelimit(cmd.Process.Pid, s.TimeLimit/1000)

	// Send signal SIGALRM to the process every tick.
	go func() {
		ticker := time.NewTicker(TICK)
		for range ticker.C {
			err := cmd.Process.Signal(os.Signal(unix.SIGSTOP))
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	var rusage unix.Rusage
	for {
		_, status, err := wait(cmd.Process.Pid, unix.WSTOPPED, &rusage)
		if err != nil {
			fmt.Println("wait", err)
			return nil, err
		}
		if status.Exited() {
			return buf.Bytes(), nil
		}

		if status.Stopped() {
			switch status.StopSignal() {
			case unix.SIGSTOP:
				runningTime := RunningTime(cmd.Process.Pid)
				cpuTime := CpuTime(cmd.Process.Pid)
				if cpuTime > s.TimeLimit ||
					// Like sleep, some process consumes no cpu usage, but does
					// consume runnig time, so here limit real runnig time to
					// 150% cpu usage time.
					runningTime > 3*s.TimeLimit/2 {
					return nil, OutOfTimeError

				}

				vm := VirtualMemory(cmd.Process.Pid)
				rss := RssSize(cmd.Process.Pid)
				// RSS size dosen't include swap out memory,
				// virtual memory dosen't include memory demand-loaded int.
				// So set limit: memory < 150% * rss and vm > memory*150%
				if rss*3 > s.MemoryLimit*2 ||
					vm*2 > s.MemoryLimit*3 {
					fmt.Printf("rss %d, vm %d, limit %d", rss, vm, s.MemoryLimit)
					return nil, OutOfMemoryError

				}
			case unix.SIGXCPU:
				return nil, OutOfTimeError
			default:
				fmt.Println("default signal", status.StopSignal())

			}
		}
		syscall.Kill(cmd.Process.Pid, syscall.SIGCONT)
	}

	return buf.Bytes(), nil
}

type Sandbox interface {
	Run() (output []byte, err error)
}
