package rabbitmq

import (
	"fmt"
	"testing"
)

func TestServiceConsumerProcess_Run(t *testing.T) {
	var service ServiceConsumerProcess
	service.QueueName = "service-receiver-two"
	service.Run()
}
func TestDiscoverNumberOfQueueMassage(t *testing.T) {
	res := DiscoverNumberOfQueueMassage("service-receiver-two")
	fmt.Println(">>>>> ", res)
}
