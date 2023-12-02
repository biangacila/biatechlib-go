package rabbitmq

import (
	"log"
)

type ServiceReceiverTwo struct {
	QueueName string
}

func (s *ServiceReceiverTwo) Run() {
	// Connect to RabbitMQ
	connection, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer connection.Close()

	// Create a channel
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	// Define the names of the two queues you want to listen to
	server1 := "server_1"
	server2 := "server_2"

	// Start Goroutines to consume messages from the queues
	go consumeQueue(s.QueueName, server1, channel)
	go consumeQueue(s.QueueName, server2, channel)

	// Keep the application running
	select {}
}
