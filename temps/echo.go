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
	"{{.ProjectName}}.com/controllers"
	"github.com/spf13/cobra"

	"github.com/swaggo/echo-swagger"
	_ "{{.ProjectName}}.com/docs"
)

var (
	env string
	{{.AppName}}devechocli= &cobra.Command{
		Use:   "run",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run Gofr development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
		switch env {
		case "":
			echo_run("dev")
		default:
			echo_run(env)
		}
		},
	}
)

func otelechospanstarter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		routeName := ctx.Path() + "_" + strings.ToLower(ctx.Request().Method)
		tracer, span := observe.EchoAppSpanner(ctx, fmt.Sprintf("%v-root", routeName))
		ctx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})

		// Process request
		err := next(ctx)
		if err != nil {
			return err
		}

		span.SetAttributes(attribute.String("response", string(ctx.Response().Status)))
		span.End()
		return nil
	}
}

func dbsessioninjection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		db, err := database.ReturnSession()
		if err != nil {
			return err
		}
		ctx.Set("db", db)

		nerr := next(ctx)
		if nerr != nil {
			return nerr
		}

		return nil
	}
}

func echo_run(env string) {
	//  loading dev env file first
	configs.AppConfig.SetEnv(env)

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

	SetupRoutes(app)
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
	{{.AppName}}devechocli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	goFrame.AddCommand({{.AppName}}devechocli)

}


func SetupRoutes(app *echo.Echo) {
	// the Otel spanner middleware
	app.Use(otelechospanstarter)

	// db session injection
	app.Use(dbsessioninjection)

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
