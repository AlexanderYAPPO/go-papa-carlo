package entity

type ParsingResult struct {
	Imports     []Import
	Fields      []Field
	PackageName string
}

type Target struct {
	Name          string
	Reference     string
	Import        Import
	ParsingResult ParsingResult
}

type Import struct {
	ReferenceName string
	Path          string
	IsAlias       bool
}

type Field struct {
	Name           string
	Type           string
	FunctionalTags map[string]bool // tags related to papa-carlo builder
}
