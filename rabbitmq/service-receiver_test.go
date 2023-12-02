package rabbitmq

import "testing"

func TestServiceReceiver_Run(t *testing.T) {
	var service ServiceReceiver
	service.QueueName = "swaly-exetat-user"

	service.Run()
}
