package target

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/AlexanderYAPPO/go-papa-carlo/entity"
)

const (
	_papaCarloTagKey = "papa-carlo"
)

func Parse(targetStructName string, pathToFile string) entity.ParsingResult {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, pathToFile, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	imports := captureImports(node)
	allFields := captureFields(node, targetStructName)
	return entity.ParsingResult{
		Imports:     imports,
		Fields:      allFields,
		PackageName: node.Name.Name,
	}
}

func processField(field *ast.Field) []entity.Field {
	fields := []entity.Field{}
	typeInfo := getTypeInfo(field.Type)
	tags := parseRelevantTags(field.Tag)
	for _, name := range field.Names {
		// iterate to cover multipla names on a single field:
		// type S struct {
		// 	A, B int
		// }
		fields = append(fields, entity.Field{
			Name:               name.Name,
			Type:               typeInfo.Name,
			UsesUnexportedType: typeInfo.UsesUnexportedType,
			FunctionalTags:     tags,
		})
	}
	return fields
}

func parseRelevantTags(tag *ast.BasicLit) map[string]bool {
	if tag == nil {
		return nil
	}
	tagLiteral, err := strconv.Unquote(tag.Value)
	if err != nil {
		return nil
	}
	st := reflect.StructTag(tagLiteral)
	raw, ok := st.Lookup(_papaCarloTagKey)
	if !ok || raw == "" {
		return nil
	}
	tags := map[string]bool{}
	for _, tagName := range strings.Split(raw, ",") {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		tags[tagName] = true
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}

type parsedTypeInfo struct {
	Name               string
	UsesUnexportedType bool
}

func getTypeInfo(expr ast.Expr) parsedTypeInfo {
	switch t := expr.(type) {
	case *ast.Ident:
		return parsedTypeInfo{
			Name:               t.Name,
			UsesUnexportedType: !isPredeclaredTypeName(t.Name) && !token.IsExported(t.Name),
		}

	case *ast.ArrayType:
		elt := getTypeInfo(t.Elt)
		return parsedTypeInfo{
			Name:               "[]" + elt.Name,
			UsesUnexportedType: elt.UsesUnexportedType,
		}

	case *ast.MapType:
		keyType := getTypeInfo(t.Key)
		valType := getTypeInfo(t.Value)
		return parsedTypeInfo{
			Name:               fmt.Sprintf("map[%s]%s", keyType.Name, valType.Name),
			UsesUnexportedType: keyType.UsesUnexportedType || valType.UsesUnexportedType,
		}

	case *ast.StarExpr:
		x := getTypeInfo(t.X)
		return parsedTypeInfo{
			Name:               "*" + x.Name,
			UsesUnexportedType: x.UsesUnexportedType,
		}

	case *ast.SelectorExpr:
		// This handles pkg.Type or alias.Type
		pkgName := ""
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName = ident.Name
		}
		return parsedTypeInfo{
			Name:               pkgName + "." + t.Sel.Name,
			UsesUnexportedType: !token.IsExported(t.Sel.Name),
		}

	case *ast.ChanType:
		direction := "chan"
		switch t.Dir {
		case ast.RECV:
			direction = "<-chan"
		case ast.SEND:
			direction = "chan<-"
		}
		valueType := getTypeInfo(t.Value)
		return parsedTypeInfo{
			Name:               direction + " " + valueType.Name,
			UsesUnexportedType: valueType.UsesUnexportedType,
		}

	case *ast.FuncType:
		params := []string{}
		usesUnexportedType := false
		if t.Params != nil {
			for _, p := range t.Params.List {
				paramType := getTypeInfo(p.Type)
				params = append(params, paramType.Name)
				usesUnexportedType = usesUnexportedType || paramType.UsesUnexportedType
			}
		}

		results := []string{}
		if t.Results != nil {
			for _, r := range t.Results.List {
				resultType := getTypeInfo(r.Type)
				results = append(results, resultType.Name)
				usesUnexportedType = usesUnexportedType || resultType.UsesUnexportedType
			}
		}

		return parsedTypeInfo{
			Name: fmt.Sprintf("func(%s) (%s)",
				strings.Join(params, ", "),
				strings.Join(results, ", ")),
			UsesUnexportedType: usesUnexportedType,
		}

	case *ast.InterfaceType:
		return parsedTypeInfo{Name: "interface{}"}

	}
	return parsedTypeInfo{}
}

func isPredeclaredTypeName(name string) bool {
	_, ok := types.Universe.Lookup(name).(*types.TypeName)
	return ok
}

func captureFields(node *ast.File, targetStructName string) []entity.Field {
	allFields := []entity.Field{}
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return false
		}

		if ts.Name.Name != targetStructName {
			return true
		}

		for _, field := range st.Fields.List {
			parsedFlatField := processField(field)
			allFields = append(allFields, parsedFlatField...)
		}
		return true
	})
	return allFields
}

func captureImports(f *ast.File) []entity.Import {
	imports := []entity.Import{}

	for _, imp := range f.Imports {
		fullPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			imports = append(imports, entity.Import{ReferenceName: imp.Name.Name, Path: fullPath, IsAlias: true})
		} else {
			parts := strings.Split(fullPath, "/")
			shortName := parts[len(parts)-1]
			imports = append(imports, entity.Import{ReferenceName: shortName, Path: fullPath, IsAlias: false})
		}
	}
	return imports
}
