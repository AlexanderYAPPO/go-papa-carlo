package target

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	typeName := getTypeName(field.Type)
	tags := parseRelevantTags(field.Tag)
	for _, name := range field.Names {
		// iterate to cover multipla names on a single field:
		// type S struct {
		// 	A, B int
		// }
		fields = append(fields, entity.Field{Name: name.Name, Type: typeName, FunctionalTags: tags})
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

func getTypeName(expr ast.Expr) string {
	var typeName string
	switch t := expr.(type) {
	case *ast.Ident:
		typeName = t.Name

	case *ast.ArrayType:
		typeName = "[]" + getTypeName(t.Elt)

	case *ast.MapType:
		keyType := getTypeName(t.Key)
		valType := getTypeName(t.Value)
		typeName = fmt.Sprintf("map[%s]%s", keyType, valType)

	case *ast.StarExpr:
		typeName = "*" + getTypeName(t.X)

	case *ast.SelectorExpr:
		// This handles pkg.Type or alias.Type
		pkgName := ""
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName = ident.Name
		}
		typeName = pkgName + "." + t.Sel.Name

	case *ast.ChanType:
		direction := "chan"
		switch t.Dir {
		case ast.RECV:
			direction = "<-chan"
		case ast.SEND:
			direction = "chan<-"
		}
		typeName = direction + " " + getTypeName(t.Value)

	case *ast.FuncType:
		params := []string{}
		if t.Params != nil {
			for _, p := range t.Params.List {
				params = append(params, getTypeName(p.Type))
			}
		}

		results := []string{}
		if t.Results != nil {
			for _, r := range t.Results.List {
				results = append(results, getTypeName(r.Type))
			}
		}

		typeName = fmt.Sprintf("func(%s) (%s)",
			strings.Join(params, ", "),
			strings.Join(results, ", "))

	case *ast.InterfaceType:
		return "interface{}"

	}
	return typeName
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
