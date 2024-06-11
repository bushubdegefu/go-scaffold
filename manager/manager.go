package manager

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	goFrame = &cobra.Command{
		Use:           "AppMan",
		Short:         "AppMan – command-line tool to aid structure your projects with gorm",
		Long:          "command-line tool to aid structure your projects with gorm. It mainly suitable for fiber and echo app",
		Version:       "0.0.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	if err := goFrame.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
