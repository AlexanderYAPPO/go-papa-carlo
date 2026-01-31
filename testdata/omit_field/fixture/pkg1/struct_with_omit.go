package pkg1

// StructWithOmittedField ...
type StructWithOmittedField struct {
	RequiredInt    int
	OmmitableString string `papa-carlo:"omit"`
	OmittableStringWithOtherTag string `yaml:"field_to_omit" papa-carlo:"omit" json:"field_to_omit"`
	OmittableFieldWithOtherPapaCarloTag string `para-carl:"fake-tag" papa-carlo:"omit"`
	FieldWithOtherPapaCarloTag string `para-carl:"fake-tag"`
	RequiredBool   bool
}
