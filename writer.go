package gohm

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"
)

func toReceiver(s string) string {
	return strings.ToLower(string([]rune(s)[0]))
}
func toParams(fields []Field) string {
	_fields := []string{}
	for _, field := range fields {
		_fields = append(_fields,
			fmt.Sprintf("%s %s", strings.ToLower(field.Name), field.Type))
	}
	return strings.Join(_fields, ", ")
}

func GenCreate(w io.Writer, data EntityData) {
	funcs := template.FuncMap{
		"toReceiver": toReceiver,
		"toParams":   toParams,
	}
	tpl, err := template.New("create.go.tpl").
		Funcs(funcs).
		ParseFiles("./templates/create.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
