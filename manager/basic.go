package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	basicstruct = &cobra.Command{
		Use:   "basic",
		Short: "generate basic folder structure for project",
		Long:  `Generate basic folder structure for your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			basiccmd()
		},
	}
)

func basiccmd() {
	temps.Frame()
	temps.CommonFrame()
}

func init() {
	goFrame.AddCommand(basicstruct)

}
