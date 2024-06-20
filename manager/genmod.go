package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	genechocli = &cobra.Command{
		Use:   "genecho",
		Short: "Generate Data Models based on the GORM using provided spec on the config.json file using Echo Framework ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file Using Echo Framework`,
		Run: func(cmd *cobra.Command, args []string) {
			genecho()
		},
	}
	genfibercli = &cobra.Command{
		Use:   "genfiber",
		Short: "Generate Data Models based on the GORM using provided spec on the config.json file using Fiber Framework ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file Using Fiber Framework`,
		Run: func(cmd *cobra.Command, args []string) {
			genfiber()
		},
	}
)

func genecho() {
	temps.ModelDataFrame()
	temps.DbConnDataFrame()
	temps.CurdFrameEcho()
	temps.TestFrameEcho()
}

func genfiber() {
	temps.ModelDataFrame()
	temps.CurdFrameFiber()
	temps.TestFrameFiber()

}

func init() {
	goFrame.AddCommand(genechocli)
	goFrame.AddCommand(genfibercli)

}
