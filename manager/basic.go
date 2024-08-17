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

	standardbcli = &cobra.Command{
		Use:   "db",
		Short: "generate basic folder struct for the database connection file( mainly sqlite and postgres)",
		Long:  `generate basic folder struct for the database connection file( mainly sqlite and postgres)`,
		Run: func(cmd *cobra.Command, args []string) {
			standdatabase()
		},
	}

	standarnodbcli = &cobra.Command{
		Use:   "nodb",
		Short: "generate basic folder struct without database connection file( with out databse conn struct)",
		Long:  `generate basic folder struct without database connection file( with out databse conn struct)`,
		Run: func(cmd *cobra.Command, args []string) {
			standnodbversion()
		},
	}
	standarnosqldbcli = &cobra.Command{
		Use:   "nosql",
		Short: "generate basic folder struct for nosql models for app logic( with out databse n)",
		Long:  `generate basic folder struct for nosql models for app logic`,
		Run: func(cmd *cobra.Command, args []string) {
			standnosqlmongo()
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
	temps.HaproxyFrame()
	temps.ServiceFrame()
	temps.CommonCMD()
}

func standardcmd() {
	temps.LoadData()
	temps.Frame()
	temps.DbConnDataFrame()
	temps.StandardTracerFrame()
	temps.CommonFrame()
	temps.RabbitFrame()
	temps.CommonCMD()
}

func standrabbit() {
	temps.LoadData()
	temps.Frame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.ConsumeFrame()
	temps.RunConsumeFrame()
	temps.CommonCMD()

}
func standpublish() {
	temps.LoadData()
	temps.Frame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.CommonCMD()
}

func standdatabase() {
	temps.LoadData()
	temps.Frame()
	temps.DbConnDataFrame()
	temps.CommonCMD()
}

func standnodbversion() {
	temps.LoadData()
	temps.Frame()
	temps.CommonCMD()
}

func standnosqlmongo() {
	temps.LoadData()
	temps.Frame()
	temps.MongoDataBaseFrame()
	temps.NoSQLModelDataFrame()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(basicstruct)
	goFrame.AddCommand(standardstruct)
	goFrame.AddCommand(standardrabbit)
	goFrame.AddCommand(standarpubcli)
	goFrame.AddCommand(standardbcli)
	goFrame.AddCommand(standarnodbcli)
	goFrame.AddCommand(standarnosqldbcli)
}
