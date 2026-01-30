package repository

import "context"

type TransactionManager interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}
