package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/immanuelhume/gohm"
)

// gohm
type User struct {
	Name string
	Age  int
}

// gohm
type Meme struct{}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./main.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var v Visitor
	ast.Walk(v, f)
}

type Visitor int

func (v Visitor) Visit(node ast.Node) ast.Visitor {
	t, s := gohm.ParseModel(node)
	if t == nil || s == nil {
		return v
	}
	spew.Dump(t.Name)
	for _, field := range s.Fields.List {
		spew.Dump(field.Names)
		t := field.Type.(*ast.Ident)
		spew.Dump(t.Name)
	}
	return v
}
