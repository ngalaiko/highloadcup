package schema

import (
	"fmt"
)

const (
	maxCountryLen = 50
	maxCityLen    = 50
)

// Locations is a view of array of locations
type Locations struct {
	Locations []*Location `json:"locations"`
}

// Location is a location view from db
type Location struct {
	ID       uint32 `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint32 `json:"distance"`

	VisitIDs []uint32 `json:"-"`
}

// Entity return entity
func (l *Location) Entity() Entity {
	return EntityLocations
}

// IntID return entity id
func (l *Location) IntID() uint32 {
	return l.ID
}

// Validate validates location fields
func (l *Location) Validate() error {
	switch {
	case l.ID == 0:
		return fmt.Errorf("User.ID is null")
	case len(l.Country) > maxCountryLen:
		return fmt.Errorf("Location.Country longer than %d", maxCountryLen)
	case len(l.City) > maxCityLen:
		return fmt.Errorf("Location.City longer than %d", maxCityLen)
	default:
		return nil
	}
}
