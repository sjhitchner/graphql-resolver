package config

import (
	"encoding/json"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

const ConfigPath = "../config/models.yml"

func Test(t *testing.T) {
	TestingT(t)
}

type ConfigSuite struct {
	Config *Config
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpSuite(c *C) {
	config, err := LoadConfigFile(ConfigPath)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	s.Config = config
}

func (s *ConfigSuite) Test_Parsing(c *C) {
	c.Skip("skipping")

	b, err := json.MarshalIndent(s.Config, "", "  ")
	c.Assert(err, IsNil)

	fmt.Println(string(b))

}

func (s *ConfigSuite) Test_Types(c *C) {
	for _, typ := range s.Config.Types {
		primative, impt := s.Config.TypePrimative(typ.Name)
		c.Assert(primative, Equals, typ.Primative)

		if primative == "time.Time" {
			c.Assert(impt, Equals, "time")
		}
	}

	{
		primative, impt := s.Config.TypePrimative("integer")
		c.Assert(primative, Equals, "int64")
		c.Assert(impt, Equals, "")
	}
	{
		primative, impt := s.Config.TypePrimative("float")
		c.Assert(primative, Equals, "float64")
		c.Assert(impt, Equals, "")
	}

	{
		primative, impt := s.Config.TypePrimative("boolean")
		c.Assert(primative, Equals, "bool")
		c.Assert(impt, Equals, "")
	}
	{
		primative, impt := s.Config.TypePrimative("string")
		c.Assert(primative, Equals, "string")
		c.Assert(impt, Equals, "")
	}
	{
		primative, impt := s.Config.TypePrimative("timestamp")
		c.Assert(primative, Equals, "time.Time")
		c.Assert(impt, Equals, "time")
	}

	c.Assert(func() {
		s.Config.TypePrimative("asdf")
	}, PanicMatches, "No type definition for asdf")
}

func (s *ConfigSuite) Test_FindModelByName(c *C) {
	model := s.Config.FindModelByName("user")
	c.Assert(model.Name, Equals, "user")
	c.Assert(model.Internal, Equals, "user")
	c.Assert(model.Fields, HasLen, 5)

	field := model.FindFieldByName("username")
	c.Assert(field.Name, Equals, "username")
	c.Assert(field.Internal, Equals, "username")

	c.Assert(func() {
		model.FindFieldByName("asdf")
	}, PanicMatches, "No field asdf in model user")
}
