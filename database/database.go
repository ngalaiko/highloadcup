package database

import (
	"context"
	"fmt"

	"golang.org/x/sync/syncmap"

	"github.com/ngalayko/highloadcup/schema"
)

const (
	ctxKey ctxDbKey = "ctx_key_for_db"
)

type ctxDbKey string

type DB struct {
	usersMap     *syncmap.Map
	locationsMap *syncmap.Map
	visitsMap    *syncmap.Map
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
	return &DB{
		usersMap:     &syncmap.Map{},
		locationsMap: &syncmap.Map{},
		visitsMap:    &syncmap.Map{},
	}
}

func (db *DB) mapByEntity(entity schema.Entity) *syncmap.Map {
	switch entity {
	case schema.EntityUsers:
		return db.usersMap
	case schema.EntityVisits:
		return db.visitsMap
	case schema.EntityLocations:
		return db.locationsMap
	default:
		return nil
	}
}

// CreateOrUpdate creates entity in db
func (db *DB) CreateOrUpdate(entity schema.IEntity) error {

	m := db.mapByEntity(entity.Entity())

	m.Store(entity.IntID(), entity)

	return nil
}

func (db *DB) Get(entity schema.Entity, id uint32) (interface{}, error) {
	m := db.mapByEntity(entity)

	if result, ok := m.Load(id); ok {
		return result, nil
	}

	return nil, fmt.Errorf("%s with id %d not exists", entity.String(), id)
}

func (db *DB) GetIds(entity schema.Entity, ids []uint32) ([]interface{}, error) {
	var result []interface{}
	for _, id := range ids {
		v, err := db.Get(entity, id)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, nil
}
