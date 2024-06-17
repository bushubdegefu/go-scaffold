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
	for i := 0; i < len(data.Models); i++ {
		data.Models[i].LowerName = strings.ToLower(data.Models[i].Name)
		data.Models[i].AppName = data.AppName
		data.Models[i].ProjectName = data.ProjectName
		rl_list := make([]Relationship, 0)
		for k := 0; k < len(data.Models[i].RlnModel); k++ {
			rmf := strings.Split(data.Models[i].RlnModel[k], "$")
			cur_relation := Relationship{
				ParentName:      data.Models[i].Name,
				LowerParentName: data.Models[i].LowerName,
				FieldName:       rmf[0],
				LowerFieldName:  strings.ToLower(rmf[0]),
				MtM:             rmf[1] == "mtm",
				OtM:             rmf[1] == "otm",
				MtO:             rmf[1] == "mto",
			}
			rl_list = append(rl_list, cur_relation)
			data.Models[i].Relations = rl_list
		}

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
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"{{.ProjectName}}.com/configs"
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
	
configs.NewEnvFile("./configs")
	app := fiber.New()

	//load config file
	// recover from panic attacks middlerware
	app.Use(recover.New())

	// allow cross origin request
	app.Use(cors.New())
	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")

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
