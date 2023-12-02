package kafka

import "testing"

func TestServiceConnection1_Consumer(t *testing.T) {
	var hub ServiceConnection1
	hub.Consumer("EasiPath")
}
func TestServiceConnection1_Prducer(t *testing.T) {
	var hub ServiceConnection1
	hub.Producer("Biacibenga", `{"action":"welcome notice"}`)
}

//   N3M3NDk6
//  Ticket ID : 86277
