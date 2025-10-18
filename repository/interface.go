package repository

import "context"

type RepoInterface interface {
	Increment(ctx context.Context, key string, expiry int) (int, error)
	GetByKey(ctx context.Context, key string) (int, error)
}
