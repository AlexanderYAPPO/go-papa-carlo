package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func Parse(targetStructName string, pathToFile string) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, pathToFile, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

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

		fmt.Printf("Found Struct: %s\n", ts.Name.Name)

		fmt.Println("DEBUG: ---------- going through fields ----------------------")
		for _, field := range st.Fields.List {
			processField(field)
		}
		fmt.Println("DEBUG: ---------- ended going through fields ----------------------")
		return true
	})
}

func processField(field *ast.Field) {
	typeName := ""

	// Determine the type of the field
	switch t := field.Type.(type) {
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
			direction = "<-chan" // Receive-only
		case ast.SEND:
			direction = "chan<-" // Send-only
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
	}

	for _, name := range field.Names {
		fmt.Printf("%s %s\n", name.Name, typeName)
	}

}

func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeName(t.X)
	case *ast.ArrayType:
		return "[]" + getTypeName(t.Elt)
	case *ast.SelectorExpr:
		// This handles pkg.Type or alias.Type
		pkgName := ""
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName = ident.Name
		}
		return pkgName + "." + t.Sel.Name
	default:
		return "unknown"
	}
}
