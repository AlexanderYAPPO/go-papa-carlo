package target

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateTarget(t *testing.T) {
	structAPath := filepath.Join("..", "testdata", "reference", "success", "struct_a.go")
	structBPath := filepath.Join("..", "testdata", "reference", "success", "pkg", "struct_b.go")
	structCPath := filepath.Join("..", "testdata", "reference", "success", "pkg", "sub_pkg", "struct_c.go")
	structBadPath := filepath.Join("..", "testdata", "reference", "bad_mod_file", "struct_bad.go")

	outputAPath := filepath.Join("..", "testdata", "reference", "success", "a_builder_gen.go")
	outputABuildersPath := filepath.Join("..", "testdata", "reference", "success", "builders", "a_builder_gen.go")
	outputBPath := filepath.Join("..", "testdata", "reference", "success", "builders", "b_builder_gen.go")
	outputCPath := filepath.Join("..", "testdata", "reference", "success", "builders", "c_builder_gen.go")
	outputBadPath := filepath.Join("..", "testdata", "reference", "bad_mod_file", "builders", "bad_builder_gen.go")

	testCases := []struct {
		name          string
		targetName    string
		structPath    string
		outputPath    string
		parsingResult entity.ParsingResult
		want          entity.Target
		wantErr       error
	}{
		{
			name:       "struct A in module root with default output path",
			targetName: "A",
			structPath: structAPath,
			outputPath: outputAPath,
			parsingResult: entity.ParsingResult{
				Imports:     []entity.Import{},
				Fields:      []entity.Field{{Name: "Name", Type: "string"}},
				PackageName: "success",
			},
			want: entity.Target{
				Name:      "A",
				Reference: "A",
				ParsingResult: entity.ParsingResult{
					Imports:     []entity.Import{},
					Fields:      []entity.Field{{Name: "Name", Type: "string"}},
					PackageName: "success",
				},
			},
			wantErr: nil,
		},
		{
			name:       "struct B in nested pkg with custom output path",
			targetName: "B",
			structPath: structBPath,
			outputPath: outputBPath,
			parsingResult: entity.ParsingResult{
				Imports:     []entity.Import{},
				Fields:      []entity.Field{{Name: "Count", Type: "int"}},
				PackageName: "pkg",
			},
			want: entity.Target{
				Name:      "B",
				Reference: "targetpkg.B",
				Import: entity.Import{
					ReferenceName: "targetpkg",
					Path:          "example.com/reference/pkg",
					IsAlias:       true,
				},
				ParsingResult: entity.ParsingResult{
					Imports:     []entity.Import{},
					Fields:      []entity.Field{{Name: "Count", Type: "int"}},
					PackageName: "builders",
				},
			},
			wantErr: nil,
		},
		{
			name:       "struct A in module root with output in builders",
			targetName: "A",
			structPath: structAPath,
			outputPath: outputABuildersPath,
			parsingResult: entity.ParsingResult{
				Imports:     []entity.Import{},
				Fields:      []entity.Field{{Name: "Name", Type: "string"}},
				PackageName: "success",
			},
			want: entity.Target{
				Name:      "A",
				Reference: "targetpkg.A",
				Import: entity.Import{
					ReferenceName: "targetpkg",
					Path:          "example.com/reference",
					IsAlias:       true,
				},
				ParsingResult: entity.ParsingResult{
					Imports:     []entity.Import{},
					Fields:      []entity.Field{{Name: "Name", Type: "string"}},
					PackageName: "builders",
				},
			},
			wantErr: nil,
		},
		{
			name:       "struct C in sub package with custom output path",
			targetName: "C",
			structPath: structCPath,
			outputPath: outputCPath,
			parsingResult: entity.ParsingResult{
				Imports:     []entity.Import{},
				Fields:      []entity.Field{{Name: "Enabled", Type: "bool"}},
				PackageName: "sub_pkg",
			},
			want: entity.Target{
				Name:      "C",
				Reference: "targetpkg.C",
				Import: entity.Import{
					ReferenceName: "targetpkg",
					Path:          "example.com/reference/pkg/sub_pkg",
					IsAlias:       true,
				},
				ParsingResult: entity.ParsingResult{
					Imports:     []entity.Import{},
					Fields:      []entity.Field{{Name: "Enabled", Type: "bool"}},
					PackageName: "builders",
				},
			},
			wantErr: nil,
		},
		{
			name:       "malformed go.mod returns module directive not found error",
			targetName: "Bad",
			structPath: structBadPath,
			outputPath: outputBadPath,
			parsingResult: entity.ParsingResult{
				Imports:     []entity.Import{},
				Fields:      []entity.Field{{Name: "Value", Type: "int"}},
				PackageName: "bad_mod_file",
			},
			want:    entity.Target{},
			wantErr: errors.New("module directive not found in go.mod"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CreateTarget(tc.parsingResult, tc.targetName, tc.structPath, tc.outputPath)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
