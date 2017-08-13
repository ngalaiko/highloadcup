package schema

import (
	"errors"
)

// Gender is emun type for user.gender
type Gender uint16

const (
	// GenderMale is a `m` gender
	GenderMale = iota
	// GenderFemale is a `f` gender
	GenderFemale
)

// String returns the string value of the Gender.
func (t Gender) String() string {
	var enumVal string

	switch t {
	case GenderMale:
		enumVal = "m"

	case GenderFemale:
		enumVal = "f"

	}

	return enumVal
}

// MarshalText marshals Gender into text.
func (t Gender) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText unmarshall Gender from text.
func (t *Gender) UnmarshalText(text []byte) error {
	switch string(text) {
	case "m":
		*t = GenderMale

	case "f":
		*t = GenderFemale

	default:
		return errors.New("invalid Gender")
	}

	return nil
}
