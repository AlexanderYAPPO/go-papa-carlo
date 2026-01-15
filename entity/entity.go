package entity

type ParsingResult struct {
	Imports []Import
	Fields  []Field
}

type Import struct {
	Alias string
	Path  string
}

type Field struct {
	Name string
	Type string
}
