package entity

type ParsingResult struct {
	Imports []Import
	Fields  []Field
	TargetTypeName string // TODO: not populated
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
