package user

import (
	"fmt"
)

func (s *userService) GetSettings(id string) (*Settings, error) {
	u, err := s.repo.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("service|GetSettings couldn't get user: %s", err)
	}
	settings := u.Settings
	return &settings, nil
}

func (s *userService) UpdateSettings(id string, settings *Settings) error {
	u, err := s.repo.GetUser(id)
	if err != nil {
		return fmt.Errorf("service|UpdateSettings couldn't get user: %s", err)
	}
	u.Settings = *settings
	err = s.UpdateUser(u)
	if err != nil {
		return fmt.Errorf("service|UpdateSettings couldn't update user: %s", err)
	}

	return nil
}

func (s *userService) ResetDefaults(id string) (*Settings, error) {
	u, err := s.repo.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("service|ResetDefaults couldn't get user: %s", err)
	}

	u.Settings = Defaults.Settings
	u.Subscriptions = Defaults.Subscriptions
	err = s.UpdateUser(u)
	if err != nil {
		return nil, fmt.Errorf("service|ResetDefaults couldn't update user: %s", err)
	}
	return &u.Settings, nil
}
