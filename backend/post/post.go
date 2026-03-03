package post

import "time"

type Post struct {
	ID        int       `json:"id"`
	ImageURL  string    `json:"imageUrl"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}
