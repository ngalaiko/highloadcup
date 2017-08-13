package database

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ngalayko/highloadcup/schema"
)

// CreateLocation creates given location objects
func (db *DB) CreateLocations(locations *schema.Locations) error {
	for _, location := range locations.Locations {
		if err := db.CreateOrUpdate(location); err != nil {
			return err
		}
	}

	return nil
}

// LoadAllLocations return all locations from db
func (db *DB) LoadAllLocations() (locations map[uint32]*schema.Location, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(schema.LocationsBucketName)

		locations = map[uint32]*schema.Location{}

		b.ForEach(func(k, v []byte) error {
			location := &schema.Location{}
			if err := json.Unmarshal(v, location); err != nil {
				return err
			}

			locations[location.ID] = location
			return nil
		})
		return nil
	})
	return
}
