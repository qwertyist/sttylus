package session

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/qwertyist/tabula/ws"
)

var validate *validator.Validate

type SessionService interface {
	CreateSession(sess Session) (*Session, error)
	GetSession(id string) (*Session, error)
	GetSessions(ids []string) ([]*Session, error)
	UpdateSession(sess Session) (*Session, error)
	DeleteSession(id string) (*Session, error)
	CreateUser(u User) *User
	GetUser(u User) *User
	GetUsers() []*User
	GetCaption(id string) string
	SetAuthToken(token string)
	CheckAuthToken(token string) bool
	ResetSessions() int
}

type sessionService struct {
	AuthToken string
	repo      Repository
	Sessions  map[string]*Session
	pools     map[string]*ws.Pool
	UserIDs   []string
	Users     []User
}

func NewSessionService(repo Repository, pools map[string]*ws.Pool) SessionService {
	validate = validator.New()
	return &sessionService{repo: repo, Sessions: make(map[string]*Session), pools: pools}
}

func (s *sessionService) Validate(sess Session) *ErrorLog {
	errorlog := &ErrorLog{}
	err := validate.Struct(sess)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			validationError := createValidationError(e)
			errorlog.Errors = append(errorlog.Errors, validationError)
		}
		return errorlog
	}
	return nil
}

func (s *sessionService) CreateSession(sess Session) (*Session, error) {
	e := &ErrorLog{}
	now := time.Now()
	if sess.ID == "" || len(sess.ID) < 6 || len(sess.ID) > 16 {
		sess.ID = createIDNumber(8)
	} else {
		_, ok := s.Sessions[sess.ID]
		//log.Println(ok)
		if ok {
			return nil, fmt.Errorf("%s", DuplicateSessionID)
		}
	}
	e.ID = sess.ID

	if sess.Type < 0 {
		sess.Type = Tabula
		err := &Error{Error: WarningNoSessionType, Time: now}
		e.Errors = append(e.Errors, err)
	}
	if sess.From.IsZero() || sess.To.IsZero() {
		sess.From = time.Now()
		sess.To = time.Now().Add(3 * time.Hour)
	}
	if sess.To.Before(sess.From) {
		err := &Error{Error: InvalidTimeToBeforeFrom, Time: now}
		e.Errors = append(e.Errors, err)
		sess.ErrorLog = e
		return &sess, fmt.Errorf("%s", string(InvalidTimeToBeforeFrom))
	}

	if sess.From.Before(now.AddDate(0, 0, -1)) {
		err := &Error{Error: InvalidTimeFromBeforeToday, Time: now}
		e.Errors = append(e.Errors, err)
		sess.ErrorLog = e
		return &sess, fmt.Errorf("%s", InvalidTimeFromBeforeToday)
	}

	if sess.Name == "" {
		if sess.Type == Tabula {
			sess.Name = "Distanstolkning " + sess.ID
		} else if sess.Type == Zoom {
			sess.Name = "Zoomtolkning " + sess.ID
		} else {
			sess.Name = "Distanstextning " + sess.ID
		}
	}
	if sess.Ref == "" && len(sess.Clients) == 0 && len(sess.Itprs) == 0 {
		err := &Error{Error: WarningNoReference, Time: now}
		e.Errors = append(e.Errors, err)
	}
	s.Validate(sess)
	s.Sessions[sess.ID] = &sess
	/*	err = s.repo.CreateSession(sess)
		if err != nil {
			log.Printf("yikes:", err)
		}
	*/
	sess.ErrorLog = e
	return &sess, nil
}

func (s *sessionService) GetSessions(ids []string) ([]*Session, error) {
	e := &Error{Time: time.Now()}
	ss := []*Session{}
	//log.Println("Get public sessions")
	for _, sess := range s.Sessions {
		if sess.Public {
			ss = append(ss, sess)
		}
	}
	for _, id := range ids {
		for _, sess := range s.Sessions {
			if sess.Ref == id {
				ss = append(ss, sess)
			}
			for _, c := range sess.Clients {
				if c == id {
					ss = append(ss, sess)
				}
			}
			for _, i := range sess.Itprs {
				if i == id {
					ss = append(ss, sess)
				}
			}
		}
	}
	if e.Error == "" {
		return ss, nil
	}

	return nil, fmt.Errorf("%s", e.Error)
}

func (s *sessionService) GetSession(id string) (*Session, error) {
	sess, ok := s.Sessions[id]
	if !ok {
		return nil, fmt.Errorf("%s", SessionNotFound)
	}

	return sess, nil
}

func (s *sessionService) UpdateSession(sess Session) (*Session, error) {
	target, err := s.GetSession(sess.ID)
	if target == nil {
		return target, err
	}
	errorlog := s.Validate(sess)
	if errorlog != nil {
		sess.ErrorLog = errorlog
		return &sess, fmt.Errorf("%s", ValidationError)
	}
	s.Sessions[sess.ID] = &sess
	return &sess, nil

}

func (s *sessionService) DeleteSession(id string) (*Session, error) {
	target, _ := s.GetSession(id)
	if target != nil {
		delete(s.Sessions, id)
		return target, nil
	}
	return nil, fmt.Errorf("%s", SessionNotFound)
}
