package pkg1

// StructWithOptionalField ...
type StructWithOptionalField struct {
	RequiredInt    int
	OptionalString string `papa-carlo:"optional"`
	OptionalBool   bool   `json:"optional_bool" papa-carlo:"optional"`
	RequiredMap    map[string]int
}
