package database

import (
	"testing"

	"github.com/ngalayko/highloadcup/schema"
	. "gopkg.in/check.v1"
)

type DbTestSuite struct{}

var _ = Suite(&DbTestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *DbTestSuite) Test_parseFileName__should_parse_entity_by_file_name(c *C) {

	cases := []struct {
		FileName string
		Entity   schema.Entity
		Error    bool
	}{
		{
			FileName: "users_1.json",
			Entity:   schema.EntityUsers,
			Error:    false,
		},
		{
			FileName: "locations_2.json",
			Entity:   schema.EntityLocations,
			Error:    false,
		},
		{
			FileName: "visits_1.json",
			Entity:   schema.EntityVisits,
			Error:    false,
		},
		{
			FileName: "not_valid.json",
			Error:    true,
		},
	}

	for _, cas := range cases {
		result, err := parseFileName(cas.FileName)
		if cas.Error {
			c.Assert(err, NotNil)
			continue
		}

		c.Assert(err, IsNil)
		c.Assert(result, Equals, cas.Entity)
	}
}
