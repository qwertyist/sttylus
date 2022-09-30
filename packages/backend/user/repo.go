package user

type UserRepository interface {
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id string) error

	GetSettings(id string) (*Settings, error)
	UpdateSettings(id string, s *Settings) error
	Login(email, password string) (*User, error)
}
