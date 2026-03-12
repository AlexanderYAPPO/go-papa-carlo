package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AlexanderYAPPO/go-papa-carlo/pipeline"
)

const goModContent = `module example.com/fixture

go 1.20

require github.com/stretchr/testify v1.9.0
`

func TestBuilderScenarios(t *testing.T) {
	t.Run("simple struct with few fields", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "simple_struct", "fixture", "pkg1", "struct_with_few_fields.go"),
			"StructWithFewFields",
			filepath.Join("testdata", "simple_struct", "consumer", "build_struct_with_few_fields_test.go"),
			"",
		)
	})
	t.Run("omit fields by omit tag", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "omit_field", "fixture", "pkg1", "struct_with_omit.go"),
			"StructWithOmittedField",
			filepath.Join("testdata", "omit_field", "consumer", "build_struct_with_omit_test.go"),
			"",
		)
	})
	t.Run("struct where some of the fields are optional", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "optional_field", "fixture", "pkg1", "struct_with_optional.go"),
			"StructWithOptionalField",
			filepath.Join("testdata", "optional_field", "consumer", "build_struct_with_optional_test.go"),
			"",
		)
	})
	t.Run("struct where all fields are optional", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "optional_field", "fixture", "pkg1", "struct_with_only_optional.go"),
			"StructWithOnlyOptional",
			filepath.Join("testdata", "optional_field", "consumer", "build_struct_with_only_optional_test.go"),
			"",
		)
	})
	t.Run("generate a builder in a different from source package", func(t *testing.T) {
		runScenarioWithCustomOutputPackage(t,
			filepath.Join("testdata", "custom_output_path", "fixture", "pkg1", "struct_with_few_fields.go"),
			"StructWithFewFields",
			filepath.Join("testdata", "custom_output_path", "consumer", "build_struct_with_custom_output_path_test.go"),
		)
	})
	t.Run("struct with omitted private fields", func(t *testing.T) {
		runScenario(t,
			filepath.Join("testdata", "struct_with_private_fields", "fixture", "pkg", "struct_with_private_fields.go"),
			"StructWithPrivateFieldsThatWorks",
			filepath.Join("testdata", "struct_with_private_fields", "consumer", "build_struct_private_fields.go"),
			"",
		)
	})
	t.Run("throw an error when a struct has not omitted private fields", func(t *testing.T) {
		runScenarioWithExpectedError(t,
			filepath.Join("testdata", "struct_with_private_fields", "fixture", "pkg", "struct_with_private_fields.go"),
			"StructWithPrivateFieldsThatErrors",
			"",
			"field privateField is a private field and cannot be used in builder generation unless omitted with tag papa-carlo:\"omit\"",
		)
	})
	t.Run("throw an error when a struct uses an unexported type", func(t *testing.T) {
		runScenarioWithExpectedError(t,
			filepath.Join("testdata", "struct_with_unexported_types", "fixture", "pkg", "struct_with_unexported_types.go"),
			"StructWithUnexportedType",
			"",
			"field PublicField uses a type containing unexported identifiers: privateType",
		)
	})

	t.Run("throw an error when a file is found but not the struct", func(t *testing.T) {
		runScenarioWithExpectedError(t,
			filepath.Join("testdata", "simple_struct", "fixture", "pkg1", "struct_with_few_fields.go"),
			"UnknownStruct",
			"",
			"requested struct is not found in the file",
		)
	})

	t.Run("throw an error if the source file doesn't exist", func(t *testing.T) {
		expectedErrMsg := "open /file/that/doesnt/exist.go: no such file or directory"
		err := pipeline.GenerateToFile("FakeStructName", "/file/that/doesnt/exist.go", "")
		if err == nil {
			t.Fatalf("expected error containing %q, got nil", expectedErrMsg)
		}
		if !strings.Contains(err.Error(), expectedErrMsg) {
			t.Fatalf("expected error containing %q, got %q", expectedErrMsg, err.Error())
		}
	})
}

func runScenario(t *testing.T, structPath, structName, usagePath, outputFileName string) {
	t.Helper()
	tempDir := createPackage(t)

	structDst := copyGoFile(t, tempDir, structPath)
	outputPath := filepath.Join(filepath.Dir(structDst), structName+"_builder_gen.go")
	if outputFileName != "" {
		outputPath = filepath.Join(filepath.Dir(structDst), outputFileName)
	}
	if err := pipeline.GenerateToFile(structName, structDst, outputPath); err != nil {
		t.Fatalf("generate builder: %v", err)
	}

	copyGoFile(t, tempDir, usagePath)
	runGoTest(t, tempDir)
}

func runScenarioWithCustomOutputPackage(t *testing.T, sourceStructPath, structName, usagePath string) {
	t.Helper()
	tempDir := createPackage(t)

	sourceStructDst := copyGoFile(t, tempDir, sourceStructPath)
	outputDir := filepath.Join(tempDir, "builders")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("create custom output dir: %v", err)
	}
	outputPath := filepath.Join(outputDir, "custom_builder_gen.go")
	if err := pipeline.GenerateToFile(structName, sourceStructDst, outputPath); err != nil {
		t.Fatalf("generate builder: %v", err)
	}

	copyGoFile(t, tempDir, usagePath)
	runGoTest(t, tempDir)
}

func runScenarioWithExpectedError(t *testing.T, structPath, structName, outputFileName, expectedErrMsg string) {
	t.Helper()
	tempDir := createPackage(t)

	structDst := copyGoFile(t, tempDir, structPath)
	outputPath := filepath.Join(filepath.Dir(structDst), structName+"_builder_gen.go")
	if outputFileName != "" {
		outputPath = filepath.Join(filepath.Dir(structDst), outputFileName)
	}

	err := pipeline.GenerateToFile(structName, structDst, outputPath)
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", expectedErrMsg)
	}
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Fatalf("expected error containing %q, got %q", expectedErrMsg, err.Error())
	}
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
