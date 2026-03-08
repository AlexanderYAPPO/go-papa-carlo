package pkg1

// StructWithOnlyOptional ...
type StructWithOnlyOptional struct {
	Name    string `papa-carlo:"optional"`
	Enabled bool   `json:"enabled" papa-carlo:"optional"`
	Count   int    `papa-carlo:"optional"`
}
