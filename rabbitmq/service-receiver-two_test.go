package rabbitmq

import "testing"

func TestServiceReceiverTwo_Run(t *testing.T) {
	var service ServiceReceiverTwo
	service.QueueName = "service-receiver-two"
	service.Run()
}
