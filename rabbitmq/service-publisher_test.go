package rabbitmq

import "testing"

func TestServiceConsumer_Run(t *testing.T) {
	var service ServicePublisher
	service.QueueName = "swaly-exetat-user"
	service.ContentType = "text/plain"
	service.Message = `{"Name":"Merveilleux Biangacila","Email":"biangacila@gmail.com","Phone":"27729139504"}`
	service.Run()
}
