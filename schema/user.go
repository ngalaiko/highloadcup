package schema

import (
	"fmt"
)

const (
	maxEmailLen     = 100
	maxFirstNameLen = 50
	maxLastNameLen  = 50
	maxBirthYear    = 915148800
	minBirthYear    = -1262304000
)

// Users is a view of array of users
type Users struct {
	Users []*User `json:"users"`
}

// User is a user view from db
type User struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    Gender `json:"gender"`
	BirthDate int64  `json:"birth_date"`

	VisitIDs []uint32 `json:"-"`
}

// Entity return entity
func (u *User) Entity() Entity {
	return EntityUsers
}

// IntID return entity id
func (u *User) IntID() uint32 {
	return u.ID
}

// Validate validates user view
func (u *User) Validate() error {
	switch {
	case len(u.Email) > maxEmailLen:
		return fmt.Errorf("User.Email longer than %d", maxEmailLen)
	case len(u.FirstName) > maxFirstNameLen:
		return fmt.Errorf("User.FirstName longer than %d", maxFirstNameLen)
	case len(u.LastName) > maxLastNameLen:
		return fmt.Errorf("User.LastName longer than %d", maxLastNameLen)
	case u.BirthDate > maxBirthYear:
		return fmt.Errorf("User.BirthDate more than %d", maxLastNameLen)
	case u.BirthDate < minBirthYear:
		return fmt.Errorf("User.BirthDate less than %d", minBirthYear)
	default:
		return nil
	}
}
