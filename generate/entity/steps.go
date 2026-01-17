package entity

import (
	"fmt"
	"unicode"

	baseentity "github.com/AlexanderYAPPO/go-papa-carlo/entity"
)

type FieldGenerationStep interface {
	Generate() []string
}

type FieldBuilderGivesFieldBuilder struct {
	FieldToGenerateBuilderFor baseentity.Field
	FieldToGive               baseentity.Field
	TargetType                string
}

func (f FieldBuilderGivesFieldBuilder) Generate() []string {
	currentBuilderName := builderName(f.FieldToGenerateBuilderFor)
	nextBuilderName := builderName(f.FieldToGive)
	methodName := withMethodName(f.FieldToGenerateBuilderFor)
	paramName := paramName(f.FieldToGenerateBuilderFor)

	return []string{
		fmt.Sprintf("type %s struct {", currentBuilderName),
		fmt.Sprintf("\tt *%s", f.TargetType),
		"}",
		"",
		fmt.Sprintf("func (b %s) %s(%s %s) %s {", currentBuilderName, methodName, paramName, f.FieldToGenerateBuilderFor.Type, nextBuilderName),
		fmt.Sprintf("\tb.t.%s = %s", f.FieldToGenerateBuilderFor.Name, paramName),
		fmt.Sprintf("\treturn %s{t: b.t}", nextBuilderName),
		"}",
		"",
	}
}

type LastField struct {
	Field      baseentity.Field
	TargetType string
}

func (l LastField) Generate() []string {
	currentBuilderName := builderName(l.Field)
	methodName := withMethodName(l.Field)
	paramName := paramName(l.Field)

	return []string{
		fmt.Sprintf("type %s struct {", currentBuilderName),
		fmt.Sprintf("\tt *%s", l.TargetType),
		"}",
		"",
		fmt.Sprintf("func (b %s) %s(%s %s) FinalizationBuilder {", currentBuilderName, methodName, paramName, l.Field.Type),
		fmt.Sprintf("\tb.t.%s = %s", l.Field.Name, paramName),
		"\treturn FinalizationBuilder{t: b.t}",
		"}",
		"",
	}
}

type Finalization struct {
	TargetType string
}

func (f Finalization) Generate() []string {
	return []string{
		"type FinalizationBuilder struct {",
		fmt.Sprintf("\tt *%s", f.TargetType),
		"}",
		"",
		fmt.Sprintf("func (b FinalizationBuilder) Build() %s {", f.TargetType),
		"\treturn *b.t",
		"}",
		"",
	}
}

type MethodNew struct {
	TargetType string
	FirstField baseentity.Field
}

func (m MethodNew) Generate() []string {
	firstBuilderName := builderName(m.FirstField)
	return []string{
		fmt.Sprintf("func New%sBuilder() %s {", m.TargetType, firstBuilderName),
		fmt.Sprintf("\temptyTarget := &%s{}", m.TargetType),
		fmt.Sprintf("\treturn %s{t: emptyTarget}", firstBuilderName),
		"}",
		"",
	}
}

func builderName(field baseentity.Field) string {
	return fmt.Sprintf("BuilderFor%s", field.Name)
}

func withMethodName(field baseentity.Field) string {
	return fmt.Sprintf("With%s", field.Name)
}

func paramName(field baseentity.Field) string {
	base := lowerFirst(field.Name)
	if base == "" {
		return "value"
	}
	return fmt.Sprintf("%sVal", base)
}

func lowerFirst(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
