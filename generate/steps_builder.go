package generate

import (
	baseentity "github.com/AlexanderYAPPO/go-papa-carlo/entity"
	generateentity "github.com/AlexanderYAPPO/go-papa-carlo/generate/entity"
)

func buildFieldGenerationSteps(result baseentity.ParsingResult) []generateentity.FieldGenerationStep {
	targetType := result.TargetTypeName
	steps := []generateentity.FieldGenerationStep{
		generateentity.Finalization{TargetType: targetType},
	}

	fields := result.Fields
	if len(fields) == 0 {
		return steps
	}

	lastField := fields[len(fields)-1]
	steps = append(steps, generateentity.LastField{Field: lastField, TargetType: targetType})

	for i := len(fields) - 2; i >= 0; i-- {
		steps = append(steps, generateentity.FieldBuilderGivesFieldBuilder{
			FieldToGenerateBuilderFor: fields[i],
			FieldToGive:               fields[i+1],
			TargetType:                targetType,
		})
	}

	steps = append(steps, generateentity.MethodNew{TargetType: targetType, FirstField: fields[0]})
	return steps
}
