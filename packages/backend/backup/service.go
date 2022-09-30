package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/botvid/webapp/repository"
)

type BackupService interface {
	DumpUsers() error
	DumpAbbs() error
	DumpManuscripts() error
	DumpAll() error
	RestoreUsers() error
	RestoreAbbs() error
	RestoreManuscripts() error
	RestoreAll() error
}

type backupService struct {
	repo repository.Repository
}

func NewBackupService(r repository.Repository) BackupService {
	return &backupService{
		r,
	}
}

func (s *backupService) DumpUsers() error {
	uu, err := s.repo.GetUsers()
	if err != nil {
		return fmt.Errorf("service|DumpUsers couldn't get users: %q", err)
	}
	file, _ := os.OpenFile("users_dump.json", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	log.Printf("Got users: %v\n", uu)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(uu)
	if err != nil {
		return fmt.Errorf("service|DumpUsers couldn't encode to file %s: %q", file.Name, err)
	}

	return nil
}

func (s *backupService) DumpAbbs() error {
	ll, err := s.repo.GetLists(nil)
	if err != nil {
		return fmt.Errorf("service|DumpAbbs couldn't get lists: %q", err)
	}
	file, _ := os.OpenFile("abbs_dump.json", os.O_CREATE, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(ll)

	return nil
}

func (s *backupService) DumpManuscripts() error {

	return nil
}

func (s *backupService) DumpAll() error {
	err := s.DumpUsers()
	if err != nil {
		return err
	}

	s.DumpAbbs()
	if err != nil {
		return err
	}

	s.DumpManuscripts()
	if err != nil {
		return err
	}

	return nil
}

func (s *backupService) RestoreAbbs() error {
	return nil
}
func (s *backupService) RestoreManuscripts() error {
	return nil
}
func (s *backupService) RestoreUsers() error {
	return nil
}
func (s *backupService) RestoreAll() error { return nil }
