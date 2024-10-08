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

	// the Otel spanner middleware
	app.Use(otelechospanstarter)

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

func setupRoutes(app *echo.Echo) {
	gapp := app.Group("/api/v1")

	// playgroundHandler := playground.Handler("GraphQL", "/query")

	gapp.POST("/admin", func(ctx echo.Context) error {
		//  Connecting to Databse
		db, err := database.ReturnSession()
		if err != nil {
			log.Errorf("Error Connecting to Database: %v\n", err)
		}
		//  Geting thracer
		tracer := ctx.Get("tracer").(*observe.RouteTracer)

		//  Schema handler
		graphqlHandler := handler.NewDefaultServer(
			graph.NewExecutableSchema(
				graph.Config{Resolvers: &graph.Resolver{DB: db, Tracer: tracer}},
			),
		)

		graphqlHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	// gapp.GET("/playground", func(ctx echo.Context) error {
	// 	playgroundHandler.ServeHTTP(ctx.Response(), ctx.Request())
	// 	return nil
	// })
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

	// the Otel spanner middleware
	app.Use(otelechospanstarter)

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


func setupRoutes(app *echo.Echo) {
	gapp := app.Group("/api/v1")

	// playgroundHandler := playground.Handler("GraphQL", "/query")

	gapp.POST("/admin", func(ctx echo.Context) error {
		//  Connecting to Databse
		db, err := database.ReturnSession()
		if err != nil {
			log.Errorf("Error Connecting to Database: %v\n", err)
		}
		//  Geting thracer
		tracer := ctx.Get("tracer").(*observe.RouteTracer)

		//  Schema handler
		graphqlHandler := handler.NewDefaultServer(
			graph.NewExecutableSchema(
				graph.Config{Resolvers: &graph.Resolver{DB: db, Tracer: tracer}},
			),
		)

		graphqlHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	// gapp.GET("/playground", func(ctx echo.Context) error {
	// 	playgroundHandler.ServeHTTP(ctx.Response(), ctx.Request())
	// 	return nil
	// })
}
`
