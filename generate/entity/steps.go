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
	TargetType     string
	OptionalFields []baseentity.Field
}

func (f Finalization) Generate() []string {
	res := []string{
		"type FinalizationBuilder struct {",
		fmt.Sprintf("\tt *%s", f.TargetType),
		"}",
		"",
	}
	for _, field := range f.OptionalFields {
		methodName := withOptionalMethodName(field)
		paramName := paramName(field)
		res = append(res,
			fmt.Sprintf("func (b FinalizationBuilder) %s(%s %s) FinalizationBuilder {", methodName, paramName, field.Type),
			fmt.Sprintf("\tb.t.%s = %s", field.Name, paramName),
			"\treturn FinalizationBuilder{t: b.t}",
			"}",
			"",
		)
	}
	res = append(res,
		fmt.Sprintf("func (b FinalizationBuilder) Build() %s {", f.TargetType),
		"\treturn *b.t",
		"}",
		"",
	)
	return res
}

type MethodNewToRequiredField struct {
	TargetType string
	FirstField baseentity.Field
}

func (m MethodNewToRequiredField) Generate() []string {
	firstBuilderName := builderName(m.FirstField)
	return []string{
		fmt.Sprintf("func New%sBuilder() %s {", m.TargetType, firstBuilderName),
		fmt.Sprintf("\temptyTarget := &%s{}", m.TargetType),
		fmt.Sprintf("\treturn %s{t: emptyTarget}", firstBuilderName),
		"}",
		"",
	}
}

type MethodNewToFinalization struct {
	TargetType string
}

func (m MethodNewToFinalization) Generate() []string {
	return []string{
		fmt.Sprintf("func New%sBuilder() FinalizationBuilder {", m.TargetType),
		fmt.Sprintf("\temptyTarget := &%s{}", m.TargetType),
		"\treturn FinalizationBuilder{t: emptyTarget}",
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

func withOptionalMethodName(field baseentity.Field) string {
	return fmt.Sprintf("WithOptional%s", field.Name)
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
