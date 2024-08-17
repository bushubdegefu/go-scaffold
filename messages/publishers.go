
package messages

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type SampleMessage struct {}

func PublishMessageQueue(posted_message RequestObject, queue_name string) error {

	//   connection and channels from rabbitmq
	connection, channel, _ := QeueConnect(queue_name)
	defer connection.Close()
	defer channel.Close()

	// Create a message to publish.
	queue_message, _ := json.Marshal(posted_message)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(queue_message),
		Type:        "REQUEST",
	}

	//send to rabbit app module qeue using channel
	// Attempt to publish a message to the queue.
	if err := channel.Publish(
		"",         // exchange
		queue_name, // queue name
		false,      // mandatory
		false,      // immediate
		message,    // message to publish
	); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

