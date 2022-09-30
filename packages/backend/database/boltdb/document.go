package boltdb

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/botvid/webapp/document"
)

func (r *repo) GetDoc(ID string) (*document.Document, error) {
	var doc *document.Document
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		v := b.Get([]byte(ID))
		if v != nil {
			err := json.Unmarshal(v, &doc)
			if err != nil {
				return fmt.Errorf("bolt GetDoc failed unmarshalling: %s", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return doc, err
}

func (r *repo) GetDocs(userID string) ([]*document.Document, error) {
	var docs []*document.Document
	err := r.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		err := b.ForEach(func(k, v []byte) error {
			var doc *document.Document
			err := json.Unmarshal(v, &doc)
			if doc.Creator == userID {
				if err != nil {
					return fmt.Errorf("bolt GetDocs (by userID) failed unmarshalling :%s", err)
				}

				docs = append(docs, doc)
			}
			return nil

		})
		return err
	})
	return docs, err
}

func (r *repo) GetGlobalDocs() ([]*document.Document, error) {
	var dd []*document.Document
	err := r.bolt.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("global"))
		docs := root.Bucket([]byte("documents"))
		if docs != nil {
			err := docs.ForEach(func(k, v []byte) error {
				var d *document.Document
				err := json.Unmarshal(v, &d)
				if err != nil {
					return fmt.Errorf("bolt GetGlobalDocs failed unmarshalling: %s", err)
				}
				dd = append(dd, d)
				return nil
			})
			return err
		}
		return nil
	})
	return dd, err
}

func (r *repo) CreateDoc(doc *document.Document) error {
	encoded, err := json.Marshal(doc)

	if err != nil {
		return fmt.Errorf("bolt CreateDoc error marshalling: %s", err)
	}

	err = r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		err := b.Put([]byte(doc.ID), encoded)
		if err != nil {
			return fmt.Errorf("bolt CreateDoc failed: %s", err)
		}
		return nil
	})
	return err
}

func (r *repo) CreateManyDocs(docs []*document.Document) error {
	return r.bolt.Batch(func(tx *bolt.Tx) error {
		for _, doc := range docs {
			encoded, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("bolt|CreateManyDocs error marshalling: %s", err)
			}
			if err != nil {
				return fmt.Errorf("bolt|CreateManyDocs error marshalling: %s", err)
			}
			b := tx.Bucket([]byte("documents"))

			err = b.Put([]byte(doc.ID), encoded)
			if err != nil {
				return fmt.Errorf("bolt|CreateManyDocs failed: %s", err)
			}
		}
		return nil
	})
}

func (r *repo) UpdateDoc(doc *document.Document) error {
	encoded, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("bolt UpdateDoc error marshalling: %s", err)
	}
	err = r.bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("documents"))
		err := b.Put([]byte(doc.ID), encoded)
		if err != nil {
			return fmt.Errorf("bolt UpdateDoc failed: %s", err)
		}
		return nil
	})
	return err
}

func (r *repo) DeleteDoc(ID string) error {
	err := r.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		err := b.Delete([]byte(ID))
		if err != nil {
			return fmt.Errorf("bolt DeleteDoc failed: %s", err)
		}
		return nil
	})
	return err
}
