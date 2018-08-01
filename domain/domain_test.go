package domain

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
)

func Test(t *testing.T) { TestingT(t) }

type DomainSuite struct{}

var _ = Suite(&DomainSuite{})

func (s *DomainSuite) Test_Load(c *C) {
	var Data = `
name: "Human"
table: "humans"
description: "A humanoid"
fields: 
  - name: "id"
    internal: "id"
    description: "The ID of the human"
    deprecated: false
    type: 
      base: ID
      index: primary`

	var model Model
	c.Assert(yaml.Unmarshal([]byte(Data), &model), IsNil)

	c.Assert(model.Name, Equals, "Human")
	c.Assert(model.Table, Equals, "humans")
	c.Assert(model.Description, Equals, "A humanoid")
	c.Assert(model.Fields, HasLen, 1)
	field := model.Fields[0]
	c.Assert(field.Name, Equals, "id")
	c.Assert(field.Internal, Equals, "id")
	c.Assert(field.Description, Equals, "The ID of the human")
	//c.Assert(field.Deprecated, Equals, "id")
	c.Assert(field.Type.Base, Equals, ID)
	c.Assert(field.Type.Index, Equals, "primary")
}
