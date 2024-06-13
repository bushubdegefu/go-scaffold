package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	// setting default value for config data file
	//  GetPostPatchPut
	// "Get$Post$Patch$Put"

	for i := 0; i < len(data.Models); i++ {
		data.Models[i].LowerName = strings.ToLower(data.Models[i].Name)
		data.Models[i].AppName = data.AppName
		data.Models[i].ProjectName = data.ProjectName

		for j := 0; j < len(data.Models[i].Fields); j++ {
			data.Models[i].Fields[j].BackTick = "`"
			cf := strings.Split(data.Models[i].Fields[j].CurdFlag, "$")

			data.Models[i].Fields[j].Get, _ = strconv.ParseBool(cf[0])
			data.Models[i].Fields[j].Post, _ = strconv.ParseBool(cf[1])
			data.Models[i].Fields[j].Patch, _ = strconv.ParseBool(cf[2])
			data.Models[i].Fields[j].Put, _ = strconv.ParseBool(cf[3])
			data.Models[i].Fields[j].AppName = data.AppName
			data.Models[i].Fields[j].ProjectName = data.ProjectName

		}
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


	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/models/controllers"
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

	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	setupRoutes(app)
	// starting on provided port
	go func(app *echo.Echo) {
		//  Http serving port
		HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
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


func setupRoutes(app *echo.Echo) {
	gapp := app.Group("/admin")
	{{range .Models}}
	gapp.GET("/{{.LowerName}}", controllers.Get{{.Name}}s).Name = "get_all_{{.LowerName}}s"
	gapp.GET("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID).Name = "get_one_{{.LowerName}}s"
	gapp.POST("/{{.LowerName}}", controllers.Post{{.Name}}).Name = "post_{{.LowerName}}"
	gapp.PATCH("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}}).Name = "patch_{{.LowerName}}"
	gapp.DELETE("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name = "delete_{{.LowerName}}"

	{{end}}
}

`
