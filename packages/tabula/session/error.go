package session

import (
	"net/http"
	"time"
)

type ErrorMessage string

const (
	FaultyRequest              ErrorMessage = "Faulty REST request"
	DuplicateSessionID                      = "Can't create a session with a duplicated ID"
	UnauthorizedAccess                      = "Unauthorized access"
	NotEnoughInformation                    = "Not enough information provided"
	InvalidTimeToBeforeFrom                 = "Invalid provided time: To and From are inverted"
	InvalidTimeFromBeforeToday              = "Invalid provided time: From is before today"
	WarningNoReference                      = "Warning: No reference selected, no information will be sent out"
	WarningNoSessionType                    = "Warning: No remote session type selected, falling back to Tabula"
	SessionNotFound                         = "No session found with provided ID"
	ValidationError
)

type Error struct {
	Time  time.Time
	Error ErrorMessage
}

type ErrorLog struct {
	ID       string
	Errors   []*Error
	SaveText bool

	Type  RemoteType
	Ref   string
	Iptrs []string

	Req *http.Request
}

func createValidationError(err error) *Error {
	return &Error{
		Time:  time.Now(),
		Error: ValidationError + ": " + ErrorMessage(err.Error()),
	}
}
