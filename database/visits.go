package database

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ngalayko/highloadcup/schema"
)

// CreateVisits creates given visits objects in db
func (db *DB) CreateVisits(visits *schema.Visits) error {
	for _, visit := range visits.Visits {
		if err := db.CreateOrUpdate(visit); err != nil {
			return err
		}
	}

	return nil
}

// LoadAllVisits return all visits from db
func (db *DB) LoadAllVisits() (visits map[uint32]*schema.Visit, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(schema.VisitsBucketName)

		visits = map[uint32]*schema.Visit{}

		b.ForEach(func(k, v []byte) error {

			visit := &schema.Visit{}
			if err := json.Unmarshal(v, visit); err != nil {
				return err
			}

			visits[visit.ID] = visit
			return nil
		})
		return nil
	})
	return
}
