package abbreviation

type query struct {
	Standard string   `json:"standard"`
	Addon    []string `json:"addon"`
	ListIDs  []string `json:"list_ids"`

	UserID string `json:"user_id"`
	User   string `json:"user"`
	Abb    string `json:"abb"`
}

type response struct {
	Word   string `json:"word"`
	Abb    string `json:"abb"`
	Missed bool   `json:"missed"`
}

type conflict struct {
	Abb string `json:"abb,omitempty"`
	Old string `json:"old,omitempty"`
	New string `json:"new,omitempty"`
}

type suggestions struct {
	Suggestions []string `json:"suggestions"`
}

type lookup struct {
	UserAbbs   []*Abbreviation
	GlobalAbbs []*Abbreviation
}
