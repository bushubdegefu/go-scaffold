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
	gqlcurdclientcli = &cobra.Command{
		Use:   "gqlclient",
		Short: "Generate Basic gql client code store.js file using the config.json file",
		Long:  `Generate Basic gql client code store.js file using the config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			gqlcurdclientgen()
		},
	}
	commonservicecli = &cobra.Command{
		Use:   "service",
		Short: "Generate Basic linux service, docker and git ignore files with basic haproxy cft file config.json file",
		Long:  `Generate Basic linux service, docker and git ignore files with basic haproxy cft file config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			servicecligen()
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

func gqlcurdclientgen() {
	temps.LoadData()
	// gql store.js file generation
	temps.GQLClientFrame()
	temps.CommonGraphQLFrame()
}

func servicecligen() {
	temps.LoadData()
	temps.GitDockerFrame()
	temps.HaproxyFrame()
	temps.ServiceFrame()
	temps.CommonGraphQLFrame()
}

func init() {
	goFrame.AddCommand(gqlcli)
	goFrame.AddCommand(gqlcurdcli)
	goFrame.AddCommand(commonservicecli)
	goFrame.AddCommand(gqlcurdclientcli)

}
