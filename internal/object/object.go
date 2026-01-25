package object

import "context"

type SaverInVault interface {
	PutObject(ctx context.Context, path string) error
}
