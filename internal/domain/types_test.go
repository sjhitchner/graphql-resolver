package domain

import (
	//"encoding/json"
	//"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

const ConfigPath = "../config/models.yml"

func Test(t *testing.T) {
	TestingT(t)
}

type DomainSuite struct {
	Config *config.Config
}

var _ = Suite(&DomainSuite{})

func (s *DomainSuite) SetUpSuite(c *C) {
	config, err := config.LoadConfigFile(ConfigPath)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	s.Config = config
}

func (s *DomainSuite) Test_Parse(c *C) {
	c.Assert(s.Config, NotNil)

	/*
		models, repos, types := ProcessConfig(s.Config)

		b, _ := json.MarshalIndent(models, "", "  ")
		fmt.Println(string(b))

		b, _ = json.MarshalIndent(repos, "", "  ")
		fmt.Println(string(b))

		b, _ = json.MarshalIndent(types, "", "  ")
		fmt.Println(string(b))
	*/
	/*
		c.Assert(models, HasLen, 0)
		c.Assert(repos, HasLen, 0)
		c.Assert(types, HasLen, 0)
	*/
}
