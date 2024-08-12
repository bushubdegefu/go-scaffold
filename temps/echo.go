package temps

import (
	"os"
	"text/template"
)

func EchoFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	echo_tmpl, err := template.New("RenderData").Parse(devechoTemplate)
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

	err = echo_tmpl.Execute(devecho_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ##########################################

	prod_tmpl, err := template.New("RenderData").Parse(prodEchoTemplate)
	if err != nil {
		panic(err)
	}

	prod_file, err := os.Create("manager/prodecho.go")
	if err != nil {
		panic(err)
	}
	defer prod_file.Close()

	err = prod_tmpl.Execute(prod_file, RenderData)
	if err != nil {
		panic(err)
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
	"{{.ProjectName}}.com/models/controllers"
	"github.com/spf13/cobra"

	"github.com/swaggo/echo-swagger"
	_ "{{.ProjectName}}.com/docs"
)

var (
	{{.AppName}}devechocli= &cobra.Command{
		Use:   "dev",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run Gofr development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			echo_run()
		},
	}
)

func echo_run() {
	//  loading dev env file first
	configs.AppConfig.SetEnv("dev")

	// starting the app
	app := echo.New()

	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	app.GET("/docs/*", echoSwagger.WrapHandler)

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
	goFrame.AddCommand({{.AppName}}devechocli)

}


func setupRoutes(app *echo.Echo) {
	gapp := app.Group("/admin")
	{{range .Models}}
	gapp.GET("/{{.LowerName}}", controllers.Get{{.Name}}s).Name = "get_all_{{.LowerName}}s"
	gapp.GET("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID).Name = "get_one_{{.LowerName}}s"
	gapp.POST("/{{.LowerName}}", controllers.Post{{.Name}}).Name = "post_{{.LowerName}}"
	gapp.PATCH("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}}).Name = "patch_{{.LowerName}}"
	gapp.DELETE("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name = "delete_{{.LowerName}}"

	{{range .Relations}}
	gapp.POST("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Add{{.FieldName}}{{.ParentName}}s).Name = "add_{{.LowerFieldName}}{{.LowerParentName}}"
	gapp.DELETE("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Delete{{.FieldName}}{{.ParentName}}s).Name = "delete_{{.LowerFieldName}}{{.LowerParentName}}"
	{{end}}
	{{end}}


}

`

var prodEchoTemplate = `
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
	"{{.ProjectName}}.com/models/controllers"
	"github.com/spf13/cobra"

	"github.com/swaggo/echo-swagger"
	_ "{{.ProjectName}}.com/docs"
)

var (
	{{.AppName}}= &cobra.Command{
		Use:   "prod",
		Short: "Run Production Server server ",
		Long:  {{.BackTick}}Run Production server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			prod_echo()
		},
	}
)

func prod_echo()) {
	configs.AppConfig.SetEnv("prod")
	app := echo.New()
	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	app.GET("/docs/*", echoSwagger.WrapHandler)

	setupRoutesEchoProd(app)
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
	goFrame.AddCommand({{.AppName}}prodechocli)

}


func setupRoutesEchoProd(app *echo.Echo) {
	gapp := app.Group("/admin")
	{{range .Models}}
	gapp.GET("/{{.LowerName}}", controllers.Get{{.Name}}s).Name = "get_all_{{.LowerName}}s"
	gapp.GET("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID).Name = "get_one_{{.LowerName}}s"
	gapp.POST("/{{.LowerName}}", controllers.Post{{.Name}}).Name = "post_{{.LowerName}}"
	gapp.PATCH("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}}).Name = "patch_{{.LowerName}}"
	gapp.DELETE("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name = "delete_{{.LowerName}}"

	{{range .Relations}}
	gapp.POST("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Add{{.FieldName}}{{.ParentName}}s).Name = "add_{{.LowerFieldName}}{{.LowerParentName}}"
	gapp.DELETE("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Delete{{.FieldName}}{{.ParentName}}s).Name = "delete_{{.LowerFieldName}}{{.LowerParentName}}"
	{{end}}
	{{end}}


}

`
