package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func StartKafka(topic string) {
	conf := kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%v:%v", KAFKA_HOST, KAFKA_PORT)},
		GroupID:  "1",
		Topic:    topic,
		MaxBytes: 1000,
	}
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Some error occured", err)
			continue
		}
		fmt.Printf("Message is: %v\n", string(m.Value))
	}
}
