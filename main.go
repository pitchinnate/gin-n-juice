package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type RunCommand struct {
	Command string
	Long    bool
}

var (
	flags = flag.NewFlagSet("gin-n-juice", flag.ExitOnError)
	debug = flags.Bool("debug", false, "Enable debug mode")
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	flags.Parse(os.Args[1:])
	args := flags.Args()

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if *debug {
		log.Printf("Working Directory: %s", directory)
	}

	commandExtension := ""
	if runtime.GOOS == "windows" {
		commandExtension = ".exe"
	}

	path := fmt.Sprintf("%s/tmp", directory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(args) == 0 || args[0] == "serve" {
		var wg sync.WaitGroup
		runCommand(fmt.Sprintf("cd %s/cmd/server && go build -o ../../tmp/app%s", directory, commandExtension), wg)
		runCommand(fmt.Sprintf("%s/tmp/app%s", directory, commandExtension), wg)
		wg.Wait()
	} else if len(args) > 0 && args[0] == "migrate" {
		var wg sync.WaitGroup
		runCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate%s", directory, commandExtension), wg)
		runCommand(fmt.Sprintf("%s/tmp/migrate%s %s", directory, commandExtension, strings.Join(args[1:], " ")), wg)
		wg.Wait()
	} else if len(args) > 0 && args[0] == "generator" {
		var wg sync.WaitGroup
		runCommand(fmt.Sprintf("cd %s/cmd/generator && go build -o ../../tmp/generator%s", directory, commandExtension), wg)
		runCommand(fmt.Sprintf("%s/tmp/generator%s %s", directory, commandExtension, strings.Join(args[1:], " ")), wg)
		wg.Wait()
	} else if len(args) > 0 && args[0] == "test" {
		var wg sync.WaitGroup
		os.Remove(fmt.Sprintf("%s/tmp/test.db", directory))
		runCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate%s", directory, commandExtension), wg)
		runCommand(fmt.Sprintf("%s/tmp/migrate%s -testing up", directory, commandExtension), wg)
		runCommand(fmt.Sprintf("cd %s && go test ./routes/... %s", directory, strings.Join(args[1:], " ")), wg)
		wg.Wait()
	}
}

func runCommand(command string, wg sync.WaitGroup) {
	var cmd *exec.Cmd
	wg.Add(1)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if *debug {
		log.Print("Running command: ", command)
	}
	err := cmd.Start()
	if err != nil && *debug {
		log.Print("command error: ", err, " command: ", command)
	}
	err = cmd.Wait()
	if err != nil && *debug {
		log.Print("command wait error: ", err, " command: ", command)
	}
	wg.Done()
}
