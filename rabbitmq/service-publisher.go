package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type ServicePublisher struct {
	QueueName   string
	Message     interface{}
	ContentType string
}

func (s *ServicePublisher) Run() {
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

	// Declare a queue (queue name should match the one used in your consumer)
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

	// Publish a message to the queue
	message, _ := json.Marshal(s.Message) // "Hello, RabbitMQ!"
	err = ch.Publish(
		"",          // exchange
		s.QueueName, // routing key (queue name)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: s.ContentType, // "text/plain",
			Body:        message,
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent: %s\n", string(message))
}
