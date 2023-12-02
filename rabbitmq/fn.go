package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func connectRabbitMQ() (*amqp.Connection, error) {
	// RabbitMQ connection URL
	rabbitMQURL := "amqp://username:password@localhost:5672/"
	rabbitMQURL = fmt.Sprintf("amqp://%v:%v@%v:%v/", USERNAME, PASSWORD, SERVER, PORT)

	// Establish a connection to RabbitMQ server
	connection, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func consumeQueue(queueName string, serverName string, ch *amqp.Channel) {
	// Declare the queue
	_, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // Queue name
		"",        // Consumer name (empty string for auto-generated)
		true,      // Auto-acknowledge messages
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process incoming messages
	for msg := range msgs {
		fmt.Printf("%v | Received a message: %s\n", serverName, msg.Body)
	}
}

func ConsumeQueueWithProcess(agent string, queueName string, ch *amqp.Channel, done chan bool) {
	// Declare the queue
	_, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // Queue name
		agent,     // Consumer name (empty string for auto-generated)
		false,     // Manual acknowledgement mode (auto-acknowledge: false)
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process incoming messages
	for msg := range msgs {
		// Simulate message processing (replace with your own logic)
		processMessage(msg.Body, agent)

		// Acknowledge the message
		msg.Ack(false)
	}

	done <- true // Signal that this consumer is done
}
func processMessage(body []byte, agent string) {
	fmt.Printf("Processing message: %s | %v\n", body, agent)
	// Add your message processing logic here
	// For example, you can insert the logic to distribute messages to call center agents.
	<-time.After(30 * time.Second)

}

func DiscoverNumberOfQueueMassage(queueName string) int {
	// Replace with your RabbitMQ Management API URL.
	apiUrl := "http://" + SERVER + ":1" + PORT + "/api/queues/%2f/" + queueName // Replace "my_queue" with your queue name.

	// Replace with your RabbitMQ Management API credentials.
	username := USERNAME
	password := PASSWORD

	// Create an HTTP client with basic authentication.
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.SetBasicAuth(username, password)

	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code.
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse JSON response.
	var queueInfo map[string]interface{}
	if err := json.Unmarshal(body, &queueInfo); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Extract the number of messages in the queue.
	messagesCount := int(queueInfo["messages"].(float64))

	fmt.Printf("Number of messages in the queue: %d\n", messagesCount)

	return messagesCount
}
