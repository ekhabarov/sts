package sts

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Represents field tags with values.
	Tags map[string]string

	// Field contains information about one structure field without field name.
	Field struct {
		// internal type info.
		Type types.Type
		// Equals true if field is a pointer.
		IsPointer bool
		// Field tags with values.
		Tags Tags
		// Order inside structure. It's used for printing fields with correct order.
		Ord uint8
	}

	// Fields is a set of fields of one structure.
	// Key is field name.
	Fields map[string]Field

	// File contains structures from parsed file.
	File struct {
		Package string
		Structs map[string]Fields
	}
)

// String return structure information as a string.
func (s Fields) String() string {
	c := ""
	for k, v := range s {
		c += fmt.Sprintf(
			"// Field: %q, Type: %q, isPointer: %t, Tags: %v\n",
			k, v.Type, v.IsPointer, v.Tags,
		)
	}
	c += "\n"
	return c
}

// baseType returns type name in format <package>.<type> for FQTN like
// github.com/ekhabarov/sts/examples/nulls.Time.
func baseType(t types.Type) string {
	switch typ := t.(type) {
	case *types.Named:
		splt := strings.Split(typ.String(), "/")

		switch true {
		case len(splt) > 1: // it's a FQTN
			return splt[len(splt)-1]
		case strings.Contains(typ.String(), "."):
			return fmt.Sprintf("%s.%s", typ.Obj().Pkg().Name(), strings.Split(typ.String(), ".")[1])
		default:
			return fmt.Sprintf("%s.%s", typ.Obj().Pkg().Name(), typ)
		}

	case *types.Basic:
		return typ.String()

	default:
		return typ.String()
	}
}

func (fi Field) String() string {
	s := strings.Split(fi.Type.String(), "/")
	t := s[len(s)-1]

	if fi.IsPointer {
		return "*" + t
	}
	return t
}

// inspect is a function which is run for each node in source file. See go/ast
// package for details.
func inspect(
	output *File, info *types.Info, tags []string,
) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		if p, ok := n.(*ast.File); ok {
			output.Package = p.Name.Name
			return true
		}

		spec, ok := n.(*ast.TypeSpec)
		if !ok || spec.Type == nil { // skip non-types and empty types
			return true
		}

		s, ok := spec.Type.(*ast.StructType)
		if !ok { // skip non-struct types
			return true
		}

		sname := spec.Name.String()
		if output.Structs == nil {
			output.Structs = map[string]Fields{}
		}

		if _, ok := output.Structs[sname]; !ok {
			output.Structs[sname] = Fields{}
		}

		embeddedCounter := 0
		ord := uint8(0)
		for _, field := range s.Fields.List {
			fname := "embedded_"
			// Embedded structures have no names.
			if field.Names != nil {
				fname = field.Names[0].Name
			} else {
				fname += strconv.Itoa(embeddedCounter)
				embeddedCounter++
			}

			var ftags Tags

			if t := field.Tag; t != nil {
				ftags = fieldTags(t.Value, tags)
			}

			id, fn, _, ptr := typsw(field.Type)
			if id == nil {
				continue
			}

			if fn != "" {
				fname = fn
			}

			output.Structs[sname][fname] = Field{
				Type:      info.TypeOf(id),
				IsPointer: ptr,
				Tags:      ftags,
				Ord:       ord,
			}
			ord++
		}
		return false
	}
}

func fieldTags(tagValue string, list []string) Tags {
	tags := Tags{}

	rawtag := strings.Replace(tagValue, "`", "", -1)

	for _, t := range list {
		if v := reflect.StructTag(rawtag).Get(t); v != "" {
			tags[t] = v
		}
	}

	return tags
}

// typsw is a recursive type switch which is used by inspect.
func typsw(fieldType ast.Expr) (id *ast.Ident, fname, typ string, ptr bool) {
	switch t := fieldType.(type) {
	case *ast.Ident: // simple types e.g. int, string, etc.
		id = t
		typ = t.Name
		return

	case *ast.SelectorExpr: // types like time.Time, time.Duration, nulls.String
		id = t.Sel
		typ = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
		return

	case *ast.StarExpr: // pointer to something
		id, fname, typ, _ = typsw(t.X)
		ptr = true
		return

	case *ast.ArrayType:
		id, fname, typ, ptr = typsw(t.Elt)
		return
	}
	return
}

// Parse gets path to source file, imports whole file package and run inspect
// functions on given file. Parsing whole package is necessary because
// structures and field types can be defined in different files. Function
// returns list of structures with information their about fields.
func Parse(path string, tags []string) (*File, error) {
	fset := token.NewFileSet()
	file := filepath.Base(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %q doesn't exist", path)
	}

	// Parse whole package.
	pkgs, err := parser.ParseDir(fset, filepath.Dir(path),
		func(f os.FileInfo) bool {
			return !strings.HasSuffix(f.Name(), "_test.go") // skip test files.
		}, 0)
	if err != nil {
		return nil, err
	}

	// AST representation of the file passed via path parameter.
	var node *ast.File
	ok := false
	astFiles := []*ast.File{} // parsed package files

	for pn, p := range pkgs {
		for _, f := range p.Files {
			// collect all files from within a package
			astFiles = append(astFiles, f)
		}

		// determine passed file.
		for n := range p.Files {
			if node != nil {
				break
			}

			// Possible BUG: if there are two files with same name from different
			// packages and it sts is run from parent directory.
			if !strings.HasSuffix(n, file) {
				continue
			}

			name := file
			if strings.Contains(n, "/") { // filename can be of format package/file.go
				name = strings.TrimSuffix(n, file) + file
			}

			node, ok = p.Files[name]
			if !ok {
				return nil, fmt.Errorf("file %q not found in package %q: ", file, pn)
			}
		}
	}

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	// Important info about importer
	// https://github.com/golang/go/issues/11415#issuecomment-283445198
	//
	// Basically, importer.Default() doesn't work when package like
	// "github.com/ekhabarov/sts/example/nulls" is imported.
	//
	// TODO(ekhabarov): import structs from vendor.
	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}

	_, err = conf.Check("", fset, astFiles, &info)
	if err != nil {
		return nil, fmt.Errorf("conf.check call error: %#v", err)
	}

	f := &File{}

	ast.Inspect(node, inspect(f, &info, tags))

	return f, nil
}

// Lookup return structure by name from parsed source file or an error if
// structure with such name not found.
func Lookup(f *File, name string) (Fields, error) {
	flds, ok := f.Structs[name]
	if !ok {
		return flds, fmt.Errorf("structure %q not found", name)
	}

	return flds, nil
}
