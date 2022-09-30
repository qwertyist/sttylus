package document

import "time"

type Document struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Creator string    `json:"creator"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
