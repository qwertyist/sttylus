package session

import (
	"time"
)

type RemoteType int

const (
	Tabula RemoteType = iota
	Zoom
	GETlivecap
)

type Session struct {
	ID       string    `json:"id" validate:"required,alphanum"`
	Name     string    `json:"name" validate:"required,alphanumunicode`
	From     time.Time `json:"from" validate:"required"`
	To       time.Time `json:"to" validate:"required"`
	Password string    `json:"password,omitempty" validate:"omitempty,alphanumunicode"`
	Public   bool      `json:"public"`

	Type      RemoteType `json:"type" validate:"gt=-1,lt=3"`
	Recurring bool       `json:"recurring"`

	SaveText bool     `json:"save_text"`
	Ref      string   `json:"ref"`
	Itprs    []string `json:"itprs"`
	Clients  []string `json:"clients"`

	Description string `json:"description"`
	MaxClients  int    `json:"max_clients"`
	Token       string `json:"token"`
	Breakout    string `json:"token"`

	ErrorLog *ErrorLog `json:"omitempty"`
}

type ClientType int

type User struct {
	ID          string `json:"id" validate: "required, alphanum"`
	Name        string `json:"name" validate: "required"`
	Email       string `json:"email" validate: "email"`
	Phone       string `json:"phone"`
	Creator     bool   `json:"creator"`
	Description string `json:"description"`
}
