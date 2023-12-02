package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type ServiceReceiver struct {
	QueueName string
}

func (s *ServiceReceiver) Run() {
	// Define the RabbitMQ connection URL. Replace with your RabbitMQ server details.
	connURL := "amqp://username:password@rabbitmq-server:5672/"
	connURL = fmt.Sprintf("amqp://%v:%v@%v:%v/", USERNAME, PASSWORD, SERVER, PORT)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue (queue name should match the one used in your producer)
	queueName := s.QueueName
	_, err = ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-acknowledge
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	fmt.Println(" [*] Waiting for messages. To exit, press CTRL+C")

	for msg := range msgs {
		fmt.Printf(" [x] Received: %s\n", msg.Body)
	}
}
