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
)

func echogen() {
	temps.LoadData()
	temps.EchoFrame()
}

func init() {
	goFrame.AddCommand(echocli)

}
