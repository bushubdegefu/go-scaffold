package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	fibercli = &cobra.Command{
		Use:   "fiber",
		Short: "generate the basic structure file to start app using fiber",
		Long:  `generate the basic structure file to start app using fiber`,
		Run: func(cmd *cobra.Command, args []string) {
			fibergen()
		},
	}
	fibertestcli = &cobra.Command{
		Use:   "tfiber",
		Short: "generate the basic curd test suites for the fiber handler functions",
		Long:  `generate the basic curd test suites for the fiber handler functions`,
		Run: func(cmd *cobra.Command, args []string) {
			fibertestgen()
		},
	}
)

func fibergen() {
	temps.LoadData()
	temps.FiberFrame()
	temps.FiberFrame()
	temps.CommonCMD()
}

func fibertestgen() {
	temps.LoadData()
	temps.TestFrameFiber()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(fibercli)
	goFrame.AddCommand(fibertestcli)

}
