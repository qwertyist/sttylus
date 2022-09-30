package user

import "time"

//UserRole represents the levels of authorization within the application
type Role string

const (
	UserNormal    Role = "user"
	UserModerator      = "moderator"
	UserAdmin          = "admin"
	UserTester         = "tester"
	UserClosed         = "closed"
)

type User struct {
	ID            string `json:"id,omitempty"`
	MachineID     string `json:"machine_id,omitempty"`
	LicenseKey    string `json:"license_key,omitempty"`
	PasswordHash  []byte `json:"password_hash,omitempty"`
	Salt          []byte `json:"salt,omitempty"`
	ResetPassword bool   `json:"reset_password"`
	Company       string `json:"company,omitempty"`
	Description   string `json:"description,omitempty"`
	Role          Role   `json:"role,omitempty"`
	Local         bool   `json:"local,omitempty"`

	Creator  bool      `json:"creator"`
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
	LastSync time.Time `json:"last_sync,omitempty"`
	active   bool      `json:"active"`

	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Inform string `json:"inform,omitempty"`

	Lists         []string      `json:"lists,omitempty"`
	Documents     []string      `json:"documents,omitempty"`
	Subscriptions Subscriptions `json:"subscriptions,omitempty"`
	Settings      Settings      `json:"settings,omitempty"`
}

type OldTextSettings struct {
}

type Settings struct {
	Font struct {
		Family     string  `json:"family,omitempty"`
		Size       int     `json:"size,omitempty"`
		LineHeight float32 `json:"lineHeight,omitempty"`

		ColorID      int    `json:"colorID"`
		Background   string `json:"background,omitempty"`
		Foreground   string `json:"foreground,omitempty"`
		CustomColors struct {
			Background string `json:"background,omitempty"`
			Foreground string `json:"foreground,omitempty"`
			Valid      bool   `json:"valid"`
		} `json:"customColors"`

		Margins struct {
			Top    int `json:"top,omitempty"`
			Right  int `json:"right,omitempty"`
			Bottom int `json:"bottom,omitempty"`
			Left   int `json:"left,omitempty"`
		} `json:"margins,omitempty"`
	} `json:"font,omitempty"`
	Behaviour struct {
		CapitalizeOnNewLine bool `json:"capitalizeOnNewLine"`
	} `json:"behaviour,omitempty"`
	SelectedLists struct {
		Standard string   `json:"standard,omitempty"`
		Addon    []string `json:"addon,omitempty"`
	} `json:"selectedLists,omitempty"`
	SelectedManuscripts []Manuscript `json:"selectedManuscripts"`
}

type Subscriptions struct {
	GlobalLists     []string `json:"global_lists,omitempty"`
	GlobalDocuments []string `json:"global_documents,omitempty"`

	UserLists     []string `json:"user_lists,omitempty"`
	UserDocuments []string `json:"user_documents,omitempty"`
}

type Manuscript struct {
	Abb string `json:"abb,omitempty"`
	ID  string `json:"id,omitempty"`
}
