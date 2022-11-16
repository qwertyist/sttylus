package user

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/botvid/webapp/abbreviation"
	"github.com/botvid/webapp/document"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(l loginRequest) (*User, error)
	Auth(id string, role Role) error
	RegisterWeb(l loginRequest) (*User, error)
	CreateUser(u *User) (*User, error)
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]*User, error)
	GetMachineID() (string, error)
	UpdateUser(u *User) error
	DeleteUser(id string) error
	CreateDemoUser(u *User) error
	GetSettings(id string) (*Settings, error)
	UpdateSettings(id string, settings *Settings) error
	ResetDefaults(id string) (*Settings, error)
	ResetPassword(email, password string) error
	UpdatePassword(id, password string) error

	Offline() bool
	LocalLogin(l loginRequest) (*User, error)
	SyncUser(l loginRequest) (*User, error)
	SyncUserLists(user *User) ([]*abbreviation.List, error)
	SyncUserAbbs(userID string, lists []*abbreviation.List) ([]*abbreviation.Abbreviation, error)
	SyncUserDocs(user *User) ([]*document.Document, error)
}

type userService struct {
	repo       UserRepository
	abbService abbreviation.AbbService
	docService document.DocService
	mode       string
}

//NewUserService returns a userService connected to the repo db
func NewUserService(repo UserRepository, mode string, abbService abbreviation.AbbService, docService document.DocService) UserService {
	Defaults.Settings.Font.Family = "Times New Roman"
	Defaults.Settings.Font.Size = 32
	Defaults.Settings.Font.LineHeight = 1
	Defaults.Settings.Font.ColorID = 0
	Defaults.Settings.Font.Background = "#000000"
	Defaults.Settings.Font.Foreground = "#FFFFFF"
	Defaults.Settings.Behaviour.CapitalizeOnNewLine = true
	Defaults.Settings.Font.Margins.Top = 20
	Defaults.Settings.Font.Margins.Right = 20
	Defaults.Settings.Font.Margins.Bottom = 20
	Defaults.Settings.Font.Margins.Left = 20

	return &userService{
		repo:       repo,
		mode:       mode,
		abbService: abbService,
		docService: docService,
	}
}

//GetUsers returns a list of User IDs
func GetUsers(s UserService) ([]string, error) {
	var userIDs []string
	uu, err := s.GetUsers()
	for _, u := range uu {
		userIDs = append(userIDs, u.ID)
	}
	return userIDs, err
}

func (s *userService) CreateUser(u *User) (*User, error) {
	if u.Email != "" {
		exist, _ := s.GetUserByEmail(u.Email)
		if exist != nil {
			return nil, fmt.Errorf("e-mail taken")
		}
		if u.ID != "" {
			return s.GetUser(u.ID)
		}
		u.ID = uuid.New().String()
		u.Created = time.Now()
		u.Updated = time.Now()

		if u.Role == "" {
			u.Role = UserNormal
		}

		u.Settings = Defaults.Settings
		err := s.repo.CreateUser(u)

		return u, err
	}
	return u, fmt.Errorf("no email adress provided")
}

func (s *userService) CreateDemoUser(u *User) error {
	return nil
}

func (s *userService) GetUser(id string) (*User, error) {
	u, err := s.repo.GetUser(id)
	return u, err
}

func (s *userService) GetUserByEmail(email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(s.mode)
	if user == nil {
		log.Println("No user found, try online")
		if s.mode == "desktop" {
			resp, err := http.Get("https://sttylus.se/api2/registered/" + email)
			if err != nil {
				return nil, fmt.Errorf("couldn't reach server %s", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(body))
			if string(body) == "login" {
				return &User{}, nil
			}
			return nil, nil
		}
		log.Println(" online")
		return nil, fmt.Errorf("No user with that email")
	}
	return user, nil
}

func (s *userService) GetUsers() ([]*User, error) {
	return s.repo.GetUsers()

}

func (s *userService) UpdateUser(u *User) error {
	old, err := s.GetUser(u.ID)
	log.Println("old user settings:", old.Settings.Font)
	if err != nil {
		return err
	}

	log.Print("new user settings:", u.Settings.Font)
	u.PasswordHash = old.PasswordHash
	u.Salt = old.Salt

	u.Updated = time.Now()
	err = s.repo.UpdateUser(u)
	if err != nil {
		log.Println("Could update user:", err)
	}
	return err
}

func (s *userService) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	return err
}

func (s *userService) UpdatePassword(userID, pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service|UpdatePassword couldn't generate hash: %q", err)

	}
	if userID != "" {

		u, err := s.repo.GetUser(userID)
		if err != nil {
			return fmt.Errorf("service|UpdatePassword couldn't find the user with provided ID: %s", userID)
		}
		if pwd != "" {
			u.PasswordHash = hash
			u.ResetPassword = false
		} else {
			u.ResetPassword = true
			u.PasswordHash = nil
		}

		err = s.repo.UpdateUser(u)
		if err != nil {
			return fmt.Errorf("service|UpdatePassword failed updating user: %q", err)
		}

		return nil
	}
	return fmt.Errorf("service|UpdatePassword - no user id provided")

}

func (s *userService) ResetPassword(email, password string) error {
	u, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}
	err = s.UpdatePassword(u.ID, password)
	if err != nil {
		return err
	}
	return nil
}
