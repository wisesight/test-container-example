package post

import (
	"time"
)

type Post struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	DateTime time.Time `json:"date_time"`
}
