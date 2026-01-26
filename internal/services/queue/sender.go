package queue

import "context"

type Sender interface {
	SendMessage(ctx context.Context, payload string) error
}
