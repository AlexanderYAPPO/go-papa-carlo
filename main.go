package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlexanderYAPPO/go-papa-carlo/pipeline"
)

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: %s <struct_name> <path_to_struct> [output_path]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	structName := os.Args[1]
	pathToStruct := os.Args[2]
	outputPathArg := ""
	if len(os.Args) == 4 {
		outputPathArg = os.Args[3]
	}
	outputPath := resolveOutputPath(structName, pathToStruct, outputPathArg)

	if err := pipeline.GenerateToFile(structName, pathToStruct, outputPath); err != nil {
		log.Fatal(err)
	}
}

func resolveOutputPath(structName, pathToStruct, cliOutputPath string) string {
	if cliOutputPath != "" {
		return cliOutputPath
	}

	outputName := structName + "_builder_gen.go"
	return filepath.Join(filepath.Dir(pathToStruct), outputName)
}
