package manager

import (
	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	migrationdcli = &cobra.Command{
		Use:   "migration",
		Short: "Generate Data Models based on the GORM using provided spec on the config.json file ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file`,
		Run: func(cmd *cobra.Command, args []string) {
			migrationmod()
		},
	}
)

func migrationmod() {
	temps.MigrationFrame()
}

func init() {
	goFrame.AddCommand(migrationdcli)

}
