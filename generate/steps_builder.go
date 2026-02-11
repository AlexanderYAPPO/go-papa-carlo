package generate

import (
	baseentity "github.com/AlexanderYAPPO/go-papa-carlo/entity"
	generateentity "github.com/AlexanderYAPPO/go-papa-carlo/generate/entity"
)

const (
	_omitTag     = "omit"
	_optionalTag = "optional"
)

func buildFieldGenerationSteps(result baseentity.ParsingResult) []generateentity.FieldGenerationStep {
	targetType := result.TargetTypeName
	requiredFields, optionalFields := splitFieldsByTags(result.Fields)

	steps := []generateentity.FieldGenerationStep{
		generateentity.Finalization{TargetType: targetType, OptionalFields: optionalFields},
	}

	if len(requiredFields) == 0 {
		steps = append(steps, generateentity.MethodNewToFinalization{TargetType: targetType})
		return steps
	}

	lastField := requiredFields[len(requiredFields)-1]
	steps = append(steps, generateentity.LastField{Field: lastField, TargetType: targetType})

	for i := len(requiredFields) - 2; i >= 0; i-- {
		steps = append(steps, generateentity.FieldBuilderGivesFieldBuilder{
			FieldToGenerateBuilderFor: requiredFields[i],
			FieldToGive:               requiredFields[i+1],
			TargetType:                targetType,
		})
	}

	steps = append(steps, generateentity.MethodNewToRequiredField{TargetType: targetType, FirstField: requiredFields[0]})
	return steps
}

func splitFieldsByTags(fields []baseentity.Field) (requiredFields []baseentity.Field, optionalFields []baseentity.Field) {
	requiredFields = []baseentity.Field{}
	optionalFields = []baseentity.Field{}

	for _, field := range fields {
		if field.FunctionalTags[_omitTag] {
			continue
		}
		if field.FunctionalTags[_optionalTag] {
			optionalFields = append(optionalFields, field)
			continue
		}
		requiredFields = append(requiredFields, field)
	}

	return requiredFields, optionalFields
}
