package schema

import (
	"encoding/json"
	"fmt"

	"github.com/ngalayko/highloadcup/helper"
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

	AvgMark float32 `json:"avg"`
}

// Bucket return bucket name
func (l *Location) Bucket() []byte {
	return BucketsMap[EntityLocations]
}

// ByteID is a byte view of id
func (l *Location) ByteID() []byte {
	return helper.Itob(l.ID)
}

// IntID return entity id
func (l *Location) IntID() uint32 {
	return l.ID
}

// Bytes returns bytes view of Location
func (l *Location) Bytes() []byte {
	data, _ := json.Marshal(l)
	return data
}

// Validate validates location fields
func (l *Location) Validate() error {
	switch {
	case len(l.Country) > maxCountryLen:
		return fmt.Errorf("Location.Country longer than %d", maxCountryLen)
	case len(l.City) > maxCityLen:
		return fmt.Errorf("Location.City longer than %d", maxCityLen)
	default:
		return nil
	}
}
