package manager

import (
	"fmt"

	"github.com/spf13/cobra"
	"scaffold.com/temps"
)

var (
	rpccli = &cobra.Command{
		Use:   "rpc",
		Short: "Generate Basic rpc service structure",
		Long:  `Generate Basic rpc service structure`,
		Run: func(cmd *cobra.Command, args []string) {
			rpcgen()
		},
	}
)

func rpcgen() {
	temps.LoadData()
	temps.RPCFrame()
	temps.RPCManagerFrame()
	temps.RPCServiceFrame()
	fmt.Println("After compeleting generating run the below command after installing protoc on you machine")
	fmt.Println("protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative protoapp/sample.proto")
}

func init() {
	goFrame.AddCommand(rpccli)

}
