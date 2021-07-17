package gohm

import (
	"go/ast"
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
	t := GetTypeSpec(d)
	if t == nil {
		return nil, nil
	}
	s := GetStructType(t)
	if s == nil {
		return nil, nil
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

func GetTypeSpec(d *ast.GenDecl) *ast.TypeSpec {
	if d == nil {
		return nil
	}
	if len(d.Specs) == 0 {
		return nil
	}
	t, ok := d.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	return t
}

func GetStructType(t *ast.TypeSpec) *ast.StructType {
	s, ok := t.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	return s
}

type Field struct {
	Name string
	Type string
}

func ParseFields(thing *ast.StructType) []Field {
	fields := []Field{}
	_fields := thing.Fields.List
	if _fields == nil {
		return fields
	}
	for _, field := range _fields {
		fields = append(fields, Field{
			Name: field.Names[0].Name,
			Type: field.Type.(*ast.Ident).Name})
	}
	return fields
}

type EntityData struct {
	Package string
	Name    string
	Fields  []Field
}
