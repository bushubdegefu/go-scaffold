package temps

import (
	"os"
	"text/template"
)

func RabbitFrame() {
	// ####################################################
	//  rabbit template
	rab_tmpl, err := template.New("RenderData").Parse(rabbitConnectionTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("messages", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rab_file, err := os.Create("messages/connection.go")
	if err != nil {
		panic(err)
	}
	defer rab_file.Close()

	err = rab_tmpl.Execute(rab_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func PublishFrame() {
	// ####################################################
	//  rabbit template
	rab_tmpl, err := template.New("RenderData").Parse(pusbsttuctTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("messages", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rab_file, err := os.Create("messages/publishers.go")
	if err != nil {
		panic(err)
	}
	defer rab_file.Close()

	err = rab_tmpl.Execute(rab_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func ConsumeFrame() {
	// ####################################################
	//  rabbit template
	rab_tmpl, err := template.New("RenderData").Parse(constumerBasicTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("messages", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rab_file, err := os.Create("messages/consumer.go")
	if err != nil {
		panic(err)
	}
	defer rab_file.Close()

	err = rab_tmpl.Execute(rab_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func RunConsumeFrame() {
	// ####################################################
	//  rabbit template
	rab_tmpl, err := template.New("RenderData").Parse(rabbitRunTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("messages", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rab_file, err := os.Create("manager/consumer.go")
	if err != nil {
		panic(err)
	}
	defer rab_file.Close()

	err = rab_tmpl.Execute(rab_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var rabbitConnectionTemplate = `
package rabbit

import (
	"crypto/tls"
	"fmt"

	"github.com/streadway/amqp"
	"{{.ProjectName}}.com/configs"
)

// creating connection to the rabbit message broker
// returns the connection based on the connection string
// needs to be closed after using by functions using it
// returns connection and channel struct
func BrokerConnect() (*amqp.Connection, *amqp.Channel, error) {

	con_str := configs.AppConfig.Get("RABBIT_URI")

	// connection, err := amqp.Dial(config.Config("RABBIT_BROKER_URL"))
	// connection, err := amqp.Dial(con_str)
	// if err != nil {
	// 	fmt.Printf("connectin to %v failed due to : %v\n", con_str, err)
	// }

	// RabbitMQ TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to false for production use
	}

	// Dial RabbitMQ server with TLS
	connection, err := amqp.DialTLS(con_str, tlsConfig)
	if err != nil {
		fmt.Printf("connectin to %v failed due to : %v \n", con_str, err)
	}

	// creating a channel to create a queue
	// instance over the connection we have already
	// established.
	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("connectin to channel failed due to : %v\n", err)
	}

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channel.QueueDeclare(
		"{{.ProjectName}}", // queue name
		true,        // durable
		false,       // auto delete
		false,       // exclusive
		false,       // no wait
		nil,         // arguments
	)

	if err != nil {
		connection.Close() // Close the connection if queue declaration fails
		channel.Close()    // Close the channel
		fmt.Printf("creating queue to %v failed due to : %v\n",con_str, err)
	}
	return connection, channel, nil

}
`

var pusbsttuctTemplate = `
package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type SampleMessage struct {}

func PublishMessage(posted_message SampleMessage) error {

	//   connection and channels from rabbitmq
	connection, channel, _ := BrokerConnect()
	defer connection.Close()
	defer channel.Close()

	// Create a message to publish.
	queue_message, _ := json.Marshal(posted_message)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(queue_message),
		Type:        "BULK_MAIL",
	}

	//send to rabbit app module qeue using channel
	// Attempt to publish a message to the queue.
	if err := channel.Publish(
		"",          // exchange
		"{{.ProjectName}}", // queue name
		false,       // mandatory
		false,       // immediate
		message,     // message to publish
	); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// //  Injecting OtelDetails to Header for tracing
// header := make(map[string][]string)
// propagator := propagation.TraceContext{}
// propagator.Inject(tracer.Tracer, propagation.HeaderCarrier(header))
//  use the header above to assign value to tp when initalizing the request object below

type RequestObject struct {
	Host     string
	Endpoint string
	Method   string
	Body     string
	Tp       map[string][]string
}

func PublishRequest(req RequestObject) error {

	//   connection and channels from rabbitmq
	connection, channel, _ := BrokerConnect()
	defer connection.Close()
	defer channel.Close()

	//
	request, _ := json.Marshal(req)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(request),
		Type:        "REQUEST",
	}

	//send to rabbit app module qeue using channel
	// Attempt to publish a message to the queue.
	if err := channel.Publish(
		"",      // exchange
		"esb",   // queue name
		false,   // mandatory
		false,   // immediate
		message, // message to publish
	); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

`

var constumerBasicTemplate = `
package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"{{.ProjectName}}.com/configs"
	"{{.ProjectName}}.com/observe"

)

type sample_message struct {}

func RabbitConsumer() {
	configs.AppConfig.SetEnv("dev")
	//  tracer
	tp := observe.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
		
	// Getting app connection and channel
	connection, channel, err := BrokerConnect()
	if err != nil {
		fmt.Println("Failed to establish connection:", err)
		return
	}
	defer connection.Close()
	defer channel.Close()

	// ########################################
	// Declaring consumer with its properties over the channel opened
	msgs, err := channel.Consume(
		"{{.ProjectName}}", // queue
		"",          // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)

	// ###########################################

	if err != nil {
		fmt.Println("Failed to consume messages:", err)
		return
	}

	// Process received messages based on their types
	// Using a goroutine for asynchronous message consumption
	go func(msg <-chan amqp.Delivery) {
		ctx := context.Background()

		for msg := range msgs {
			// Extract the span context out of the AMQP header.
			switch msg.Type {
			case "BULK_MAIL":    // make sure provide the type in the published message so to switch
				var message sample_message
				err := json.Unmarshal(msg.Body, &message)
				if err != nil {
					fmt.Println("Failed to unmarshal message:", err)
					continue
				}
				fmt.Println(message)
			case "REQUEST":
				var reqData RequestObject
				err := json.Unmarshal(msg.Body, &reqData)
				if err != nil {
					fmt.Println(err.Error())
				}
					
				// extracting request from otel propagator
				propagator := propagation.TraceContext{}
				ctx = propagator.Extract(ctx, propagation.HeaderCarrier(reqData.Tp))

				// starting span for http client request
				trace, span := observe.AppTracer.Start(ctx, fmt.Sprintf("started %v", rand.Intn(100)))

				// client with otel http middleware
				client := http.Client{
					Transport: otelhttp.NewTransport(
						http.DefaultTransport,
						// By setting the otelhttptrace client in this transport, it can be
						// injected into the context after the span is started, which makes the
						// httptrace spans children of the transport one.
						otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
							return otelhttptrace.NewClientTrace(ctx)
						}),
					),
				}

				// Build request to be sent
				req, _ := http.NewRequestWithContext(trace, reqData.Method, fmt.Sprintf("http://%v%v", reqData.Host, reqData.Endpoint), nil)


				// generating uuid
				gen, _ := uuid.NewV7()
				id := gen.String()
				//  performing the request
				req_body := fmt.Sprintf("Method: %v\t %v \nBody: %v\n", req.Method, req.URL, req.Body)
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("failed to perform request: %v\n", err)
					span.SetAttributes(attribute.String("esb-id", id))
					span.SetAttributes(attribute.String("esb-request", req_body))
					span.SetAttributes(attribute.String("esb-error", err.Error()))
					span.End()
					break
				}
				//

				body, _ := io.ReadAll(resp.Body)
				span.SetAttributes(attribute.String("esb-id", id))
				span.SetAttributes(attribute.String("esb-request", req_body))
				span.SetAttributes(attribute.String("esb-response", string(body)))

				span.End()
				
			default:
				fmt.Println("Unknown Task Type")
			}
		}
	}(msgs)

	fmt.Println("Waiting for messages...")
	select {}
}
`

var rabbitRunTemplate = `
package manager

import (
	"github.com/spf13/cobra"
	"{{.ProjectName}}.com/messages"
)

var (
	startconsumercli = &cobra.Command{
		Use:   "start",
		Short: "start rabbit consumer",
		Long:  "Start rabbit app consumer",
		Run: func(cmd *cobra.Command, args []string) {
			startconsumer()
		},
	}
)

func startconsumer() {
	rabbit.RabbitConsumer()
}

func init() {
	goFrame.AddCommand(startconsumercli)

}
`
