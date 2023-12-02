package kafka

type NewMessage struct {
	Msg    interface{}
	Topics string
}

func (obj *NewMessage) Send() {

}
