package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

const KAFKA_HOST = "voip.easipath.com" // "129.232.145.82"
const KAFKA_PORT = "9092"

type ServiceConnection1 struct {
}

func (s *ServiceConnection1) Producer(topic string, info interface{}) {

	var urlConnection = fmt.Sprintf("%v:%v", KAFKA_HOST, KAFKA_PORT)
	// Kafka producer configuration
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for only local acknowledgment
	producerConfig.Producer.Return.Successes = true          // Receive acknowledgment for successful messages

	// Create a Kafka producer
	producer, err := sarama.NewAsyncProducer([]string{urlConnection}, producerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	// Send a message to a Kafka topic
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(interfaceToString(info)),
	}

	// Send the message asynchronously
	producer.Input() <- message

	// Wait for acknowledgment
	select {
	case success := <-producer.Successes():
		log.Printf("Message sent successfully: %v", success.Offset)
	case err := <-producer.Errors():
		log.Printf("Error sending message: %v", err.Err)
	}

}

func (s *ServiceConnection1) Consumer(topic string) {
	var urlConnection = fmt.Sprintf("%v:%v", KAFKA_HOST, KAFKA_PORT)
	// Kafka consumer configuration
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	// Create a Kafka consumer
	consumer, err := sarama.NewConsumer([]string{urlConnection}, consumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to a Kafka topic
	//topic := "your-topic"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Consume messages from the topic
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err.Err)
		}
	}
}

func interfaceToString(info interface{}) string {
	i, _ := json.Marshal(info)
	return string(i)
}
