package database

import (
	"fmt"

	"github.com/ngalayko/highloadcup/schema"
)

func (db *DB) GetLocation(id uint32) (*schema.Location, error) {
	val, err := db.Get(schema.EntityLocations, id)
	if err != nil {
		return nil, err
	}

	if location, ok := val.(*schema.Location); ok {
		return location, nil
	}

	return nil, fmt.Errorf("error on casting %v to location", val)
}

func (db *DB) GetAllLocations() (map[uint32]*schema.Location, error) {
	m := db.mapByEntity(schema.EntityLocations)

	locations := map[uint32]*schema.Location{}
	m.Range(func(k, v interface{}) bool {
		if location, ok := v.(*schema.Location); ok {
			locations[location.ID] = location
		}

		return true
	})

	return locations, nil
}
