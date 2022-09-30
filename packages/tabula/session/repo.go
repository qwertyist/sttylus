package session

type Repository interface {
	GetUser(id string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(u User) error
	UpdateUser(u User) error
	DeleteUser(id string) error
}
