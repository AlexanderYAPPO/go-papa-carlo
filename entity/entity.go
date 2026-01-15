package entity

type ParsingResult struct {
	Imports []Import
	Fields  []Field
}

type Import struct {
	ReferenceName string
	Path          string
	IsAlias       bool
}

type Field struct {
	Name string
	Type string
}
