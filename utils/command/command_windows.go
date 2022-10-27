package command

import (
	"fmt"
	"gin-n-juice/config"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

type process struct {
	Pid    int
	Handle uintptr
}

type ProcessExitGroup windows.Handle

func NewProcessExitGroup() (ProcessExitGroup, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return ProcessExitGroup(handle), nil
}

func (g ProcessExitGroup) Dispose() error {
	return windows.CloseHandle(windows.Handle(g))
}

func (g ProcessExitGroup) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(
		windows.Handle(g),
		windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}

func RunCommand(command string, appState chan string) {
	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] Running command: %s\n", command)
	}

	commandState := make(chan bool)
	g, err := NewProcessExitGroup()
	if err != nil {
		panic(err)
	}
	defer g.Dispose()

	var cmd *exec.Cmd
	go func() {
		for {
			done := false
			select {
			case state := <-appState:
				if state == "stop" {
					g.Dispose()
				}
			case <-commandState:
				done = true
			}
			if done {
				break
			}
		}
	}()

	cmd = exec.Command("cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil && config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] command error: %s command: %s\n", err.Error(), command)
	}

	if err := g.AddProcess(cmd.Process); err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil && config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] command wait error: %s command: %s\n", err.Error(), command)
	}

	commandState <- false
}
