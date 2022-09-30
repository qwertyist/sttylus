package boltdb

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/botvid/webapp/abbreviation"
)

func (r *repo) GetAbb(listID, abb string) (*abbreviation.Abbreviation, error) {
	var a *abbreviation.Abbreviation

	err := r.bolt.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		v := l.Get([]byte(abb))
		if v != nil {
			err := json.Unmarshal(v, &a)
			if err != nil {
				return fmt.Errorf("repo|GetAbb failed unmarshalling: %s", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return a, err
}

func (r *repo) GetAbbs(listID string) ([]*abbreviation.Abbreviation, error) {
	var abbs []*abbreviation.Abbreviation
	var err error
	err = r.bolt.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		if l != nil {
			err := l.ForEach(func(k, v []byte) error {
				var abb *abbreviation.Abbreviation
				err := json.Unmarshal(v, &abb)
				if err != nil {
					return fmt.Errorf("repo|GetAbbs - (by listID) failed unmarshalling :%s", err)
				}
				abbs = append(abbs, abb)
				return nil
			})
			return err
		}
		return fmt.Errorf("repo|GetAbbs - no list bucket found")
	})
	return abbs, err
}

func (r *repo) CreateAbb(listID string, abb *abbreviation.Abbreviation) error {
	encoded, err := json.Marshal(abb)

	if err != nil {
		return fmt.Errorf("repo|CreateAbb error marshalling: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		err := l.Put([]byte(abb.Abb), encoded)

		if err != nil {
			return fmt.Errorf("repo|CreateAbb failed: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) CreateManyAbbs(listID string, abbs []*abbreviation.Abbreviation) error {
	return r.bolt.Batch(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		for _, abb := range abbs {
			encoded, err := json.Marshal(abb)
			if err != nil {
				return fmt.Errorf("repo|CreateManyAbbs error marshalling: %s", err)
			}
			err = l.Put([]byte(abb.Abb), encoded)
			if err != nil {
				return fmt.Errorf("repo|CreateManyAbbs failed: %s", err)
			}
		}
		return nil
	})

}

func (r *repo) UpdateAbb(listID string, abb *abbreviation.Abbreviation) error {
	encoded, err := json.Marshal(abb)

	if err != nil {
		return fmt.Errorf("repo|UpdateAbb error marshalling: %s", err)
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		found := l.Get([]byte(abb.Abb))
		if found == nil {
			return nil
		}

		err := l.Put([]byte(abb.Abb), encoded)
		if err != nil {
			return fmt.Errorf("repo|UpdateAbb failed: %s", err)
		}

		return nil
	})

	return err
}

func (r *repo) DeleteAbb(listID, abb string) (*abbreviation.Abbreviation, error) {
	a, err := r.GetAbb(listID, abb)
	if err != nil {
		return a, fmt.Errorf("repo|DeleteAbb pre-fetch failed: %s", err.Error())
	}
	if a == nil {
		return a, fmt.Errorf("repo|DeleteAbb didnt' find abb: %s", abb)
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		err := l.Delete([]byte(abb))
		if err != nil {
			return fmt.Errorf("repo|DeleteAbb failed: %s", err.Error())
		}
		return nil
	})
	return a, err
}

func (r *repo) DeleteAbbByID(listID, ID string) (*abbreviation.Abbreviation, error) {
	var a *abbreviation.Abbreviation
	err := r.bolt.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		b := root.Bucket([]byte(listID))
		err := b.ForEach(func(k, v []byte) error {
			var abb *abbreviation.Abbreviation
			err := json.Unmarshal(v, &abb)
			if err != nil {
				return fmt.Errorf("repo|getLists (get all) failed unmarshalling: %s", err)
			}
			if abb.ID == ID {
				err := b.Delete([]byte(abb.Abb))
				a = abb
				if err != nil {
					return fmt.Errorf("repo|DeleteAbb failed: %s", err.Error())
				}
			}
			return nil
		})

		return err
	})
	return a, err
}

func (r *repo) GetList(id string) (*abbreviation.List, error) {
	var list *abbreviation.List
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lists"))
		meta := b.Get([]byte(id))
		err := json.Unmarshal(meta, &list)
		if err != nil {
			return fmt.Errorf("repo|GetList failed unmarshalling: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repo) GetLists(listIDs []string) ([]*abbreviation.List, error) {
	var lists []*abbreviation.List
	if len(listIDs) > 0 {
		err := r.bolt.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("lists"))
			for _, ID := range listIDs {
				var list *abbreviation.List
				meta := b.Get([]byte(ID))
				if meta != nil {
					err := json.Unmarshal(meta, &list)
					if err != nil {
						return fmt.Errorf("repo|getLists (by IDs) failed unmarshalling: %s", err)
					}

					lists = append(lists, list)
				}
			}
			return nil
		})
		return lists, err
	}

	err := r.bolt.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("lists"))

		err := b.ForEach(func(k, v []byte) error {
			var list *abbreviation.List
			err := json.Unmarshal(v, &list)
			if err != nil {
				return fmt.Errorf("repo|getLists (get all) failed unmarshalling: %s", err)
			}
			lists = append(lists, list)
			return nil
		})
		return err
	})

	return lists, err
}

func (r *repo) GetUserLists(userID string) ([]*abbreviation.List, error) {
	var lists []*abbreviation.List

	err := r.bolt.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("lists"))

		err := b.ForEach(func(k, v []byte) error {
			var list *abbreviation.List
			err := json.Unmarshal(v, &list)
			if err != nil {
				return fmt.Errorf("repo|getLists (get all) failed unmarshalling: %s", err)
			}
			if list.Creator == userID {
				lists = append(lists, list)
			}
			return nil
		})
		return err
	})

	return lists, err
}

func (r *repo) CreateList(list *abbreviation.List) error {
	encoded, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("repo|CreateList failed marshalling request: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lists"))

		root := tx.Bucket([]byte("abbreviations"))
		_, err := root.CreateBucketIfNotExists([]byte(list.ID))
		if err != nil {
			return fmt.Errorf("repo|CreateList failed creating list's abb bucket: %s", err)
		}

		err = b.Put([]byte(list.ID), encoded)

		if err != nil {
			return fmt.Errorf("repo|CreateList failed: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) CreateManyLists(lists []*abbreviation.List) error {
	return r.bolt.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lists"))
		root := tx.Bucket([]byte("abbreviations"))
		for _, list := range lists {
			encoded, err := json.Marshal(list)
			if err != nil {
				return fmt.Errorf("repo|CreateManyLists error marshalling: %s", err)
			}
			_, err = root.CreateBucketIfNotExists([]byte(list.ID))
			if err != nil {
				return fmt.Errorf("repo|CreateManyLists failed creating list's abb bucket: %s", err)
			}
			err = b.Put([]byte(list.ID), encoded)
			if err != nil {
				return fmt.Errorf("repo|CreateManyLists failed: %s", err)
			}
		}
		return nil
	})
}
func (r *repo) UpdateList(list *abbreviation.List) error {
	encoded, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("repo|UpdateList error marshalling: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lists"))
		err := b.Put([]byte(list.ID), encoded)
		if err != nil {
			return fmt.Errorf("repo|UpdateList failed: %s", err)
		}
		return nil
	})

	return err
}

func (r *repo) DeleteList(listID string) (*abbreviation.List, error) {
	l, err := r.GetList(listID)
	if err != nil {
		return nil, fmt.Errorf("repo|DeleteList pre-fetch failed: %s", err)
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lists"))
		err := b.Delete([]byte(listID))
		if err != nil {
			return fmt.Errorf("repo|DeleteList failed: %s", err)
		}
		return nil
	})
	if err != nil {
		return l, err
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		err = root.DeleteBucket([]byte(listID))
		if err != nil {
			return fmt.Errorf("repo|DeleteList failed deleting list's abb bucket: %s", err)
		}

		return nil
	})

	return l, err
}

func (r *repo) ImportAbbsToList(listID string, abbs []*abbreviation.Abbreviation) error {
	err := r.bolt.Batch(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("abbreviations"))
		l := root.Bucket([]byte(listID))
		for _, abb := range abbs {
			encoded, err := json.Marshal(abb)
			if err != nil {
				return fmt.Errorf("repo|ImportAbbsToList failed marshalling: %q", err)
			}
			if err = l.Put([]byte(abb.Abb), encoded); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
