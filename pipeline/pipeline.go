package pipeline

import (
	"os"
	"path/filepath"

	"github.com/AlexanderYAPPO/go-papa-carlo/generate"
	"github.com/AlexanderYAPPO/go-papa-carlo/target"
)

func GenerateToFile(structName, pathToStruct, outputPath string) error {
	res := target.Parse(structName, pathToStruct)
	adaptedTarget, err := target.CreateTarget(res, structName, pathToStruct, outputPath)
	if err != nil {
		return err
	}
	generatedCode := generate.Generate(adaptedTarget)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(generatedCode), 0644)
}
