package temps

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func GitDockerFrame() {

	// ############################################################
	docker_tmpl, err := template.New("RenderData").Parse(dockerConfig)
	if err != nil {
		panic(err)
	}

	docker_file, err := os.Create("api.Dockerfile")
	if err != nil {
		panic(err)
	}
	defer docker_file.Close()

	err = docker_tmpl.Execute(docker_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ############################################################
	git_tmpl, err := template.New("RenderData").Parse(gitIgnore)
	if err != nil {
		panic(err)
	}

	git_file, err := os.Create(".gitignore")
	if err != nil {
		panic(err)
	}
	defer git_file.Close()

	err = git_tmpl.Execute(git_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ############################################################
	dockig_tmpl, err := template.New("RenderData").Parse(dockerIgnore)
	if err != nil {
		panic(err)
	}

	dockig_file, err := os.Create(".dockerignore")
	if err != nil {
		panic(err)
	}
	defer dockig_file.Close()

	err = dockig_tmpl.Execute(dockig_file, RenderData)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var dockerIgnore = `

`
var dockerConfig = `
FROM golang:latest

RUN apt install -y libc6 libc-bin

RUN apt -y update && apt -y upgrade

RUN apt -y install build-essential pkg-config g++ git cmake yasm

RUN apt install build-essential pkg-config git

WORKDIR /playground/

COPY docs /playground/

COPY main /playground/

COPY server.pem  /playground/

COPY server-key.pem  /playground/ 

COPY configs /playground/

RUN chmod +x main

CMD ["./main","prod"]
 `
var gitIgnore = `
configs/*
*.logs
tests/.env
tests/.test.env
`
