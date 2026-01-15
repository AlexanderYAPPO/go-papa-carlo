package generate

import (
	"fmt"
	"strings"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
)

func Generate(result entity.ParsingResult) string {
	codeLines := []string{}
	codeLines = append(codeLines, generateImports(result.Imports)...)

	return strings.Join(codeLines, "\n")
}

func generateImports(imports []entity.Import) []string {
	codeLines := []string{}
	codeLines = append(codeLines, "import (")
	for _, i := range imports {
		newLine := fmt.Sprintf("    %s", i.ReferenceName)
		if i.IsAlias {
			newLine = fmt.Sprintf("    %s \"%s\"", i.ReferenceName, i.Path)
		}
		codeLines = append(codeLines, newLine)
	}
	codeLines = append(codeLines, ")")
	return codeLines
}

// TODO: 
// We should generate the builders from tail to the head.
// So if we have a struct like this:
// type S struct {
// 	A int
// 	B int
// }
// we need to generate the builder for B first, and then for A. It's so because the builder for A will be building a builder for B 
// i.e. it needs to be aware of the builder name for B:

// First, the finalization builder should be built
// type FinalizationBuilder struct {
// 	t *TargetStruct
// }
// func (b FinalizationBuilder).Build() {
// 	TargetStruct
// }

// type BuilderForB struct {
// 	t *TargetStruct
// }
// // on this step we are only aware of the finalization which is a given
// func (b BuilderForB).WithB(bVal B) FinalizationBuilder { 
// 	t.B = bVal
// 	return FinalizationBuilder{
// 		t: t
// 	}
// }

// type BuilderForA struct {
// 	t *TargetStruct
// }
// // on this step, we've already built BuilderForB so it should return this object
// func (b BuilderForB).WithA(aVal A) BuilderForB {
// 	t.A = aVal
// 	return FinalizationBuilder{
// 		t: t
// 	}
// }

// // here we reached the end so we need to generate the New function
// func New() BuilderForA {
// 	var emptyTarget t
// 	return BuilderForA{t: emptyTarget}
// }

// It's another question how to do that nicely. I'm thinking about reverting the list of fields, and creating a new list of pairs that 
// will tell us which objects is expected to be generated and which one we are processing now:

// type S struct {
// 	A int
// 	B int
// 	C int
// }

// Fields: []{("A", "int"), ("B", "int"), ("C", "int")}
// reverse: Fields: []{("C", "int"), ("B", "int"), ("A", "int")}
// use two pointers to build the processing sequence. The result for the first will be a special "symbol". Same for the end

// []{("C", END), ("B", "C"), ("A", "B"), (START, "A")}

// Maybe we shouldn't revert the list but instead just iterate backwards for easier reading

// so the intermediate struct should look something like:
// type FieldGenerationStep struct {
// 	fieldToGenerateBuilderFor optional.Optional[entity.Field] // this field is optional because it is empty for the final element
// 	nextFieldToProduce entity.Field
// 	isLastField bool
// 	isFirstField bool // or these 2 fields can be abstracted behind an interface. That would be even better perhaps.
// }

// I think it's better to have 3 implementations of a single interface:

// type FieldGenerationStep interface {
// 	Generate() []string
// }

// // generates a builder that gives another builder: BuilderForA.WithA -> BuilderForB
// type FieldBuilderGivesFieldBuilder struct {
// 	fieldToGenerateBuilderFor entity.Field
// 	fieldToGive entity.Field
// 	targetType ???
// }

// // generates a builder for the last field and gives a finalization: BuilderForB.WithB -> FinalizationBuilder
// type LastField struct {
// 	field entity.Field
// }

// // builds the target type FinalizationBuilder.Build() -> <Target Type>
// type Finalization struct {
// 	targetType ???
// }

// // starts the chain New() -> BuilderForA
// type MethodNew struct {
// 	targetType ???
// }

// these described types should probably live in an entity package under "generate" pkg.