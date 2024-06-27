package temps

import (
	"os"
	"text/template"
)

func RabbitFrame() {

	// ####################################################
	//  rabbit template
	rab_tmpl, err := template.New("RenderData").Parse(rabbitTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("messages", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rab_file, err := os.Create("messages/connection.go")
	if err != nil {
		panic(err)
	}
	defer rab_file.Close()

	err = rab_tmpl.Execute(rab_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var rabbitTemplate = `

`
