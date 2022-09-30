package repository

import (
	"github.com/botvid/webapp/abbreviation"
	"github.com/botvid/webapp/document"
	"github.com/botvid/webapp/user"
)

type Repository interface {
	Init() error
	CreateUser(u *user.User) error
	GetUser(id string) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)

	GetUsers() ([]*user.User, error)
	UpdateUser(u *user.User) error
	DeleteUser(id string) error

	Login(email, password string) (*user.User, error)

	GetSettings(id string) (*user.Settings, error)
	UpdateSettings(id string, settings *user.Settings) error

	CreateAbb(listID string, abb *abbreviation.Abbreviation) error
	CreateManyAbbs(listID string, abbs []*abbreviation.Abbreviation) error
	GetAbb(listID, abb string) (*abbreviation.Abbreviation, error)
	GetAbbs(listID string) ([]*abbreviation.Abbreviation, error)
	UpdateAbb(listID string, abb *abbreviation.Abbreviation) error
	DeleteAbb(listID, abb string) (*abbreviation.Abbreviation, error)
	DeleteAbbByID(listID, ID string) (*abbreviation.Abbreviation, error)

	ImportAbbsToList(listID string, abbs []*abbreviation.Abbreviation) error

	CreateList(list *abbreviation.List) error
	CreateManyLists(lists []*abbreviation.List) error
	GetList(listID string) (*abbreviation.List, error)
	GetLists(listIDs []string) ([]*abbreviation.List, error)
	GetUserLists(userID string) ([]*abbreviation.List, error)
	UpdateList(list *abbreviation.List) error
	DeleteList(listID string) (*abbreviation.List, error)

	GetDoc(ID string) (*document.Document, error)
	GetDocs(userID string) ([]*document.Document, error)
	CreateDoc(doc *document.Document) error
	CreateManyDocs(docs []*document.Document) error

	UpdateDoc(doc *document.Document) error
	DeleteDoc(ID string) error
}
