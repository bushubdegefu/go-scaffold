package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func FiberFrame() {
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
	devf_tmpl, err := template.New("data").Parse(devfTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devf_file, err := os.Create("manager/devfiber.go")
	if err != nil {
		panic(err)
	}
	defer devf_file.Close()

	err = devf_tmpl.Execute(devf_file, data)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var devfTemplate = `
package manager

import (
	"fmt"
	
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"{{.ProjectName}}.com/configs"

	"github.com/spf13/cobra"
)

var (
	{{.AppName}}cli= &cobra.Command{
		Use:   "fiber",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run {{.AppName}} development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			fiber_run()
		},
	}
)


func fiber_run() {
	app := fiber.New()
	//load config file
	configs.NewEnvFile("./configs")
	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	
	// starting on provided port
	go func(app *fiber.App) {
		app.Listen("0.0.0.0:" + HTTP_PORT)
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
	goFrame.AddCommand({{.AppName}}cli)

}
`
