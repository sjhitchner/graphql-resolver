package domain

import (
	"fmt"
	"strings"
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
description: "A humanoid"
fields: 
  - name: "id"
    internal: "id"
    description: "The ID of the human"
    deprecated: false
    type: ID
    indexes: 
      - primary
`

	var model Model
	c.Assert(yaml.Unmarshal([]byte(Data), &model), IsNil)

	c.Assert(model.Name, Equals, "Human")
	c.Assert(model.Description, Equals, "A humanoid")
	c.Assert(model.Fields, HasLen, 1)
	field := model.Fields[0]
	c.Assert(field.Name, Equals, "id")
	c.Assert(field.Internal, Equals, "id")
	c.Assert(field.Description, Equals, "The ID of the human")
	//c.Assert(field.Deprecated, Equals, "id")
	c.Assert(field.Type, Equals, ID)
	c.Assert(field.Indexes, HasLen, 1)
	c.Assert(field.Indexes[0], Equals, "primary")
}

const Models = `---
- name: "Human"
  description: "A humanoid creature from the Star Wars universe"
  generate: 
    - sql
    - graphql
    - protobuf
  graphql:
    package: "resolvers"
  sql:
    table: "humans"
    dialect: "postgres"
  fields: 
    - name: "id"
      internal: "id"
      description: "The ID of the human"
      deprecated: false
      type: ID
      indexes: 
        - primary
    - name: "name"
      internal: "name"
      description: "What this human calls themselves"
      deprecated: false
      type: string
      indexes:
        - name_index
        - name_height_index
    - name: "height"
      internal: "height"
      description: "Height in the preferred unit, default is meters"
      deprecated: false
      type: string
      indexes:
        - name_height_index
    - name: "mass"
      internal: "mass"
      description: "Mass in kilograms, or null if unknown"
      deprecated: false
      type: string
`

func (s *DomainSuite) Test_Indexes(c *C) {
	reader := strings.NewReader(Models)

	models, err := ReadSchema(reader)
	c.Assert(err, IsNil)

	c.Assert(models, HasLen, 1)
	model := models[0]

	c.Assert(model.Name, Equals, "Human")

	indexes := model.Indexes()

	fmt.Println(indexes)

	c.Assert(indexes, HasLen, 3)
	c.Assert(indexes[0].Name, Equals, "primary")
	c.Assert(indexes[1].Name, Equals, "name_index")
	c.Assert(indexes[2].Name, Equals, "name_height_index")
}
