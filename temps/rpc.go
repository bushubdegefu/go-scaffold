package temps

import (
	"os"
	"text/template"
)

func RPCFrame() {
	// ####################################################
	//  rabbit template
	rpc_tmpl, err := template.New("RenderData").Parse(rpcTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("protoapp", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rpc_file, err := os.Create("protoapp/sample.proto")
	if err != nil {
		panic(err)
	}
	defer rpc_file.Close()

	err = rpc_tmpl.Execute(rpc_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func RPCServiceFrame() {
	// ####################################################
	//  rabbit template
	rpc_tmpl, err := template.New("RenderData").Parse(rpcGolangTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("protoapp", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rpc_file, err := os.Create("protoapp/protoapp.go")
	if err != nil {
		panic(err)
	}
	defer rpc_file.Close()

	err = rpc_tmpl.Execute(rpc_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func RPCManagerFrame() {
	// ####################################################
	//  rabbit template
	rpc_tmpl, err := template.New("RenderData").Parse(rpcServerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("protoapp", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rpc_file, err := os.Create("manager/rpcapp.go")
	if err != nil {
		panic(err)
	}
	defer rpc_file.Close()

	err = rpc_tmpl.Execute(rpc_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var rpcTemplate = `
syntax = "proto3";

option go_package = "./protoapp";

message SampleSalt {
    string salt_a = 1;
    string salt_b = 2;
}

message SampleAppID { 
    string app_id = 1;
}

message SampleAppRoles {
    repeated string roles =1;
}

service SampleService {
    rpc GetSalt(SampleAppID) returns (SampleSalt) {}
    rpc GetSampleAppRoles(SampleAppID) returns (SampleAppRoles) {}
}
`

var rpcGolangTemplate = `
package protoapp

import (
	"fmt"

	"golang.org/x/net/context"
)

type SampleRPCServiceServer struct {
	SampleServiceServer
}

func (server *SampleRPCServiceServer) GetSalt(ctx context.Context, message *SampleAppID) (*SampleSalt, error) {
	fmt.Printf("The APP ID: %v\n", message.AppId)
	salt_a, salt_b := "A build", "B build"
	return &SampleSalt{SaltA: salt_a, SaltB: salt_b}, nil
}

`

var rpcServerTemplate = `
package manager

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"{{.ProjectName}}.com/protoapp"
	"{{.ProjectName}}.com/configs"
)

var (
	startrpc = &cobra.Command{
		Use:   "rpcserve",
		Short: "Start RPC server from the app at the provied Port",
		Long:  "Start RPC server from the app at the provied Port",
		Run: func(cmd *cobra.Command, args []string) {
			RpcServe()
		},
	}

	rpcclient = &cobra.Command{
		Use:   "rpcclient",
		Short: "Make  RPC call server from the app at the provied Port. For testing",
		Long:  "Make RPC call to server from the app at the provied Port through .env variable",
		Run: func(cmd *cobra.Command, args []string) {
			RpcClient()
		},
	}
)

func RpcServe() {
	configs.AppConfig.SetEnv("dev")
	tp := observe.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:"+configs.AppConfig.Get("RPC_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v\n", err, configs.AppConfig.Get("RPC_PORT"))
	}

	protoappserver := protoapp.SampleRPCServiceServer{}
	grpcServer := grpc.NewServer()

	// to registered the defined service server
	protoapp.RegisterSampleServiceServer(grpcServer, &protoappserver)

	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port: %v", err)
	}
	fmt.Println("Started RPC Server for BLUE")

}

func RpcClient() {
	var conn *grpc.ClientConn
	configs.AppConfig.SetEnv("dev")
	tp := observe.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	conn, err := grpc.NewClient(configs.AppConfig.Get("RPC_ADDRESS"),  grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := protoapp.NewSampleServiceClient(conn)

	message := protoapp.SampleAppID{
		AppId: uuid.New().String(),
	}

	response, err := c.GetSalt(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Getting Salt: %s", err)
	}

	log.Printf("Response from Server:\n %s\n", response)

}

func init() {
	goFrame.AddCommand(startrpc)
	goFrame.AddCommand(rpcclient)

}
`
