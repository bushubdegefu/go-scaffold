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
)

func fibergen() {
	temps.LoadData()
	temps.FiberFrame()
	temps.TracerFrame()
}

func init() {
	goFrame.AddCommand(fibercli)

}
