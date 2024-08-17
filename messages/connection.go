
package messages

import (
	"crypto/tls"
	"fmt"

	"github.com/streadway/amqp"
	"mongo-play.com/configs"
)


func QeueConnect(queue_name string) (*amqp.Connection, *amqp.Channel, error) {

	// Getting Rabbit URI from the ENV file
	con_str := configs.AppConfig.Get("RABBIT_URI")

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
		queue_name, // queue name
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
