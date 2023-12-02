package rabbitmq

import (
	"fmt"
	"log"
)

type ServiceConsumerProcess struct {
	QueueName string
}

func (s *ServiceConsumerProcess) Run() {

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

	// Define the name of the queue you want to consume from
	queueName := s.QueueName

	// Number of concurrent consumers
	numConsumers := 3

	// Create a channel to signal when consumers are done
	done := make(chan bool)

	// Start multiple consumers
	for i := 0; i < numConsumers; i++ {
		go ConsumeQueueWithProcess(fmt.Sprintf("Agent %v", i+1), queueName, channel, done)
	}

	// Wait for all consumers to finish
	for i := 0; i < numConsumers; i++ {
		<-done
	}

}
