package generate

import (
	baseentity "github.com/AlexanderYAPPO/go-papa-carlo/entity"
	generateentity "github.com/AlexanderYAPPO/go-papa-carlo/generate/entity"
)

const (
	_omitTag     = "omit"
	_optionalTag = "optional"
)

func buildFieldGenerationSteps(target baseentity.Target) []generateentity.FieldGenerationStep {
	targetTypeName := target.Name
	targetTypeRef := target.Reference
	if targetTypeRef == "" {
		targetTypeRef = targetTypeName
	}
	result := target.ParsingResult
	requiredFields, optionalFields := splitFieldsByTags(result.Fields)

	steps := []generateentity.FieldGenerationStep{
		generateentity.Finalization{TargetTypeRef: targetTypeRef, OptionalFields: optionalFields},
	}

	if len(requiredFields) == 0 {
		steps = append(steps, generateentity.MethodNewToFinalization{TargetTypeName: targetTypeName, TargetTypeRef: targetTypeRef})
		return steps
	}

	lastField := requiredFields[len(requiredFields)-1]
	steps = append(steps, generateentity.LastField{Field: lastField, TargetTypeRef: targetTypeRef})

	for i := len(requiredFields) - 2; i >= 0; i-- {
		steps = append(steps, generateentity.FieldBuilderGivesFieldBuilder{
			FieldToGenerateBuilderFor: requiredFields[i],
			FieldToGive:               requiredFields[i+1],
			TargetTypeRef:             targetTypeRef,
		})
	}

	steps = append(steps, generateentity.MethodNewToRequiredField{TargetTypeName: targetTypeName, TargetTypeRef: targetTypeRef, FirstField: requiredFields[0]})
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
