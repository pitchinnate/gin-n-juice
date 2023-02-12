package command

import (
	"fmt"
	"gin-n-juice/config"
	"github.com/go-cmd/cmd"
	"math/rand"
	"strings"
	"time"
)

func RunCommand(command string, appState chan string) {
	id := rand.Intn(100000)
	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] Running command: %s %d \n", command, id)
	}

	pieces := strings.Split(command, " ")
	c := cmd.NewCmd(pieces[0], pieces[1:]...)

	lastStdout := 0
	lastStderr := 0
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			status := c.Status()
			n := len(status.Stdout)
			if n > 0 {
				for i := lastStdout; i < n; i++ {
					fmt.Println(status.Stdout[i])
				}
				lastStdout = n
			}
			n = len(status.Stderr)
			if n > 0 {
				for i := lastStderr; i < n; i++ {
					fmt.Println(status.Stderr[i])
				}
				lastStderr = n
			}
		}
	}()

	statusChan := c.Start()
	commandChan := make(chan string)

	go func() {
		complete := false
		for {
			if complete {
				fmt.Printf("[GIN-N-JUICE] Command complete complete loop: %s %d \n", command, id)
				break
			}
			select {
			case state := <-appState:
				if state == "restart" || state == "stop" {
					//fmt.Printf("[GIN-N-JUICE] Killing command: %s %d \n", command, id)
					err := c.Stop()
					if err != nil {
						fmt.Printf("[GIN-N-JUICE] Error Killing command: %s \n", err.Error())
					}
				}
			case <-commandChan:
				fmt.Printf("[GIN-N-JUICE] Command complete kill loop: %s %d \n", command, id)
				complete = true
			}
		}
	}()

	// this waits until the command is done or killed
	<-statusChan

	commandChan <- "done"

	// code below will not fire till after statusChan is updated
	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] command done: %s %d \n", command, id)
	}

	// dump the remaining contents in stdout and stderr
	status := c.Status()
	n := len(status.Stdout)
	if n > 0 {
		for i := lastStdout; i < n; i++ {
			fmt.Println(status.Stdout[i])
		}
		lastStdout = n
	}
	n = len(status.Stderr)
	if n > 0 {
		for i := lastStderr; i < n; i++ {
			fmt.Println(status.Stderr[i])
		}
		lastStderr = n
	}
}
