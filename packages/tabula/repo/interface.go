package repo

import "github.com/qwertyist/tabula/session"

type Repository interface {
	Init() error
	CreateUser(u session.User) error
	GetUser(id string) (*session.User, error)
	GetUsers() ([]*session.User, error)
	UpdateUser(u session.User) error
	DeleteUser(id string) error
}
