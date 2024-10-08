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
		Short: "Generate Data Models & CURD handler based on the GORM using provided spec on the config.json file using Fiber Framework ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file Using Fiber Framework`,
		Run: func(cmd *cobra.Command, args []string) {
			genfiber()
		},
	}
	fibcurdcli = &cobra.Command{
		Use:   "curdfiber",
		Short: "Generate CURD handlers the GORM using provided spec on the config.json file using Fiber Framework ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file Using Fiber Framework`,
		Run: func(cmd *cobra.Command, args []string) {
			genfibercurd()
		},
	}
	gormmodelscli = &cobra.Command{
		Use:   "models",
		Short: "Generate Models the GORM using provided spec on the config.json file using Fiber Framework ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file Using Fiber Framework`,
		Run: func(cmd *cobra.Command, args []string) {
			gengorm()
		},
	}
)

func genecho() {
	temps.LoadData()
	temps.ModelDataFrame()
	temps.DbConnDataFrame()
	temps.CurdFrameEcho()
	temps.TestFrameEcho()
	temps.CommonCMD()
}

func genfiber() {
	temps.LoadData()
	temps.ModelDataFrame()
	temps.CurdFrameFiber()
	temps.TestFrameFiber()
	temps.CommonCMD()
}
func genfibercurd() {
	temps.LoadData()
	temps.CurdFrameFiber()
	temps.CommonCMD()
}
func gengorm() {
	temps.LoadData()
	temps.ModelDataFrame()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(genechocli)
	goFrame.AddCommand(genfibercli)
	goFrame.AddCommand(fibcurdcli)
	goFrame.AddCommand(gormmodelscli)

}
