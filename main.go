package main

import (
	"errors"
	"flag"
	"fmt"
	"gin-n-juice/config"
	"gin-n-juice/utils/command"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	flags = flag.NewFlagSet("gin-n-juice", flag.ExitOnError)
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	loadEnv()
	config.SetupEnv()

	flags.Parse(os.Args[1:])
	args := flags.Args()

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if config.DEBUG {
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

	c := make(chan os.Signal, 1)
	appState := make(chan string)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			select {
			case <-c:
				appState <- "stop"
			}
		}
	}()

	if len(args) == 0 || args[0] == "serve" {
		command.RunCommand(fmt.Sprintf("cd %s/cmd/server && go build -o ../../tmp/app%s", directory, commandExtension), appState)
		command.RunCommand(fmt.Sprintf("%s/tmp/app%s", directory, commandExtension), appState)
	} else if len(args) > 0 && args[0] == "migrate" {
		command.RunCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate%s", directory, commandExtension), appState)
		command.RunCommand(fmt.Sprintf("%s/tmp/migrate%s %s", directory, commandExtension, strings.Join(args[1:], " ")), appState)
	} else if len(args) > 0 && args[0] == "generator" {
		command.RunCommand(fmt.Sprintf("cd %s/cmd/generator && go build -o ../../tmp/generator%s", directory, commandExtension), appState)
		command.RunCommand(fmt.Sprintf("%s/tmp/generator%s %s", directory, commandExtension, strings.Join(args[1:], " ")), appState)
	} else if len(args) > 0 && args[0] == "test" {
		os.Remove(fmt.Sprintf("%s/tmp/test.db", directory))
		command.RunCommand(fmt.Sprintf("cd %s/cmd/migrations && go build -o ../../tmp/migrate%s", directory, commandExtension), appState)
		command.RunCommand(fmt.Sprintf("%s/tmp/migrate%s -testing up", directory, commandExtension), appState)
		command.RunCommand(fmt.Sprintf("cd %s && go test ./routes/... ./models/... %s", directory, strings.Join(args[1:], " ")), appState)
	} else if len(args) > 0 && args[0] == "rename" {
		if len(args) == 1 {
			log.Fatal("Must pass a new name")
		}
		RenamePackage(directory, args[1])
	}
}

func RenamePackage(directory string, newName string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		//log.Printf("Error reading dir %s", err)
	} else {
		for _, file := range files {
			path := filepath.Join(directory, file.Name())
			if file.IsDir() {
				RenamePackage(path, newName)
			} else {
				lastThree := file.Name()[len(file.Name())-3:]
				if lastThree == ".go" || file.Name() == "go.mod" {
					currentFile, err := os.ReadFile(path)
					if err == nil {
						fileString := string(currentFile)
						fileString = strings.ReplaceAll(fileString, "gin-n-juice", newName)
						if err := os.WriteFile(path, []byte(fileString), 0644); err != nil {
							log.Printf("ERROR - updating file: %s", path)
						} else {
							log.Printf("Updated file: %s", path)
						}
					}
				}
			}
		}
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}
