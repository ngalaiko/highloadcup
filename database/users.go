package database

import (
	"fmt"

	"github.com/ngalayko/highloadcup/schema"
)

func (db *DB) GetUser(id uint32) (*schema.User, error) {
	val, err := db.Get(schema.EntityUsers, id)
	if err != nil {
		return nil, err
	}

	if user, ok := val.(*schema.User); ok {
		return user, nil
	}

	return nil, fmt.Errorf("error on casting %v to user", val)
}

func (db *DB) GetUsers(ids []uint32) (map[uint32]*schema.User, error) {
	val, err := db.GetIds(schema.EntityUsers, ids)
	if err != nil {
		return nil, err
	}

	usersMap := map[uint32]*schema.User{}
	for _, v := range val {
		if user, ok := v.(*schema.User); ok {
			usersMap[user.ID] = user
		}
	}

	return usersMap, nil
}
