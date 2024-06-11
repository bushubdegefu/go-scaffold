package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func EchoFrame() {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure
	var data Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// #####################
	// this is where using the data will come to existance

	//  this is creating manger file inside the manager folder
	// ############################################################
	echo_tmpl, err := template.New("data").Parse(devechoTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devecho_file, err := os.Create("manager/devecho.go")
	if err != nil {
		panic(err)
	}
	defer devecho_file.Close()

	err = echo_tmpl.Execute(devecho_file, data)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var devechoTemplate = `
package manager

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"{{.ProjectName}}.com/configs"

	"github.com/spf13/cobra"
)

var (
	{{.AppName}}= &cobra.Command{
		Use:   "echo",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run Gofr development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			echo_run()
		},
	}
)

func echo_run() {

	app := echo.New()
	configs.NewEnvFile("./configs")
	
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Access an environment variable
	apiKey := os.Getenv("TEST_NAME")
	environ := os.Getenv("APP_ENV")
	fmt.Printf(apiKey + ": \t" + environ + "\n")
	// register route greet
	app.GET("/greet", func(ctx echo.Context) error {
		apiKey := os.Getenv("TEST_NAME")
		environ := os.Getenv("APP_ENV")
		fmt.Printf(apiKey + ": \t" + environ + "\n")

		return ctx.String(http.StatusOK, "Hello, World!")
	})

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs

	// configs.NewEnvFile("./configs", )

	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
	// HTTP_PORT := "7500"
	// starting on provided port
	go func(app *echo.Echo) {
		app.Logger.Fatal(app.Start("0.0.0.0:" + HTTP_PORT))
		// log.Fatal(app.ListenTLS(":" + port_1, "server.pem", "server-key.pem"))
	}(app)

	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	fmt.Println("{{.AppName}} was successful shutdown.")

}

func init() {
	goFrame.AddCommand({{.AppName}})

}

`
