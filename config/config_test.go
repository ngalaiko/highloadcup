package config

import (
	"testing"

	. "gopkg.in/check.v1"
)

type ConfigTestSuite struct {
	ConfigPath string
}

var _ = Suite(&ConfigTestSuite{
	ConfigPath: "../config.yaml",
})

func Test(t *testing.T) { TestingT(t) }

func (s *ConfigTestSuite) Test_Parse__should_parse_config(c *C) {
	config := NewConfig(s.ConfigPath)

	c.Assert(config.DataPath, Equals, "./data")
	c.Assert(config.DbPath, Equals, "/tmp/highloadcup.db")
}
