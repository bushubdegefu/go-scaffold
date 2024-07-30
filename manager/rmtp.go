package manager

import (
	"fmt"

	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	rmtpcli = &cobra.Command{
		Use:   "rpc",
		Short: "Generate Basic rpc service structure",
		Long:  `Generate Basic rpc service structure`,
		Run: func(cmd *cobra.Command, args []string) {
			rmtpcligen()
		},
	}
)

func rmtpcligen() {
	temps.RtmpNginxFrame()
	temps.NginxConfFrame()
	fmt.Println("After compeleting generating run the below command to build after cd into realtime")
	fmt.Println("docker build -t rmtp-nginx -f rmtp.Dockerfile .")
}

func init() {
	goFrame.AddCommand(rmtpcli)

}
