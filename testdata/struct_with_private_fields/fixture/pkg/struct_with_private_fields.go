package pkg

// StructWithPrivateFieldsThatErrors is a struct that should error because the private field is not omitted.
type StructWithPrivateFieldsThatErrors struct {
	PublicField int
	privateField int
}

// StructWithPrivateFieldsThatWorks is a struct that should work because the private field is omitted.
type StructWithPrivateFieldsThatWorks struct {
	PublicField int
	privateField int `papa-carlo:"omit"`
}

