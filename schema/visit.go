package schema

import (
	"fmt"
)

const (
	maxVisitedAt = 1420070400
	minVisitedAt = 946684800
)

// Visits is a view for array of visits
type Visits struct {
	Visits []*Visit `json:"visits"`
}

// Visit is a visit view from db
type Visit struct {
	ID         uint32 `json:"id"`
	LocationID uint32 `json:"location"`
	UserID     uint32 `json:"user"`
	VisitedAt  int64  `json:"visited_at"`
	Mark       uint8  `json:"mark"`
}

// Entity return entity
func (v *Visit) Entity() Entity {
	return EntityVisits
}

// IntID return entity id
func (v *Visit) IntID() uint32 {
	return v.ID
}

// Validate validates Visit view
func (v *Visit) Validate() error {
	switch {
	case v.ID == 0:
		return fmt.Errorf("Visit.ID is null")
	case v.VisitedAt > maxVisitedAt:
		return fmt.Errorf("Visit.VisitedAt more than %d", maxVisitedAt)
	case v.VisitedAt < minVisitedAt:
		return fmt.Errorf("Visit.VisitedAt less than %d", minVisitedAt)
	default:
		return nil
	}
}
