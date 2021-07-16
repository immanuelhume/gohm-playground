package gohm_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/immanuelhume/gohm"
)

func TestParseModel(t *testing.T) {
	cases := []struct {
		Name   string
		Node   ast.Node
		Spec   *ast.TypeSpec
		Struct *ast.StructType
	}{
		{"Struct with gohm tag",
			gohmStruct,
			gohmStructTypeSpec,
			gohmStructTypeInfo},
		{"Struct without gohm tag",
			noGohmStruct,
			nil, nil},
		{"Interface with gohm tag",
			interfaceStub,
			nil, nil},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			spec, struc := gohm.ParseModel(test.Node)
			assertDeepEqual(t, spec, test.Spec)
			assertDeepEqual(t, struc, test.Struct)
		})
	}
}

func TestHasGohmTag(t *testing.T) {
	cases := []struct {
		Name    string
		Doc     *ast.CommentGroup
		HasGohm bool
	}{
		{"No comments", &ast.CommentGroup{}, false},
		{"Nil pointer to GenDecl.Doc", nil, false},
		{"Comments without gohm", noGohmDoc, false},
		{"Comments with gohm", gohmDoc, true},
	}

	for _, test := range cases {
		got := gohm.HasGohmTag(test.Doc)
		assertBool(t, got, test.HasGohm)
	}
}

func TestParseStruct(t *testing.T) {
	cases := []struct {
		Name   string
		Node   *ast.GenDecl
		Spec   *ast.TypeSpec
		Struct *ast.StructType
	}{
		{"var declaration",
			&ast.GenDecl{Tok: token.VAR},
			nil, nil},
		{"interface type declaration",
			interfaceStub,
			nil, nil,
		},
		{"struct type declaration",
			gohmStruct,
			gohmStructTypeSpec,
			gohmStructTypeInfo,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			spec, struc := gohm.ParseStruct(test.Node)
			assertDeepEqual(t, spec, test.Spec)
			assertDeepEqual(t, struc, test.Struct)
		})
	}
}

func TestWriteCreate(t *testing.T) {
	b := bytes.Buffer{}
	gohm.WriteCreate(&b, gohmStructTypeSpec, gohmStructTypeInfo)
	want := `func (u *User) Create(name string, age int) {`
	got := b.String()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseFields(t *testing.T) {
	t.Run("places fields in a map", func(t *testing.T) {
		got := gohm.ParseFields(gohmStructTypeInfo.Fields.List)
		want := map[string]string{
			"Name": "string",
			"Age":  "int",
		}
		assertDeepEqual(t, got, want)
	})

	t.Run("empty struct", func(t *testing.T) {
		got := gohm.ParseFields([]*ast.Field{})
		want := map[string]string{}
		assertDeepEqual(t, got, want)
	})
}

// assertion helpers
func assertBool(t testing.TB, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
func assertDeepEqual(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// test doubles
var expr, _ = parser.ParseFile(token.NewFileSet(), "", `
package main

// gohm
// big chungus
type User struct {
	Name string
	Age  int
}

// big chungus
type Song struct {
	Name string
	Length  int
}

// gohm
type Chungus interface {
	Fight()
}
`, parser.ParseComments)

var gohmStruct = expr.Decls[0].(*ast.GenDecl)
var gohmStructTypeSpec = gohmStruct.Specs[0].(*ast.TypeSpec)
var gohmStructTypeInfo = gohmStructTypeSpec.Type.(*ast.StructType)

var noGohmStruct = expr.Decls[1].(*ast.GenDecl)
var noGohmStructTypeSpec = noGohmStruct.Specs[0].(*ast.TypeSpec)
var noGohmStructTypeInfo = noGohmStructTypeSpec.Type.(*ast.StructType)

var gohmDoc = gohmStruct.Doc
var noGohmDoc = noGohmStruct.Doc

var interfaceStub = expr.Decls[2].(*ast.GenDecl)
