package temps

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func Frame() {

	//  this is creating manger file inside the manager folder
	// ############################################################
	manager_tmpl, err := template.New("RenderData").Parse(managerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	manager_file, err := os.Create("manager/manager.go")
	if err != nil {
		panic(err)
	}
	defer manager_file.Close()

	err = manager_tmpl.Execute(manager_file, RenderData)
	if err != nil {
		panic(err)
	}

	//this is creating the main.go file
	// ############################################################
	main_tmpl, err := template.New("RenderData").Parse(mainTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	main_file, err := os.Create("main.go")
	if err != nil {
		panic(err)
	}
	defer main_file.Close()

	err = main_tmpl.Execute(main_file, RenderData)
	if err != nil {
		panic(err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	config_tmpl, err := template.New("RenderData").Parse(configsTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("configs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		panic(err)
	}

	config_file, err := os.Create("configs/configs.go")
	if err != nil {
		panic(err)
	}
	defer config_file.Close()

	err = config_tmpl.Execute(config_file, RenderData)
	if err != nil {
		panic(err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	env_tmpl, err := template.New("RenderData").Parse(envTemplate)
	if err != nil {
		panic(err)
	}

	env_file, err := os.Create("configs/.env")
	if err != nil {
		panic(err)
	}
	defer env_file.Close()

	tenv_file, err := os.Create("tests/.env")
	if err != nil {
		panic(err)
	}
	defer tenv_file.Close()

	err = env_tmpl.Execute(tenv_file, RenderData)
	if err != nil {
		panic(err)
	}

	err = env_tmpl.Execute(env_file, RenderData)
	if err != nil {
		panic(err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	devenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		panic(err)
	}

	devenv_file, err := os.Create("configs/.dev.env")
	if err != nil {
		panic(err)
	}
	defer devenv_file.Close()

	err = devenv_tmpl.Execute(devenv_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ############################################################
	prodenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		panic(err)
	}

	prodenv_file, err := os.Create("configs/.prod.env")
	if err != nil {
		panic(err)
	}
	defer prodenv_file.Close()

	err = prodenv_tmpl.Execute(prodenv_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ##########################################################
	testenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		panic(err)
	}

	testenv_file, err := os.Create("tests/.test.env")
	if err != nil {
		panic(err)
	}
	defer testenv_file.Close()

	err = testenv_tmpl.Execute(testenv_file, RenderData)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var mainTemplate = `
package main

import (
	"{{.ProjectName}}.com/manager"
)

//	@title			Swagger {{.ProjectName}} API
//	@version		0.1
//	@description	This is {{.ProjectName}} API OPENAPI Documentation.
//	@termsOfService	http://swagger.io/terms/
//  @BasePath  /api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-APP-TOKEN
//	@description				Description for what is this security definition being used

//	@securityDefinitions.apikey Refresh
//	@in							header
//	@name						X-REFRESH-TOKEN
//	@description				Description for what is this security definition being used

func main() {
	manager.Execute()
}
`

var managerTemplate = `
package manager
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	goFrame = &cobra.Command{
		Use:           "AppMan",
		Short:         "AppMan – command-line tool to aid structure you fiber backend projects with gorm",
		Long:          "Manager File Framed by go frame",
		Version:       "0.0.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	if err := goFrame.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
`

var configsTemplate = `
package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultFileName         = "/.env"
	defaultOverrideFileName = "/.local.env"
)

type EnvConfig struct {
	defaultPath string
	prodFlag    string
}

type Config interface {
	Get(string) string
	GetOrDefault(string, string) string
}

var AppConfig EnvConfig

func NewEnvFile(configFolder string) {
	AppConfig = EnvConfig{
		defaultPath: configFolder,
	}
	AppConfig.read()
}

func (e *EnvConfig) read() {
	var (
		defaultFile  = e.defaultPath + defaultFileName
		overrideFile = e.defaultPath + defaultOverrideFileName
	)
	err := godotenv.Overload(defaultFile)
	env := e.Get("APP_ENV")
	if err != nil {
		fmt.Printf("WARNING: Failed to load config from file: %v, Err: %v \n", defaultFile, err)
	} else {
		fmt.Printf("INFO: Loaded config from file: %v\n", defaultFile)
	}

	// If 'APP_ENV' is set to x, then App will read '.env' from configs directory, and then it will be overwritten
	// by configs present in file '.x.env'
	overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, env)
	if env == "" && e.prodFlag == "" {
		overrideFile = fmt.Sprintf("%s/dev.env", e.defaultPath)
	}

	if e.prodFlag != "" {
		overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, e.prodFlag)
	}
	err = godotenv.Overload(overrideFile)
	if err != nil {
		fmt.Printf("WARNING: to load config from file: %v, Err: %v \n", overrideFile, err)
	} else {
		fmt.Printf("INFO: config from file: %v \n", overrideFile)
	}
}

func (e *EnvConfig) Get(key string) string {
	return os.Getenv(key)
}

func (e *EnvConfig) SetEnv(key string) {
	AppConfig = EnvConfig{
		defaultPath: "./configs",
		prodFlag: key,
	}
	AppConfig.read()
}

func (e *EnvConfig) GetOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
`

var envTemplate = `
APP_ENV=dev
`
var devenvTemplate = `
APP_NAME=dev
HTTP_PORT=7500
TEST_NAME="Development Development"
BODY_LIMIT=70
READ_BUFFER_SIZE=40
RATE_LIMIT_PER_SECOND=5000


#Database config settings
DB_TYPE=postgres
POSTGRES_URI="host=localhost user=blueuser password=default dbname=learning_one port=5432 sslmode=disable"
#DB_TYPE="sqlite"
#SQLLITE_URI="goframe.db"

#Messeage qeue settings specifically rabbit
RABBIT_URI="amqps://fldooyrt:bQLS-hvGReUQPy_kUSexUYOhzJHwf08u@lionfish.rmq.cloudamqp.com:5671/fldooyrt"

# JWT token settings 
JWT_SALT_LIFE_TIME=60 #in minutes
JWT_SALT_LENGTH=25

#RPC settings
RPC_PORT=6500


#Observability settings
TRACE_EXPORTER=jaeger
TRACER_HOST=localhost
TRACER_PORT=14317
`
