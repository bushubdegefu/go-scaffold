package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func MigrationFrame() {
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

	//  this is creating manger file inside the manager folder
	// ############################################################
	migration_tmpl, err := template.New("data").Parse(migrationTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	migration_file, err := os.Create("manager/migrate.go")
	if err != nil {
		panic(err)
	}
	defer migration_file.Close()

	err = migration_tmpl.Execute(migration_file, data)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var migrationTemplate = `
package manager

import (
	"fmt"
	"{{.ProjectName}}.com/models"
	"github.com/spf13/cobra"
)

var (
	{{.AppName}}migrate= &cobra.Command{
		Use:   "migrate",
		Short: "Run Database Migration for found in init migration Models",
		Long:  {{.BackTick}}Migrate to init database{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			init_migrate()
		},
	}
)

func init_migrate() {
	models.InitDatabase()
	fmt.Println("Migrated Database Models sucessfully")
}

func init() {
	goFrame.AddCommand({{.AppName}}migrate)

}

`
