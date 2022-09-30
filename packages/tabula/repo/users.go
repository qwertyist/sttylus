package repo

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/qwertyist/tabula/session"
)

func (r *repo) GetUsers() ([]*session.User, error) {
	var uu []*session.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.ForEach(func(k, v []byte) error {
			var u *session.User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return fmt.Errorf("repo.GetUsers failed unmarshalling: %s", err)
			}
			uu = append(uu, u)
			return nil
		})
		return err
	})
	return uu, err
}

func (r *repo) GetUser(id string) (*session.User, error) {
	var user *session.User
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get([]byte(id))
		if v != nil {
			err := json.Unmarshal(v, &user)
			if err != nil {
				return fmt.Errorf("repo.GetUser failed marshalling: %s", err)
			}
			return nil
		}
		return fmt.Errorf("repo.GetUser couldn't find user with id: %s", id)

	})
	return user, err
}

func (r *repo) CreateUser(u session.User) error {
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

func (r *repo) UpdateUser(u session.User) error {
	encoded, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("repo|CreateUser couldn't marshal user struc, err: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Put([]byte(u.ID), encoded)
		if err != nil {
			return fmt.Errorf("repo.UpdateUser couldn't put: %s", err.Error())
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
			return fmt.Errorf("repo.DeleteUser failed: %s", err)
		}
		return nil
	})
	return err
}
