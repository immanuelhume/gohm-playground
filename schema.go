package gohm

import (
	"fmt"
	"go/ast"
	"io"
	"strings"
)

func ParseModel(n ast.Node) (*ast.TypeSpec, *ast.StructType) {
	d, ok := n.(*ast.GenDecl)
	if !ok {
		return nil, nil // it's not a GenDecl type
	}
	if !HasGohmTag(d.Doc) {
		return nil, nil // not tagged with // gohm
	}
	t, s := ParseStruct(d)
	if t == nil || s == nil {
		return nil, nil // it's not a struct type
	}
	return t, s
}

func HasGohmTag(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.TrimSpace(comment.Text) == "// gohm" {
			return true
		}
	}
	return false
}

func ParseStruct(d *ast.GenDecl) (*ast.TypeSpec, *ast.StructType) {
	if d == nil {
		return nil, nil
	}
	if len(d.Specs) == 0 {
		return nil, nil
	}
	t, ok := d.Specs[0].(*ast.TypeSpec) // cast to TypeSpec
	if !ok {
		return nil, nil // not a type declaration
	}
	s, ok := t.Type.(*ast.StructType) // check if it's a struct
	if !ok {
		return nil, nil
	}
	return t, s
}

func WriteCreate(w io.Writer, t *ast.TypeSpec, s *ast.StructType) {
	entity := t.Name.Name
	fields := ParseFields(s.Fields.List)
	_fields := []string{}
	for k, v := range fields {
		_fields = append(_fields, fmt.Sprintf("%s %s", strings.ToLower(k), v))
	}
	fmt.Fprintf(w, `func (%s *%s) Create(%s) {`,
		strings.ToLower(string(entity[0])), entity, strings.Join(_fields, ", "))
}

func ParseFields(fields []*ast.Field) map[string]string {
	fmap := make(map[string]string)
	for _, field := range fields {
		_name := field.Names[0].Name
		_type := field.Type.(*ast.Ident).Name
		fmap[_name] = _type
	}
	return fmap
}
