package post

import (
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks-post/service.go -package=mockspost
type Service interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPostByID(ctx context.Context, id string) (*Post, error)
	GetPostsByDateTime(ctx context.Context, dateTime time.Time) ([]*Post, error)
}

type DefaultService struct {
	repo Repository
}

func (s *DefaultService) CreatePost(ctx context.Context, post *Post) error {
	err := s.repo.New(ctx, post)
	return err
}

func (s *DefaultService) GetPostByID(ctx context.Context, id string) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	return post, err
}

func (s *DefaultService) GetPostsByDateTime(ctx context.Context, dateTime time.Time) ([]*Post, error) {
	posts, err := s.repo.FindByDateTime(ctx, dateTime)
	return posts, err
}

func NewDefaultService(repo Repository) *DefaultService {
	return &DefaultService{
		repo: repo,
	}
}
