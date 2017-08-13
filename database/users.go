package database

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ngalayko/highloadcup/schema"
)

// CreateUsers creates given user objects in db
func (db *DB) CreateUsers(users *schema.Users) error {
	for _, user := range users.Users {
		if err := db.CreateOrUpdate(user); err != nil {
			return err
		}
	}

	return nil
}

// LoadAllUsers return all users from db
func (db *DB) LoadAllUsers() (users map[uint32]*schema.User, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(schema.UsersBucketName)

		users = map[uint32]*schema.User{}

		b.ForEach(func(k, v []byte) error {
			user := &schema.User{}
			if err := json.Unmarshal(v, user); err != nil {
				return err
			}

			users[user.ID] = user
			return nil
		})
		return nil
	})
	return
}
