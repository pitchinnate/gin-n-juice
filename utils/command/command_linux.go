package command

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"tagdeploy/config"
)

func RunCommand(command string, appState chan string) {
	var cmd *exec.Cmd

	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] Running command: %s\n", command)
	}

	commandState := make(chan bool)
	killed := false
	go func() {
		for {
			done := false
			select {
			case state := <-appState:
				if state == "stop" {
					killed = true
					syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				}
			case <-commandState:
				done = true
			}
			if done {
				break
			}
		}
	}()

	cmd = exec.Command("bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil && config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] command error: %s command: %s\n", err.Error(), command)
	}
	if err == nil && config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] Running command with a PID of:  %d\n", cmd.Process.Pid)
	}

	err = cmd.Wait()
	if err != nil && config.DEBUG && !killed {
		fmt.Printf("[GIN-N-JUICE] command wait error: %s command: %s", err.Error(), command)
	}

	commandState <- false
}
