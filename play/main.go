package main

import (
	"go/ast"

	"github.com/davecgh/go-spew/spew"
	"github.com/immanuelhume/gohm"
	"golang.org/x/tools/go/packages"
)

// gohm
type User struct {
	Name string
	Age  int
}

// gohm
type Meme struct{}

func main() {
	cfg := &packages.Config{Mode: packages.NeedTypes |
		packages.NeedTypesInfo | packages.NeedFiles | packages.NeedSyntax |
		packages.NeedName}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		var v Visitor
		for _, file := range pkg.Syntax {
			ast.Walk(&v, file)
		}
		for _, ent := range v.Names {
			t := pkg.TypesInfo.Defs[ent]
			spew.Dump(t)
		}
	}
}

type Visitor struct {
	Names []*ast.Ident
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	t, s := gohm.ParseModel(node)
	if t == nil || s == nil {
		return v
	}
	v.Names = append(v.Names, t.Name)
	return v
}
