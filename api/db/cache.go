package db

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
)

type Cache struct {
	DB *bolt.DB
}

func InitBoltDB() error {
	cache, err := Open()
	if err != nil {
		return err
	}
	cache, err = CreateBucket(cache.DB)
	if err != nil {
		return err
	}
	return cache.Close()
}

func Open() (*Cache, error) {
	db, err := bolt.Open("cache.db", 0666, nil)
	if err != nil {
		return nil, err
	}

	cache := &Cache{DB: db}

	return cache, nil
}

func CreateBucket(db *bolt.DB) (*Cache, error) {
	cache := &Cache{DB: db}
	// Create a new bucket for caching data
	err := cache.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("dashboard"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// func NewCacheTx(cache *Cache) (err error) {
// 	cache.Tx, err = cache.DB.Begin(true)
// 	return
// }

func (c *Cache) Set(key string, value interface{}) error {
	return c.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("dashboard"))
		buf, err := encode(value)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), buf)
	})
}

func (c *Cache) Get(key string, value interface{}) error {
	return c.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("dashboard"))
		buf := bucket.Get([]byte(key))
		if buf == nil {
			return nil
		}
		return decode(buf, value)
	})
}

func (c *Cache) Delete(key string) error {
	return c.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("cache"))
		return bucket.Delete([]byte(key))
	})
}

func (c *Cache) Close() error {
	return c.DB.Close()
}

func encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes(), err
}

func decode(buf []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(v)
}
