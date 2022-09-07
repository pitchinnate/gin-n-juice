package main

import (
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
	flags = flag.NewFlagSet("jackwebbit", flag.ExitOnError)
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	flags.Parse(os.Args[1:])
	args := flags.Args()

	log.Print("args:", args)

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Print("Current Directory: ", directory)

	if len(args) == 0 || args[0] == "serve" {
		var wg sync.WaitGroup
		runCommand(fmt.Sprintf("cd %s/cmd/server && go build -o ../../tmp/app.exe", directory), wg)
		runCommand(fmt.Sprintf("%s/tmp/app.exe", directory), wg)
		wg.Wait()
	} else if len(args) > 0 && args[0] == "migrate" {
		var wg sync.WaitGroup
		runCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate.exe", directory), wg)
		runCommand(fmt.Sprintf("%s/tmp/migrate.exe %s", directory, strings.Join(args[1:], " ")), wg)
		wg.Wait()
	} else if len(args) > 0 && args[0] == "test" {
		var wg sync.WaitGroup
		os.Remove(fmt.Sprintf("%s/tmp/test.db", directory))
		runCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate.exe", directory), wg)
		runCommand(fmt.Sprintf("%s/tmp/migrate.exe -testing up", directory), wg)
		runCommand(fmt.Sprintf("cd %s && go test ./routes/... %s", directory, strings.Join(args[1:], " ")), wg)
		wg.Wait()
	}
}

func runCommand(command string, wg sync.WaitGroup) {
	if runtime.GOOS == "windows" {
		wg.Add(1)
		cmd := exec.Command("cmd", "/C", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Print("Running command: ", command)
		err := cmd.Start()
		if err != nil {
			log.Print("command error: ", err, " command: ", command)
		}
		err = cmd.Wait()
		if err != nil {
			log.Print("command wait error: ", err, " command: ", command)
		}
		wg.Done()
	} else {
		panic("linux not supported yet")
		//command := strings.Join(arguements, " ")
		//fullCommand := fmt.Sprintf("%s %s", path, command)
		//cmd := exec.Command("bash", "-c", fullCommand)
		//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		//if debugMode {
		//	cmd.Stdout = os.Stdout
		//	cmd.Stderr = os.Stderr
		//}
		//return cmd
	}
}
