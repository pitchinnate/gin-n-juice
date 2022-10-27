package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"gin-n-juice/config"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	flags           = flag.NewFlagSet("generator", flag.ExitOnError)
	generatorType   = flags.String("type", "", "Type of file(s) to generate: model|resource")
	generatorName   = flags.String("name", "", "Name of model or resource")
	generatorFields = flags.String("fields", "", "Fields for model: name:type,name:type,...")
)

type ModelInfo struct {
	ModelName      string
	SingleInstance string
	PluralInstance string
	PackageName    string
	Properties     []ModelProperty
}

type ModelProperty struct {
	Name string
	Type string
	Json string
}

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	if len(os.Args) == 1 {
		log.Fatalf("Type and name are required. Use --help for more info.")
	}

	loadEnv()
	config.SetupEnv()

	flags.Parse(os.Args[1:])

	pluralize := pluralize.NewClient()

	modelInfo := ModelInfo{
		ModelName:      strcase.ToCamel(*generatorName),
		SingleInstance: strcase.ToSnake(*generatorName),
		PluralInstance: pluralize.Plural(strcase.ToSnake(*generatorName)),
		PackageName:    config.PACKAGE_NAME,
	}

	if *generatorFields != "" {
		props := []ModelProperty{}
		pieces := strings.Split(*generatorFields, ",")
		for _, field := range pieces {
			pieces2 := strings.Split(field, ":")
			if len(pieces2) < 2 {
				log.Fatal("Each property needs a name and a type: ", field)
				continue
			}
			propType := pieces2[1]
			validTypes := []string{"string", "bool", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64",
				"int", "uint", "uintptr", "float32", "float64", "complex64", "complex128"}
			if !slices.Contains(validTypes, propType) {
				log.Fatal("Invalid property type: ", field)
				continue
			}

			props = append(props, ModelProperty{
				strcase.ToCamel(pieces2[0]),
				pieces2[1],
				strcase.ToSnake(pieces2[0]),
			})
		}
		modelInfo.Properties = props
	}

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if *generatorType == "resource" {
		templateDirectory := fmt.Sprintf("%s/templates/routes", directory)
		routeDirectory := fmt.Sprintf("%s/routes/%s", directory, modelInfo.PluralInstance)

		files, err := ioutil.ReadDir(templateDirectory)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(routeDirectory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(routeDirectory, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, file := range files {
			pieces := strings.Split(file.Name(), ".")
			templateFilePath := fmt.Sprintf("%s/%s", templateDirectory, file.Name())
			newFilePath := fmt.Sprintf("%s/%s.go", routeDirectory, pieces[0])

			if _, err := os.Stat(newFilePath); errors.Is(err, os.ErrNotExist) {
				tmpl, err := template.ParseFiles(templateFilePath)
				if err != nil {
					log.Printf("Error using template: %s", templateFilePath)
					log.Print(err)
					continue
				}
				writer := new(bytes.Buffer)
				err = tmpl.Execute(writer, modelInfo)
				if err != nil {
					log.Printf("Error executing template: %s", templateFilePath)
					log.Print(err)
					continue
				}
				err = os.WriteFile(newFilePath, writer.Bytes(), 0644)
				if err != nil {
					log.Printf("Error writing file: %s", newFilePath)
					log.Print(err)
					continue
				}
			} else {
				log.Printf("File already exists: %s", newFilePath)
			}
		}
	}

	if *generatorType == "model" {
		templateFilePath := fmt.Sprintf("%s/templates/model/model.tmpl", directory)
		modelDirectory := fmt.Sprintf("%s/models", directory)
		newFilePath := fmt.Sprintf("%s/%s.go", modelDirectory, modelInfo.SingleInstance)

		if _, err := os.Stat(newFilePath); errors.Is(err, os.ErrNotExist) {
			tmpl, err := template.ParseFiles(templateFilePath)
			if err != nil {
				log.Printf("Error using template: %s", templateFilePath)
				log.Fatal(err)
			}
			writer := new(bytes.Buffer)
			err = tmpl.Execute(writer, modelInfo)
			if err != nil {
				log.Printf("Error executing template: %s", templateFilePath)
				log.Fatal(err)
			}
			err = os.WriteFile(newFilePath, writer.Bytes(), 0644)
			if err != nil {
				log.Printf("Error writing file: %s", newFilePath)
				log.Fatal(err)
			}
		} else {
			log.Printf("File already exists: %s", newFilePath)
		}
	}
}

func loadEnv() {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(fmt.Sprintf("%s/.env", directory))
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}
