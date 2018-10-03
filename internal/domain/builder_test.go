package domain

import (
	"encoding/json"
	"fmt"

	. "gopkg.in/check.v1"
	//"github.com/sjhitchner/graphql-resolver/internal/config"
)

func (s *DomainSuite) Test_Builder(c *C) {
	c.Assert(s.Config, NotNil)

	models, _, err := Parse(s.Config)
	c.Assert(err, IsNil)

	b, _ := json.MarshalIndent(models, "", "  ")
	fmt.Println(string(b))

	/*
		b, _ = json.MarshalIndent(repos, "", "  ")
		fmt.Println(string(b))

		b, _ = json.MarshalIndent(types, "", "  ")
		fmt.Println(string(b))

			c.Assert(models, HasLen, 0)
			c.Assert(repos, HasLen, 0)
			c.Assert(types, HasLen, 0)
	*/
}
