package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	gqlcli = &cobra.Command{
		Use:   "gql",
		Short: "Generate Basic gql schema for the provided models on config.json file",
		Long:  `Generate Basic gql schema for the provided models on config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			gqlcligen()
		},
	}
	gqlcurdcli = &cobra.Command{
		Use:   "gqlcurd",
		Short: "Generate Basic gql curd resolver function along with common pagination functions with gorm provided models on config.json file",
		Long:  `Generate Basic gql curd resolver function along with common pagination functions with gorm provided models on config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			gqlcurdcligen()
		},
	}
)

func gqlcligen() {
	temps.LoadData()
	temps.GraphFrame()
	temps.CommonCMD()
}
func gqlcurdcligen() {
	temps.LoadData()
	temps.GraphCurdFrame()
	temps.CommonGraphQLFrame()
}

func init() {
	goFrame.AddCommand(gqlcli)
	goFrame.AddCommand(gqlcurdcli)

}
