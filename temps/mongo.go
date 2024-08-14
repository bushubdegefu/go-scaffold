package temps

import (
	"os"
	"text/template"
)

func MongoDataBaseFrame() {
	// ####################################################
	//  rabbit template
	rpc_tmpl, err := template.New("RenderData").Parse(mongoDatabaseConn)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("nosqlconn", os.ModePerm)
	if err != nil {
		panic(err)
	}

	nosqlconn_file, err := os.Create("nosqlconn/nosqlconn.go")
	if err != nil {
		panic(err)
	}
	defer nosqlconn_file.Close()

	err = rpc_tmpl.Execute(nosqlconn_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var mongoDatabaseConn = `
package nosqlconn

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"{{.ProjectName}}.com/configs"
)

// MongoSession creates and returns a MongoDB client instance.
func MongoSession() (*mongo.Database,*mongo.Client,error) {
	// Retrieve MongoDB URI and database name from environment variables
	dsn := configs.AppConfig.Get("MONGO_URI")
	databaseName := configs.AppConfig.Get("MONGO_DB_NAME")

	if dsn == "" || databaseName == "" {
		return nil,nil,fmt.Errorf("missing MongoDB URI or database name. Please set 'MONGO_URI' and 'MONGO_DB_NAME' environment variables")
	}

	// Setup options with OpenTelemetry middleware
	opts := options.Client().
		ApplyURI(dsn).
		SetMonitor(otelmongo.NewMonitor())

	// Create a new MongoDB client
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil,nil,fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

// DisconnectClient disconnects the MongoDB client.
func DisconnectClient(client *mongo.Client) error {
	if client == nil {
		return nil
	}
	if err := client.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	return nil
}
`
