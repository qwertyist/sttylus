package boltdb

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/botvid/webapp/user"
)

func (r *repo) CreateUser(u *user.User) error {
	log.Println("Create User:", u.Settings.Font)
	encoded, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("repo|CreateUser couldn't marshal user struc, err: %s", err)
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Put([]byte(u.ID), encoded)
		if err != nil {
			return fmt.Errorf("repo|CreateUser couldn't put: %s", err.Error())
		}
		return nil
	})
	return nil
}

func (r *repo) GetUser(id string) (*user.User, error) {
	var user *user.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get([]byte(id))
		if v != nil {
			err := json.Unmarshal(v, &user)
			if err != nil {
				return fmt.Errorf("repo|GetUser failed unmarshalling: %s", err)
			}
			return nil
		}
		return fmt.Errorf("repo|GetUser couldn't find user with id: %s", id)

	})
	return user, err
}

func (r *repo) GetUserByEmail(email string) (*user.User, error) {
	var u *user.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.ForEach(func(k, v []byte) error {
			var tmp *user.User
			err := json.Unmarshal(v, &tmp)
			if err != nil {
				return fmt.Errorf("repo|GetUserByEmail failed unmarshalling: %s", err)
			}
			if tmp.Email == email {
				u = tmp
				return nil
			}
			return nil
		})
		return err
	})
	return u, err
}

func (r *repo) GetUsers() ([]*user.User, error) {
	var uu []*user.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.ForEach(func(k, v []byte) error {
			var u *user.User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return fmt.Errorf("repo|GetUsers failed unmarshalling: %s", err)
			}
			uu = append(uu, u)
			return nil
		})
		return err
	})
	return uu, err
}

func (r *repo) UpdateUser(u *user.User) error {
	encoded, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("repo|UpdateUser failed marshalling: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Put([]byte(u.ID), encoded)
		if err != nil {
			return fmt.Errorf("repo|UpdateUser failed: %s", err)
		}
		return nil
	})

	return err
}

func (r *repo) DeleteUser(id string) error {
	err := r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Delete([]byte(id))
		if err != nil {
			return fmt.Errorf("repo|DeleteUser failed: %s", err)
		}
		return nil
	})

	return err
}

func (r *repo) Login(email, password string) (*user.User, error) {
	var u *user.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.ForEach(func(k, v []byte) error {
			var tmp *user.User
			err := json.Unmarshal(v, &tmp)
			if err != nil {
				return fmt.Errorf("repo|Login failed unmarshalling: %s", err)
			}
			if tmp.Email == email {
				u = tmp
				return nil
			}
			return nil
		})
		return err
	})
	return u, err
}

func (r *repo) GetSettings(id string) (*user.Settings, error) {
	return nil, nil
}

func (r *repo) UpdateSettings(id string, settings *user.Settings) error {
	return nil
}
