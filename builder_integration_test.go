package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/AlexanderYAPPO/go-papa-carlo/pipeline"
)

const goModContent = `module example.com/fixture

go 1.20

require github.com/stretchr/testify v1.9.0
`

func TestBuilderScenarios(t *testing.T) {
	t.Run("struct_with_few_fields", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "simple_struct", "fixture", "pkg1", "struct_with_few_fields.go"),
			"StructWithFewFields",
			filepath.Join("testdata", "simple_struct", "consumer", "build_struct_with_few_fields_test.go"),
		)
	})
	t.Run("struct_with_omitted_field", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "omit_field", "fixture", "pkg1", "struct_with_omit.go"),
			"StructWithOmittedField",
			filepath.Join("testdata", "omit_field", "consumer", "build_struct_with_omit_test.go"),
		)
	})
	t.Run("struct_with_optional_field", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "optional_field", "fixture", "pkg1", "struct_with_optional.go"),
			"StructWithOptionalField",
			filepath.Join("testdata", "optional_field", "consumer", "build_struct_with_optional_test.go"),
		)
	})
	t.Run("struct_with_only_optional_fields", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "optional_only_field", "fixture", "pkg1", "struct_with_only_optional.go"),
			"StructWithOnlyOptional",
			filepath.Join("testdata", "optional_only_field", "consumer", "build_struct_with_only_optional_test.go"),
		)
	})
}

func runScenario(t *testing.T, structPath, structName, usagePath string) {
	t.Helper()
	tempDir := createPackage(t)

	structDst := copyGoFile(t, tempDir, structPath)
	if err := pipeline.GenerateToFile(structName, structDst); err != nil {
		t.Fatalf("generate builder: %v", err)
	}

	copyGoFile(t, tempDir, usagePath)
	runGoTest(t, tempDir)
}

func createPackage(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	return tempDir
}

func copyGoFile(t *testing.T, tempDir, srcPath string) string {
	t.Helper()
	dstDir := filepath.Join(tempDir, filepath.Base(filepath.Dir(srcPath)))
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("create destination dir: %v", err)
	}
	dstPath := filepath.Join(dstDir, filepath.Base(srcPath))
	copyFile(t, srcPath, dstPath)
	return dstPath
}

func runGoTest(t *testing.T, tempDir string) {
	t.Helper()
	cmd := exec.Command("go", "test", "-mod=mod", "./...")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed: %v\n%s", err, output)
	}
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}
