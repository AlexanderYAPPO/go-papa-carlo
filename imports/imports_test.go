package imports

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
	"github.com/stretchr/testify/assert"
)

func TestPickUniqueAlias(t *testing.T) {
	testCases := []struct {
		name         string
		allImports   []entity.Import
		desiredAlias string
		expected     string
	}{
		{
			name: "returns desired alias when unused",
			allImports: []entity.Import{
				{ReferenceName: "fmt", Path: "fmt"},
				{ReferenceName: "time", Path: "time"},
			},
			desiredAlias: "targetpkg",
			expected:     "targetpkg",
		},
		{
			name: "returns next suffixed alias when desired alias is used",
			allImports: []entity.Import{
				{ReferenceName: "targetpkg", Path: "example.com/a"},
				{ReferenceName: "targetpkg1", Path: "example.com/b"},
				{ReferenceName: "targetpkg2", Path: "example.com/c"},
			},
			desiredAlias: "targetpkg",
			expected:     "targetpkg3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			alias := PickUniqueAlias(tc.allImports, tc.desiredAlias)
			assert.Equal(t, tc.expected, alias)
		})
	}
}

func TestBuildImportPath(t *testing.T) {
	_, filepathRelError := filepath.Rel("repo/mod", "/repo/mod/pkg")

	testCases := []struct {
		name          string
		moduleName    string
		moduleRoot    string
		targetDir     string
		expectedPath  string
		expectedError error
	}{
		{
			name:         "returns module name for module root",
			moduleName:   "example.com/mod",
			moduleRoot:   "/repo",
			targetDir:    "/repo",
			expectedPath: "example.com/mod",
		},
		{
			name:         "returns module plus relative path for subdirectory",
			moduleName:   "example.com/mod",
			moduleRoot:   filepath.Join("repo"),
			targetDir:    filepath.Join("repo", "internal", "builders"),
			expectedPath: "example.com/mod/internal/builders",
		},
		{
			name:          "returns error when target directory is outside module root",
			moduleName:    "example.com/mod",
			moduleRoot:    filepath.Join("repo", "mod"),
			targetDir:     filepath.Join("repo", "other"),
			expectedError: errors.New("target path is outside of module root"),
		},
		{
			name:          "returns library error from filepath.Rel as is",
			moduleName:    "example.com/mod",
			moduleRoot:    "repo/mod",
			targetDir:     "/repo/mod/pkg",
			expectedError: filepathRelError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modulePath, err := BuildImportPath(tc.moduleName, tc.moduleRoot, tc.targetDir)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPath, modulePath)
		})
	}
}
