package storage

import (
	"context"
	"time"
)

type Storage interface {
	AddAccess(ctx context.Context, keyType string, key string, maxAccesses int64) (bool, int64, error)
	GetBlock(ctx context.Context, keyType string, key string) (*time.Time, error)
	AddBlock(ctx context.Context, keyType string, key string, milliseconds int64) (*time.Time, error)
}
