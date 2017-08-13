package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/ngalayko/highloadcup/config"
	"github.com/ngalayko/highloadcup/helper"
	"github.com/ngalayko/highloadcup/schema"
)

const (
	ctxKey ctxDbKey = "ctx_key_for_db"
)

type ctxDbKey string

type DB struct {
	*bolt.DB
}

func NewContext(ctx context.Context, db interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := db.(*DB); !ok {
		db = NewDB(ctx)
	}

	return context.WithValue(ctx, ctxKey, db)
}

func FromContext(ctx context.Context) *DB {
	if db, ok := ctx.Value(ctxKey).(*DB); ok {
		return db
	}

	return NewDB(ctx)
}

func NewDB(ctx context.Context) *DB {
	cfg := config.FromContext(ctx)

	db, err := bolt.Open(cfg.DbPath, 0600, nil)
	if err != nil {
		log.Panicf("erorr when open bolt.db: %s", err)
	}

	if err := initDb(db); err != nil {
		log.Panicf("erorr when init bolt.db: %s", err)
	}

	return &DB{db}
}

func initDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range schema.Buckets {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return fmt.Errorf("CreateOrUpdate bucketName error: %s", err)
			}
		}

		return nil
	})
}

// CreateOrUpdate creates entity in db
func (db *DB) CreateOrUpdate(entity schema.IEntity) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(entity.Bucket())

		if err := b.Put(entity.ByteID(), entity.Bytes()); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetBytes(entity schema.Entity, id uint32) (result []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(schema.BucketsMap[entity])

		result = b.Get(helper.Itob(id))

		return nil
	})

	return
}

func (db *DB) Get(entity schema.Entity, id uint32, toPtr interface{}) error {
	data, err := db.GetBytes(entity, id)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, toPtr); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetIds(entity schema.Entity, ids []uint32, toPtr interface{}) error {
	return db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(schema.BucketsMap[entity])

		var result []byte
		for _, id := range ids {
			if len(result) > 0 {
				result = append(result, ',')
			}

			result = append(result, b.Get(helper.Itob(id))...)
		}

		toParse := []byte{'['}
		toParse = append(toParse, result...)
		toParse = append(toParse, ']')

		return json.Unmarshal(toParse, toPtr)
	})
}
