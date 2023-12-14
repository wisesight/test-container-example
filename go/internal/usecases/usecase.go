package usecases

import "github.com/wisesight/test-container-example/internal/services/post"

//go:generate mockgen -source=usecase.go -destination=mocksusecases/usecase.go -package=mocksusecases
type Usecases interface {
	GetPostByID(id string) (*post.Post, error)
}

type defaultUsecases struct {
	postService post.Service
}

func NewDefaultUsecases(postService post.Service) Usecases {
	return &defaultUsecases{
		postService: postService,
	}
}
