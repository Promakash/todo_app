package producer

import "context"

type Producer interface {
	Produce(ctx context.Context, msg any) error
}
