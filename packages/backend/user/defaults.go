package user

import (
	"log"
	"time"
)

var Defaults = User{}

func (s *userService) CreateOrUpdateDefaultUser() error {
	u, _ := s.repo.GetUser("0")
	if u == nil {
		u.ID = "0"
		u.Created = time.Now()
		if u.Role == "" {
			u.Role = UserNormal
		}
		err := s.repo.CreateUser(u)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	u.Updated = time.Now()

	u.Settings = Defaults.Settings
	u.Subscriptions = Defaults.Subscriptions

	err := s.repo.UpdateUser(u)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
