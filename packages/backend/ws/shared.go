package ws

type SharedAbb struct {
	Create   bool   `json:"create,omitempty"`
	Delete   bool   `json:"delete,omitempty"`
	Override bool   `json:"override,omitempty"`
	Abb      string `json:"abb,omitempty"`
	Word     string `json:"word,omitempty"`
}

type SharedManuscript struct {
}
