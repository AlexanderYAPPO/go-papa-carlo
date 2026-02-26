go-papa-carlo
==============

Generate builder-style Go code from a struct definition. The library generates a chained builder with `WithX(...)` methods and a final `Build(). The main difference of this builder library from most other builders for go is that it ensures strong field requirements. 

CLI Usage:
----------

```bash
go-papa-carlo <struct_name> <path_to_struct> [output_path]
```

If `output_path` is omitted, generated code is written next to the input struct file as `<struct_name>_builder_gen.go`.
If `output_path` points to another package directory, generated builders stay in that package and reference the target struct from its original package.

Problem:
--------
In Go, when you need to create a struct, there's no nice way of making required parameters. If you create a struct by initializing an object, Go doesn't require you to specify all fields. Therefore when you add a new field, you cannot make the compiler fail if there are struct initializations that don't specify this field:

```go
type A struct {
    OldField int
    NewField int // imagine this is a new field we are adding
}

// We cannot make NewField a required param. Because of that this addition of the new 
// field won't fail the compildation of the code below
myObj := A{
    OldField: 1, 
}
```

Alternatively, you may create a constructo function New that takes required fields as positional arguments. But for complex structs, the New method will have an interface which is difficult to work with. The more fields there is the harder it gets to read the code:

```go
type A struct {
    Field1 int
    Field2 int
    Field3 int
    Field4 int
    Field5 int
}

func New(Field1 int, Field2 int, Field3 int, Field4 int, Field5 int) A {}

myObj := New(1, 2, 3, 4, 5)
```

Solution:
---------

What this library does is it generates a chain of builder structs in such a way that each builder gives you only one method to call. The struct interfaces are chained down to the point until there's no required fields left. Consider the example below:

```go
type A struct {
    Field1 int
    Field2 int
    Field3 int
}

// After builders are generated, we can create an object with the following code:

myObj := NewABuilder(). // The result of this function gives you a struct with only one method
    WithField1(1). // same for the results of this call. You cannot generate the object until all fields are specified.
    WithField2(2).
    WithField3(3).
    Build() // only this call gives you the final object.
```

Features:
---------

**Omit** - omitting the field such that it won't be mentioned in the builder. Use tag `papa-carlo:omit`

```
type OmittableFields struct {
	RequiredInt    int
	OmmitableString string `papa-carlo:"omit"`
}

got := pkg1.NewOmittableFieldsBuilder().
    WithRequiredInt(11).
    Build()
```

**Optional** - mark the field as optional such that it's not required to be specified when the builder is used. The optional fields can be specified only at the end of the list (before Build() and after all required fields). Use tag `para-carlo:optional`. 
```
type OptionalFields struct {
	RequiredInt    int
	OptString string `papa-carlo:"optional"`
}

got := pkg1.NewStructWithOptionalFieldsBuilder().
    WithRequiredInt(7).
    WithOptionalOptString("hello").
    Build()
```
