package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	genmodcli = &cobra.Command{
		Use:   "genmod",
		Short: "Generate Data Models based on the GORM using provided spec on the config.json file ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			genmod()
		},
	}
)

func genmod() {
	temps.ModelDataFrame()
}

func init() {
	goFrame.AddCommand(genmodcli)

}
