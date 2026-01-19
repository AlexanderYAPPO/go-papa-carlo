package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlexanderYAPPO/go-papa-carlo/generate"
	"github.com/AlexanderYAPPO/go-papa-carlo/parse"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <struct_name> <path_to_struct>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	structName := os.Args[1]
	pathToStruct := os.Args[2]

	res := parse.Parse(structName, pathToStruct)
	generatedCode := generate.Generate(res)

	outputName := structName + "_builder_gen.go"
	outputPath := filepath.Join(filepath.Dir(pathToStruct), outputName)
	if err := os.WriteFile(outputPath, []byte(generatedCode), 0644); err != nil {
		log.Fatal(err)
	}
}
