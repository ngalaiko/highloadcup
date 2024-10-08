package schema

import (
	"errors"
)

// IEntity is an interface to any entity in project
type IEntity interface {
	Validate() error
	IntID() uint32
	Entity() Entity
}

// Entity is a enum for all schema components
type Entity uint16

const (
	// EntityUndefined is empty entity
	EntityUndefined = iota
	// EntityUsers is 'users' entity
	EntityUsers
	// EntityLocations is 'locations' entity
	EntityLocations
	// EntityVisits is 'visits' entity
	EntityVisits
)

// String returns the string value of the Entity.
func (t Entity) String() string {
	var enumVal string

	switch t {
	case EntityUndefined:
		enumVal = ""

	case EntityUsers:
		enumVal = "users"

	case EntityLocations:
		enumVal = "locations"

	case EntityVisits:
		enumVal = "visits"
	}

	return enumVal
}

// MarshalText marshals Entity into text.
func (t Entity) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText unmarshall Entity from text.
func (t *Entity) UnmarshalText(text []byte) error {
	switch string(text) {
	case "users":
		*t = EntityUsers

	case "locations":
		*t = EntityLocations

	case "visits":
		*t = EntityVisits

	default:
		return errors.New("invalid Entity")
	}

	return nil
}
