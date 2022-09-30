package abbreviation

import "time"

// Abbreviation stores information about a single isolated abbreviation
// abbreviations are stored under a list or globally
type Abbreviation struct {
	ID      string `json:"id"`
	Abb     string `json:"abb"`
	Word    string `json:"word"`
	Creator string `json:"creator"`

	Comment string `json:"comment"`
	Remind  bool   `json:"remind"`
	ListID  string `json:"listId"`

	Updated time.Time `json:"updated"`
}

// List information and metadata
type List struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Checksum    string `json:"checksum"`

	Type    int    `json:"type"`
	Counter int    `json:"counter"`
	Locale  string `json:"locale"`
	Creator string `json:"creator"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"deleted"`
}
