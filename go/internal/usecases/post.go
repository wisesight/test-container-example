package usecases

import (
	"context"
	"time"

	"github.com/wisesight/test-container-example/internal/services/post"
)

func (dfu *defaultUsecases) GetPostByID(id string) (*post.Post, error) {
	ctx, cencel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cencel()
	post, err := dfu.postService.GetPostByID(ctx, id)
	return post, err
}
