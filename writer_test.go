package gohm_test

import (
	"bytes"
	"testing"

	"github.com/immanuelhume/gohm"
)

func TestCreateTemplate(t *testing.T) {
	entityData := gohm.EntityData{
		Package: "main",
		Name:    "User",
		Fields: []gohm.Field{
			{Name: "Name", Type: "string"},
			{Name: "Age", Type: "int"}},
	}
	w := bytes.Buffer{}
	gohm.GenCreate(&w, entityData)

	got := w.String()
	want := `func (u *User) Create(name string, age int) main.User {`
	AssertString(t, got, want)
}
