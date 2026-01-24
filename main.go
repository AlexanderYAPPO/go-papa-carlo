package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlexanderYAPPO/go-papa-carlo/pipeline"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <struct_name> <path_to_struct>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	structName := os.Args[1]
	pathToStruct := os.Args[2]

	if err := pipeline.GenerateToFile(structName, pathToStruct); err != nil {
		log.Fatal(err)
	}
}
