package temps

import (
	"os"
	"text/template"
)

func MainFrame() {
	model := Model{
		Name: "User",
		Fields: []Field{
			{Name: "ID", Type: "int"},
			{Name: "Name", Type: "string"},
			{Name: "Email", Type: "string"},
		},
	}

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("models", os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("models/user.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, model)
	if err != nil {
		panic(err)
	}
}

var modelTemplate = `
package models

type {{.Name}} struct {
	{{range .Fields}}
	{{.Name}} {{.Type}}
	{{end}}
}
`
