package pipeline

import (
	"os"
	"path/filepath"

	"github.com/AlexanderYAPPO/go-papa-carlo/generate"
	"github.com/AlexanderYAPPO/go-papa-carlo/parse"
)

func GenerateToFile(structName, pathToStruct string) error {
	res := parse.Parse(structName, pathToStruct)
	generatedCode := generate.Generate(res)

	outputName := structName + "_builder_gen.go"
	outputPath := filepath.Join(filepath.Dir(pathToStruct), outputName)
	return os.WriteFile(outputPath, []byte(generatedCode), 0644)
}
