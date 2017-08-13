package database

import "github.com/ngalayko/highloadcup/schema"

func (db *DB) GetVisits(ids []uint32) (map[uint32]*schema.Visit, error) {
	val, err := db.GetIds(schema.EntityVisits, ids)
	if err != nil {
		return nil, err
	}

	visits := map[uint32]*schema.Visit{}
	for _, v := range val {
		if visit, ok := v.(*schema.Visit); ok {
			visits[visit.ID] = visit
		}
	}

	return visits, nil
}
