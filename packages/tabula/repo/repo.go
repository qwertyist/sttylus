package repo

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type repo struct {
	bolt *bolt.DB
}

func NewRepository(bolt *bolt.DB) Repository {
	r := &repo{bolt}
	err := r.Init()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func (r *repo) Init() error {
	err := r.bolt.Update(func(tx *bolt.Tx) error {
		buckets := []string{"users", "sessions"}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("bold couldn't create bucket [%s]: %s", bucket, err)
			}
			log.Printf("Successfully created bucket: %s", bucket)
		}
		return nil
	})
	return err
}

func OpenBoltDB(fileName string) *bolt.DB {
	log.Println("Opening BoltDB", fileName)
	bolt, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		log.Fatalf("Failed opening BoltDB %s, err: %s\n", fileName, err)
	}
	return bolt
}
