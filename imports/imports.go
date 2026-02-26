package imports

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
)

func PickUniqueAlias(imports []entity.Import, desiredAlias string) string {
	used := map[string]bool{}
	for _, imp := range imports {
		if imp.ReferenceName != "" {
			used[imp.ReferenceName] = true
		}
	}
	if !used[desiredAlias] {
		return desiredAlias
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s%d", desiredAlias, i)
		if !used[candidate] {
			return candidate
		}
	}
}

func BuildImportPath(moduleName, moduleFilePath, targetDirPath string) (string, error) {
	rel, err := filepath.Rel(moduleFilePath, targetDirPath)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return moduleName, nil
	}
	if strings.HasPrefix(rel, "..") {
		return "", errors.New("target path is outside of module root")
	}
	return moduleName + "/" + filepath.ToSlash(rel), nil
}
