package temps

import (
	"fmt"
	"os"
	"os/exec"
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
		panic(err)
	}
	defer devf_file.Close()

	err = devf_tmpl.Execute(devf_file, RenderData)
	if err != nil {
		panic(err)
	}

	// #################################################
	prod_tmpl, err := template.New("RenderData").Parse(prodTemplate)
	if err != nil {
		panic(err)
	}

	prod_file, err := os.Create("manager/prodfiber.go")
	if err != nil {
		panic(err)
	}
	defer prod_file.Close()

	err = prod_tmpl.Execute(prod_file, RenderData)
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
	"context"
	"encoding/json"
	"errors"
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
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/observe"
	"{{.ProjectName}}.com/models/controllers"
	_ "{{.ProjectName}}.com/docs"
	"github.com/spf13/cobra"
)

var (
	{{.AppName}}cli= &cobra.Command{
		Use:   "dev",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run {{.AppName}} development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			fiber_run()
		},
	}
)


func fiber_run() {
		configs.AppConfig.SetEnv("dev")
		tp := observe.InitTracer()
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()

		// Basic App Configs
		body_limit, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BODY_LIMIT", "70"))
		read_buffer_size, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("READ_BUFFER_SIZE", "70"))
		rate_limit_per_second, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "5000"))
		//load config file
		app := fiber.New(fiber.Config{
			// Prefork: true,
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
		//  rate limiting middleware
		app.Use(limiter.New(limiter.Config{
			Max:               rate_limit_per_second,
			Expiration:        1 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
		//app logging open telemetery
		app.Use(otelfiber.Middleware())

		// idempotency middleware
		app.Use(idempotency.New(idempotency.Config{
			Lifetime: 10 * time.Second,
				}))

		// early data midllware 
		app.Use(earlydata.New(earlydata.Config{
				Error: fiber.ErrTooEarly,
				// ...
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

		// recover from panic attacks middlerware
		app.Use(recover.New())

		// allow cross origin request
		app.Use(cors.New())
		
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})
		// swagger docs
		app.Get("/docs/*", swagger.HandlerDefault)
		app.Get("/docs/*", swagger.New()).Name("swagger_routes")
		
		// fiber native monitoring metrics endpoint
		app.Get("/lmetrics", monitor.New(monitor.Config{Title: "goBlue Metrics Page"})).Name("custom_metrics_route")
		
		// recover middlware
		
		// adding group with authenthication middleware
		admin_app := app.Group("/api/v1")
		setupRoutes(admin_app.(*fiber.Group))
		
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
	goFrame.AddCommand({{.AppName}}cli)

}

func NextFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

func setupRoutes(gapp *fiber.Group) {
	
	{{range .Models}}
	gapp.Get("/{{.LowerName}}",NextFunc).Name("get_all_{{.LowerName}}s").Get("/{{.LowerName}}", controllers.Get{{.Name}}s)
	gapp.Get("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("get_one_{{.LowerName}}s").Get("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID)
	gapp.Post("/{{.LowerName}}",NextFunc).Name("post_{{.LowerName}}").Post("/{{.LowerName}}", controllers.Post{{.Name}})
	gapp.Patch("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("patch_{{.LowerName}}").Patch("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}})
	gapp.Delete("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("delete_{{.LowerName}}").Delete("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name("delete_{{.LowerName}}")

	{{range .Relations}}
	gapp.Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s)
	{{end}}
	{{end}}


}

`

var prodTemplate = `
package manager

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/observe"
	_ "{{.ProjectName}}.com/docs"
	"{{.ProjectName}}.com/models/controllers"
	"github.com/spf13/cobra"
)

var (
	{{.AppName}}prodcli= &cobra.Command{
		Use:   "prod",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run {{.AppName}} development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			prod_run()
		},
	}
)


func prod_run() {
		configs.AppConfig.SetEnv("prod")
		tp := observe.InitTracer()
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()

		// Basic App Configs
		body_limit, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BODY_LIMIT", "70"))
		read_buffer_size, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("READ_BUFFER_SIZE", "70"))
		rate_limit_per_second, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "5000"))
		//load config file
		app := fiber.New(fiber.Config{
			Prefork: true,
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
		//  rate limiting middleware
		app.Use(limiter.New(limiter.Config{
			Max:               rate_limit_per_second,
			Expiration:        1 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
		//app logging open telemetery
		app.Use(otelfiber.Middleware())

		// idempotency middleware
		app.Use(idempotency.New(idempotency.Config{
			Lifetime: 10 * time.Second,
				}))

		// early data midllware 
		app.Use(earlydata.New(earlydata.Config{
				Error: fiber.ErrTooEarly,
				// ...
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

		// recover from panic attacks middlerware
		app.Use(recover.New())

		// allow cross origin request
		app.Use(cors.New())
		
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})
		// swagger docs
		app.Get("/docs/*", swagger.HandlerDefault)
		app.Get("/docs/*", swagger.New()).Name("swagger_routes")
		
		// fiber native monitoring metrics endpoint
		app.Get("/lmetrics", monitor.New(monitor.Config{Title: "goBlue Metrics Page"})).Name("custom_metrics_route")
		
		// recover middlware
		
		// adding group with authenthication middleware
		admin_app := app.Group("/api/v1")
		setupProdRoutes(admin_app.(*fiber.Group))
		
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
	goFrame.AddCommand({{.AppName}}prodcli)

}

func NextProdFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

func setupProdRoutes(gapp *fiber.Group) {
	
	{{range .Models}}
	gapp.Get("/{{.LowerName}}",NextProdFunc).Name("get_all_{{.LowerName}}s").Get("/{{.LowerName}}", controllers.Get{{.Name}}s)
	gapp.Get("/{{.LowerName}}/:{{.LowerName}}_id",NextProdFunc).Name("get_one_{{.LowerName}}s").Get("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID)
	gapp.Post("/{{.LowerName}}",NextProdFunc).Name("post_{{.LowerName}}").Post("/{{.LowerName}}", controllers.Post{{.Name}})
	gapp.Patch("/{{.LowerName}}/:{{.LowerName}}_id",NextProdFunc).Name("patch_{{.LowerName}}").Patch("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}})
	gapp.Delete("/{{.LowerName}}/:{{.LowerName}}_id",NextProdFunc).Name("delete_{{.LowerName}}").Delete("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name("delete_{{.LowerName}}")

	{{range .Relations}}
	gapp.Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextProdFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextProdFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s)
	{{end}}
	{{end}}


}

`
