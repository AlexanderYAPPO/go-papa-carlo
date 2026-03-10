package pkg

type StructWithUnexportedType struct {
	PublicField privateType
}

type privateType struct {
	Value string
}
