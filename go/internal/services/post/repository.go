package post

import (
	"context"
	"time"
)

//go:generate mockgen -source=repository.go -destination=mockspost/repository.go -package=mockspost
type Repository interface {
	New(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	FindByDateTime(ctx context.Context, dateTime time.Time) ([]*Post, error)
}
