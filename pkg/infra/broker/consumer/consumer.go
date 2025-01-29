package consumer

type Consumer interface {
	Consume(msg any) error
}
