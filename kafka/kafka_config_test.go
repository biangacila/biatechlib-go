package kafka

import (
	"fmt"
	"testing"
	"time"
)

func TestStartKafka(t *testing.T) {
	fmt.Println("Okay...")
	topic := "SaferSubscriptionService"
	go StartKafka(topic)
	fmt.Println("Kafka has been started.....")
	time.Sleep(10 * time.Minute)
}
