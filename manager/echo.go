package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	echocli = &cobra.Command{
		Use:   "echo",
		Short: "generate the basic structure file to start app using echo",
		Long:  `generate the basic structure file to start app using echo`,
		Run: func(cmd *cobra.Command, args []string) {
			echogen()
		},
	}
	gechocli = &cobra.Command{
		Use:   "gecho",
		Short: "generate the basic graphql structure file to start app using echo",
		Long:  `generate the basic graphql structure file to start app using echo`,
		Run: func(cmd *cobra.Command, args []string) {
			graphechogen()
		},
	}
	gechotestcli = &cobra.Command{
		Use:   "techo",
		Short: "generate the basic curd test for echo endpoints",
		Long:  `"generate the basic curd test for echo endpoints"`,
		Run: func(cmd *cobra.Command, args []string) {
			graphechotest()
		},
	}
)

func echogen() {
	temps.LoadData()
	temps.EchoFrame()
	temps.CommonCMD()
}
func graphechogen() {
	temps.LoadData()
	temps.GraphEchoFrame()
	temps.CommonCMD()
}
func graphechotest() {
	temps.LoadData()
	temps.TestFrameEcho()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(echocli)
	goFrame.AddCommand(gechocli)
	goFrame.AddCommand(gechotestcli)

}
