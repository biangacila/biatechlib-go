package rabbitmq

import "log"

type ServiceFetchOne struct {
	QueueName string
}

func (s *ServiceFetchOne) Run() {
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

	// Use basic.get to retrieve a single message from the queue.
	msg, _, err := channel.Get(queueName, true) // true indicates auto-acknowledgment.
	if err != nil {
		log.Fatalf("Failed to get a message: %v", err)
	}

	// Process the received message here.
	log.Printf("Received message: %s", msg.Body)

	// Acknowledge the message manually (if needed).
	// msg.Ack(false)

}
