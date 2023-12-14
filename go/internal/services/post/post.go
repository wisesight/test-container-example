package post

import (
	"time"
)

type Post struct {
	ID       string    `json:"id" bson:"_id,omitempty"`
	Title    string    `json:"title" bson:"title,omitempty"`
	Body     string    `json:"body" bson:"body,omitempty"`
	DateTime time.Time `json:"date_time" bson:"date_time,omitempty"`
}
