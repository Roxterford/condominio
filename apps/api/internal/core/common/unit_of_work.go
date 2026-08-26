package common

import "context"

type UnitOfWork[T any] interface {
	Do(ctx context.Context, fn func(deps T) error) error
}
