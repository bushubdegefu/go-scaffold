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
)

func gqlcligen() {
	temps.LoadData()
	temps.GraphFrame()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(gqlcli)

}
