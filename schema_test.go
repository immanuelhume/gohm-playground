package gohm_test

import (
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
			dogeDecl,
			dogeTypeSpec,
			dogeStructType},
		{"Struct without gohm tag",
			cheemsDecl,
			nil, nil},
		{"Interface with gohm tag",
			interfaceStub,
			nil, nil},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			spec, struc := gohm.ParseModel(test.Node)
			AssertDeepEqual(t, spec, test.Spec)
			AssertDeepEqual(t, struc, test.Struct)
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
		AssertBool(t, got, test.HasGohm)
	}
}

func TestParseTypeSpec(t *testing.T) {
	cases := []struct {
		Name string
		Node *ast.GenDecl
		Want *ast.TypeSpec
	}{
		{"an interface declaration",
			interfaceStub,
			interfaceStubTypeSpec},
		{"a struct declaration",
			dogeDecl,
			dogeTypeSpec},
	}

	for _, test := range cases {
		got := gohm.GetTypeSpec(test.Node)
		AssertDeepEqual(t, got, test.Want)
	}
}

func TestParseStructType(t *testing.T) {
	cases := []struct {
		Name     string
		TypeSpec *ast.TypeSpec
		Want     *ast.StructType
	}{
		{"interface type",
			interfaceStubTypeSpec,
			nil},
		{"struct type",
			dogeTypeSpec,
			dogeStructType},
	}

	for _, test := range cases {
		got := gohm.GetStructType(test.TypeSpec)
		AssertDeepEqual(t, got, test.Want)
	}

}

func TestParseFields(t *testing.T) {
	t.Run("collects fields in a slice", func(t *testing.T) {
		got := gohm.ParseFields(dogeStructType)
		want := []gohm.Field{
			{Name: "Name", Type: "string"},
			{Name: "Age", Type: "int"},
		}
		AssertDeepEqual(t, got, want)
	})

	t.Run("empty struct", func(t *testing.T) {
		got := gohm.ParseFields(emptyStructTypeInfo)
		want := []gohm.Field{}
		AssertDeepEqual(t, got, want)
	})
}

// assertion helpers
func AssertBool(t testing.TB, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
func AssertDeepEqual(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
func AssertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// manually create test doubles
var expr, _ = parser.ParseFile(token.NewFileSet(), "", `
package main

// gohm
// big chungus
type Doge struct {
	Name string
	Age  int
}

// big chungus
type Cheems struct {
	Name string
	Age  int
}

// gohm
type Chungus interface {
	Fight()
}

// empty struct
type Murphy struct {}
`, parser.ParseComments)

var dogeDecl = expr.Decls[0].(*ast.GenDecl)
var dogeTypeSpec = dogeDecl.Specs[0].(*ast.TypeSpec)
var dogeStructType = dogeTypeSpec.Type.(*ast.StructType)

var cheemsDecl = expr.Decls[1].(*ast.GenDecl)

// var cheemsTypeSpec = cheemsDecl.Specs[0].(*ast.TypeSpec)

// var cheemsTypeInfo = cheemsTypeSpec.Type.(*ast.StructType)

var gohmDoc = dogeDecl.Doc
var noGohmDoc = cheemsDecl.Doc

var interfaceStub = expr.Decls[2].(*ast.GenDecl)
var interfaceStubTypeSpec = interfaceStub.Specs[0].(*ast.TypeSpec)

var emptyStruct = expr.Decls[3].(*ast.GenDecl)
var emptyStructTypeSpec = emptyStruct.Specs[0].(*ast.TypeSpec)
var emptyStructTypeInfo = emptyStructTypeSpec.Type.(*ast.StructType)
