package post

import "time"

type Response struct {
	ID        int       `json:"id"`
	ImageURL  string    `json:"imageUrl"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p Post) ToResponse(baseURL string) Response {
	return Response{
		ID:        p.ID,
		ImageURL:  baseURL + "/images/" + p.ImageURL,
		Text:      p.Text,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
	}
}
