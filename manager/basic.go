package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	basicstruct = &cobra.Command{
		Use:   "basic",
		Short: "generate basic folder structure for project more fiber framework focused",
		Long:  `Generate basic folder structure for your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			basiccmd()
		},
	}

	standardstruct = &cobra.Command{
		Use:   "standard",
		Short: "generate basic folder structure for project standard, generic structure",
		Long:  `Generate basic folder structure for your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			standardcmd()
		},
	}

	standardrabbit = &cobra.Command{
		Use:   "rabbit",
		Short: "generate basic folder structure for project standard rabbit connection,publisher and consumer, generic structure",
		Long:  `Generate basic folder structure for your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			standrabbit()
		},
	}

	standarpubcli = &cobra.Command{
		Use:   "publish",
		Short: "generate basic folder structure for project standard rabbit connection and publisher, generic structure",
		Long:  `Generate basic folder structure for your project.`,
		Run: func(cmd *cobra.Command, args []string) {
			standpublish()
		},
	}
)

func basiccmd() {
	temps.LoadData()
	temps.Frame()
	temps.DbConnDataFrame()
	temps.CommonFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.FiberTracerFrame()
	temps.GitDockerFrame()
}

func standardcmd() {
	temps.LoadData()
	temps.Frame()
	temps.DbConnDataFrame()
	temps.StandardTracerFrame()
	temps.CommonFrame()
	temps.RabbitFrame()
}

func standrabbit() {
	temps.LoadData()
	temps.Frame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.ConsumeFrame()
	temps.RunConsumeFrame()

}
func standpublish() {
	temps.LoadData()
	temps.Frame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
}

func init() {
	goFrame.AddCommand(basicstruct)
	goFrame.AddCommand(standardstruct)
	goFrame.AddCommand(standardrabbit)
	goFrame.AddCommand(standarpubcli)
}
