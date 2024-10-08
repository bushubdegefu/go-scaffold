package temps

import (
	"fmt"
	"os"
	"text/template"
)

func FiberFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	devf_tmpl, err := template.New("RenderData").Parse(devfTemplate)
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
		fmt.Println(err)
		panic(err)
	}
	defer devf_file.Close()

	err = devf_tmpl.Execute(devf_file, RenderData)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

// https://help.sumologic.com/docs/apm/traces/get-started-transaction-tracing/opentelemetry-instrumentation/go/

var devfTemplate = `
package manager

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"fmt"
	"log"
	"strconv"
	"time"


	"os"
	"os/signal"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel/attribute"
	"github.com/gofiber/swagger"
	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/bluetasks"
	"{{.ProjectName}}.com/observe"
	"{{.ProjectName}}.com/controllers"
	_ "{{.ProjectName}}.com/docs"
	"{{.ProjectName}}.com/database"
	"github.com/spf13/cobra"
)

var (
	env string
	{{.AppName}}cli= &cobra.Command{
		Use:   "run",
		Short: "Run Development or Production server Based on Provided --env=dev flag. Defaults to dev ",
		Long:  {{.BackTick}}Run {{.AppName}} development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			switch env {
			case "":
				fiber_run("dev")
			default:
				fiber_run(env)
			}
		},
	}
)

func otelspanstarter(ctx *fiber.Ctx) error {
	//  creating trace context from span if they exist
	route_name := ctx.Path() + "_" + strings.ToLower(ctx.Route().Method)
	tracer, span := observe.FiberAppSpanner(ctx, fmt.Sprintf("%v-root", route_name))
	ctx.Locals("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
	if err := ctx.Next(); err != nil {
		return err
	}
	span.SetAttributes(attribute.String("response", ctx.Response().String()))
	span.End()
	return nil
}

func dbsessioninjection(ctx *fiber.Ctx) error {
	db, err := database.ReturnSession()
	if err != nil {
		return ctx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	ctx.Locals("db", db)
	return ctx.Next()
}

func NextFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

func fiber_run(env string) {
	prefork := false
	if env == "prod" {
		prefork = true
	}

	//  Loading Configuration
	configs.AppConfig.SetEnv(env)


		// Starting Otel Global tracer
		tp := observe.InitTracer()
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()

		// Starting Task Scheduler ( Running task that run regularly based on the provided configs)
		schd := bluetasks.ScheduledTasks()
		defer schd.Stop()


		// Basic App Configs
		body_limit, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BODY_LIMIT", "70"))
		read_buffer_size, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("READ_BUFFER_SIZE", "70"))
		rate_limit_per_second, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "5000"))
		//load config file
		app := fiber.New(fiber.Config{
			Prefork: prefork,
			// Network:     fiber.NetworkTCP,
			// Immutable:   true,
			JSONEncoder:    json.Marshal,
			JSONDecoder:    json.Unmarshal,
			BodyLimit:      body_limit * 1024 * 1024,
			ReadBufferSize: read_buffer_size * 1024,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				// Status code defaults to 500
				code := fiber.StatusInternalServerError
				// Retrieve the custom status code if it's a *fiber.Error
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}
				// Send custom error page
				err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
				if err != nil {
					// In case the SendFile fails
					return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
				}
				// Return from handler
				return nil
			},
		})

		// recover from panic attacks middlerware
		app.Use(recover.New())

		// allow cross origin request
		app.Use(cors.New())

		//  rate limiting middleware
		app.Use(limiter.New(limiter.Config{
			Max:               rate_limit_per_second,
			Expiration:        1 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))

		// idempotency middleware
		app.Use(idempotency.New(idempotency.Config{
			Lifetime: 10 * time.Second,
				}))


		// logger middle ware with the custom file writer object
		app.Use(logger.New(logger.Config{
			Format:     "\n${cyan}-[${time}]-[${ip}] -${white}${pid} ${red}${status} ${blue}[${method}] ${white}-${path}\n [${body}]\n[${error}]\n[${resBody}]\n[${reqHeaders}]\n[${queryParams}]\n",
			TimeFormat: "15:04:05",
			TimeZone:   "Local",
			Output:     os.Stdout,
		}))

		// prometheus middleware concrete instance
		prometheus := fiberprometheus.New("gobluefiber")
		prometheus.RegisterAt(app, "/metrics")

		// prometheus monitoring middleware
		app.Use(prometheus.Middleware)


		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})
		// swagger docs
		app.Get("/docs/*", swagger.HandlerDefault)
		app.Get("/docs/*", swagger.New()).Name("swagger_routes")

		// fiber native monitoring metrics endpoint
		app.Get("/lmetrics", monitor.New(monitor.Config{Title: "goBlue Metrics Page"})).Name("custom_metrics_route")


		HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
		// starting on provided port
		go func(app *fiber.App) {
			app.Listen("0.0.0.0:" + HTTP_PORT)
			}(app)

			c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	fmt.Println("{{.AppName}} was successful shutdown.")
}


func init() {
	{{.AppName}}cli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	goFrame.AddCommand({{.AppName}}cli)

}


func SetupRoutes(app *fiber.App) {

	//app logging open telemetery
	app.Use(otelfiber.Middleware())
	app.Use(otelspanstarter)

	// database session injection to local context
	app.Use(dbsessioninjection)

	gapp := app.Group("/api/v1")

	{{range .Models}}
	gapp.Get("/{{.LowerName}}",NextFunc).Name("get_all_{{.LowerName}}s").Get("/{{.LowerName}}", controllers.Get{{.Name}}s)
	gapp.Get("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("get_one_{{.LowerName}}s").Get("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID)
	gapp.Post("/{{.LowerName}}",NextFunc).Name("post_{{.LowerName}}").Post("/{{.LowerName}}", controllers.Post{{.Name}})
	gapp.Patch("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("patch_{{.LowerName}}").Patch("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}})
	gapp.Delete("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("delete_{{.LowerName}}").Delete("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name("delete_{{.LowerName}}")
	{{range .Relations}}{{if .OtM}}
	gapp.Patch("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",NextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Patch("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",NextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s){{end}}
	{{if .MtM}}gapp.Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s){{end}}{{end}}
	{{end}}
}

`
