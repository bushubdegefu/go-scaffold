package temps

import (
	"os"
	"text/template"
)

func GraphEchoFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	echo_tmpl, err := template.New("RenderData").Parse(deveGraphEchoTemplate)
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

	prod_tmpl, err := template.New("RenderData").Parse(prodGraphEchoTemplate)
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

var deveGraphEchoTemplate = `
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
	"{{.ProjectName}}.com/bluetasks"
	"github.com/spf13/cobra"

	"github.com/swaggo/echo-swagger"
)

var (
	{{.AppName}}graphdevechocli= &cobra.Command{
		Use:   "gdev",
		Short: "Run GraphQL Echo Development server ",
		Long:  {{.BackTick}}Run Gofr development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			graph_echo_run()
		},
	}
)

func graph_echo_run() {
	//  loading dev env file first
	configs.AppConfig.SetEnv("dev")

	// Starting Otel Global tracer
	tp := observe.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// starting the app
	app := echo.New()


	// Starting Task Scheduler ( Running task that run regularly based on the provided configs)
	schd := bluetasks.ScheduledTasks()
	defer schd.Stop()

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
	goFrame.AddCommand({{.AppName}}graphdevechocli)
}


func setupRoutes(gapp *echo.Echo) {
	gapp := app.Group("/api/v1")

	db, err := database.ReturnSession()
	if err != nil {
		panic("Error Connecting to Database")
	}

	playgroundHandler := playground.Handler("GraphQL", "/query")

	gapp.POST("/admin", func(c echo.Context) error {
		tracer := ctx.Get("tracer").(*observe.RouteTracer)
		graphqlHandler := handler.NewDefaultServer(
			graph.NewExecutableSchema(
				graph.Config{Resolvers: &graph.Resolver{DB: db, Tracer: tracer}},
			),
		)
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	gapp.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

`

var prodGraphEchoTemplate = `
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
	"{{.ProjectName}}.com/bluetasks"
	"github.com/spf13/cobra"

	"github.com/swaggo/echo-swagger"
)

var (
	{{.AppName}}graphprodechocli= &cobra.Command{
		Use:   "gprod",
		Short: "Run GraphQL Echo Production Server server ",
		Long:  {{.BackTick}}Run Production server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			graph_prod_echo()
		},
	}
)

func graph_prod_echo() {
	//  loading dev env file first
	configs.AppConfig.SetEnv("prod")

	// Starting Otel Global tracer
	tp := observe.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// starting the app
	app := echo.New()

	// Starting Task Scheduler ( Running task that run regularly based on the provided configs)
	schd := bluetasks.ScheduledTasks()
	defer schd.Stop()

	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

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
	goFrame.AddCommand({{.AppName}}graphprodechocli)
}


func setupRoutesEchoProd(app *echo.Echo) {
	gapp := app.Group("/api/v1")

	db, err := database.ReturnSession()
	if err != nil {
		panic("Error Connecting to Database")
	}

	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{DB: db}},
		),
	)
	playgroundHandler := playground.Handler("GraphQL", "/query")

	gapp.POST("/admin", func(c echo.Context) error {
		tracer := ctx.Get("tracer").(*observe.RouteTracer)
		graphqlHandler := handler.NewDefaultServer(
			graph.NewExecutableSchema(
				graph.Config{Resolvers: &graph.Resolver{DB: db, Tracer: tracer}},
			),
		)
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	gapp.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})


}

`
