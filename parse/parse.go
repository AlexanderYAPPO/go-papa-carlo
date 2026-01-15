package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"maps"
	"strings"
)

func Parse(targetStructName string, pathToFile string) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, pathToFile, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	imports := captureImports(node)
	fmt.Println("DEBUG: ---------- imports ----------------------")
	for alias, path := range imports {
		fmt.Printf("%s -> %s\n", alias, path)
	}
	fmt.Println("DEBUG: ---------- imports ----------------------")

	fmt.Println("")

	allFieldsMap := captureFields(node, targetStructName)
	fmt.Println("DEBUG: ---------- all fields map ----------------------")
	for field, typeName := range allFieldsMap {
		fmt.Printf("%s -> %s\n", field, typeName)
	}
	fmt.Println("DEBUG: ---------- all fields map ----------------------")
}

func processField(field *ast.Field) map[string]string {
	fieldMap := make(map[string]string)
	typeName := getTypeName(field.Type)
	for _, name := range field.Names {
		fieldMap[name.Name] = typeName
	}
	return fieldMap
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

func captureFields(node *ast.File, targetStructName string) map[string]string {
	allFieldsMap := make(map[string]string) // field_name -> type_name
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
			fieldMap := processField(field)
			maps.Copy(allFieldsMap, fieldMap)
		}
		return true
	})
	return allFieldsMap
}

func captureImports(f *ast.File) map[string]string {
	importMap := make(map[string]string)

	for _, imp := range f.Imports {
		fullPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			importMap[imp.Name.Name] = fullPath
		} else {
			parts := strings.Split(fullPath, "/")
			shortName := parts[len(parts)-1]
			importMap[shortName] = fullPath
		}
	}
	return importMap
}
