package target

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
	"github.com/AlexanderYAPPO/go-papa-carlo/imports"
)

const (
	_targetPkgAliasBase = "targetpkg"
	_goModFileName      = "go.mod"
)

func CreateTarget(parsingResult entity.ParsingResult, targetName, pathToTargetStruct, outputPath string) (entity.Target, error) {
	target := entity.Target{
		Name:          targetName,
		Reference:     targetName,
		ParsingResult: parsingResult,
	}

	sourceDirPath := filepath.Clean(filepath.Dir(pathToTargetStruct))
	outputDirPath := filepath.Clean(filepath.Dir(outputPath))
	if sourceDirPath == outputDirPath {
		return target, nil
	}

	moduleName, modFilePath, err := findModuleInfo(sourceDirPath)
	if err != nil {
		return entity.Target{}, err
	}

	sourceImportPath, err := imports.BuildImportPath(moduleName, modFilePath, sourceDirPath)
	if err != nil {
		return entity.Target{}, err
	}

	alias := imports.PickUniqueAlias(parsingResult.Imports, _targetPkgAliasBase)
	target.Import = entity.Import{
		ReferenceName: alias,
		Path:          sourceImportPath,
		IsAlias:       true,
	}
	target.Reference = alias + "." + target.Name
	target.ParsingResult.PackageName = filepath.Base(outputDirPath)

	return target, nil
}

func findModuleInfo(startDir string) (modulePath string, moduleRoot string, err error) {
	dir := filepath.Clean(startDir)
	for {
		goModPath := filepath.Join(dir, _goModFileName)
		data, readErr := os.ReadFile(goModPath)
		if readErr == nil {
			modulePath, parseErr := parseModulePath(string(data))
			if parseErr != nil {
				return "", "", parseErr
			}
			return modulePath, dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", errors.New("unable to resolve module path from go.mod")
		}
		dir = parent
	}
}

func parseModulePath(goModContent string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(goModContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.HasPrefix(line, "module ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return "", errors.New("invalid module directive in go.mod")
		}
		return fields[1], nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("module directive not found in go.mod")
}
