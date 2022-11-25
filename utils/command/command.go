package command

import (
	"fmt"
	"gin-n-juice/config"
	"github.com/go-cmd/cmd"
	"math/rand"
	"time"
)

func RunCommand(command string, appState chan string) {
	id := rand.Intn(100000)
	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] Running command: %s %d \n", command, id)
	}

	c := cmd.NewCmd(command)
	statusChan := c.Start()

	lastStdout := 0
	lastStderr := 0
	ticker := time.NewTicker(time.Second)

	go func() {
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

	for {
		killLoop := false
		select {
		case status := <-statusChan:
			if status.StopTs > 0 {
				//fmt.Printf("[GIN-N-JUICE] command done: %s %d \n", command, id)
				killLoop = true
			}
		case state := <-appState:
			if state == "restart" || state == "stop" {
				//fmt.Printf("[GIN-N-JUICE] Killing command: %s %d \n", command, id)
				c.Stop()
			}
		}
		if killLoop {
			break
		}
	}

	if config.DEBUG {
		fmt.Printf("[GIN-N-JUICE] command closed: %s %d \n", command, id)
	}
}
