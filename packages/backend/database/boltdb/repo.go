package boltdb

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/botvid/webapp/repository"
)

type repo struct {
	bolt *bolt.DB
}

// NewBoltRepository returns a repo connected to provided bolt db
func NewBoltRepository(bolt *bolt.DB) repository.Repository {
	return &repo{
		bolt,
	}
}

func (r *repo) Init() error {
	err := r.bolt.Update(func(tx *bolt.Tx) error {
		buckets := []string{"users", "abbreviations", "lists", "documents", "commends", "sessions"}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("bolt create bucket [%s]: %s", bucket, err)
			}
			log.Printf("Successfully created bucket: %s", bucket)
		}

		return nil

	})
	return err
}
