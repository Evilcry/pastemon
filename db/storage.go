package db

import (
	"log"
	"time"

	"github.com/evilcry/pastemon/configs"

	"github.com/asggo/store"
)

// GetStorageConnection function
// get storage connection
func GetStorageConnection(conf *configs.Config) *store.Store {
	var ds *store.Store
	var err error
	for tries := 1; tries < 20; tries += 2 {
		ds, err = store.NewStore(conf.DbFile)
		if err != nil {
			log.Printf("[-] Cannot open database: %s\n", err)
			time.Sleep(1 << uint(tries) * time.Millisecond)
		} else {
			break
		}
	}

	return ds
}

// InitStorage func
// Initialize storage buckets
func InitStorage(conf *configs.Config, s *store.Store) {
	for _, bucket := range conf.Buckets {
		s.CreateBucket(bucket)
	}
}

// CleanKeys func
// purge keys after each iteration
func CleanKeys(conf *configs.Config) {
	now := time.Now()
	max := time.Duration(conf.MaxTime) * time.Second

	for key, _ := range conf.Keys {
		if now.Sub(conf.Keys[key]) > max {
			delete(conf.Keys, key)
		}
	}
}
