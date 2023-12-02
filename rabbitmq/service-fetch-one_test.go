package rabbitmq

import "testing"

func TestServiceFetchOne_Run(t *testing.T) {
	var service ServiceFetchOne
	service.QueueName = "service-receiver-two"
	service.Run()
}
