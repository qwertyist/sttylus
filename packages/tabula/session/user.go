package session

import (
	"log"
)

func (s *sessionService) GetUser(request User) *User {
	if request.ID != "" {
		u, err := s.repo.GetUser(request.ID)
		if err != nil {
			log.Fatal(err)
		}
		return u
	}
	uu, err := s.repo.GetUsers()
	if err != nil {
		log.Fatalf("session.GetUser failed: %s\n", err)
	}
	for _, u := range uu {
		if u.Email != "" {
			if u.Email == request.Email {
				return u
			}
		}
		if u.Phone != "" {
			if u.Phone == request.Phone {
				return u
			}
		}
	}
	return nil
}

func (s *sessionService) GetUsers() []*User {
	uu, err := s.repo.GetUsers()
	if err != nil {
		log.Println(err)
		return nil
	}
	return uu
}

func (s *sessionService) CreateUser(u User) *User {
	u.ID = createIDNumber(10)
	s.UserIDs = append(s.UserIDs, u.ID)
	s.Users = append(s.Users, u)
	err := s.repo.CreateUser(u)
	if err != nil {
		log.Println("yikes:", err)
	}
	return &u
}
