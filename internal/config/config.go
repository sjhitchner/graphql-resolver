package config

/*
// Types
id
string
int
float
email
password
*/

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v2"
)

func LoadConfigFile(str string) (*Config, error) {
	f, err := os.Open(str)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to read config '%s'", str)
	}
	defer f.Close()

	return LoadConfig(f)
}

func LoadConfig(r io.Reader) (*Config, error) {
	dec := yaml.NewDecoder(r)
	dec.SetStrict(true)

	var m Config
	if err := dec.Decode(&m); err != nil {
		return nil, errors.Wrapf(err, "Unable to load config")
	}

	return &m, nil
}

type Config struct {
	Generate []string `yaml:"generate"`
	Models   []Model  `yaml:"models"`
	Types    []Type   `yaml:"custom_types"`

	GraphQL   *GraphQLModule   `yaml:"graphql,omitempty"`
	Resolvers *ResolversModule `yaml:"resolvers,omitempty"`
	SQL       *SQLModule       `yaml:"sql,omitempty"`
}

func (t Config) TypeMapping(base string) string {
	switch base {
	case "integer":
		return "int64"
	case "float":
		return "float64"
	case "boolean":
		return "bool"
	case "string":
		return "string"
	default:
		return strcase.UpperCamelCase(base)
	}
}

func (t Config) TypePrimative(base string) string {
	switch base {
	case "integer":
		return "int64"
	case "float":
		return "float64"
	case "boolean":
		return "bool"
	case "string":
		return "string"
	default:
		for _, b := range t.Types {
			if base == b.Name {
				return b.Primative
			}
		}
	}
	panic("No type definition for " + base)
}

func (t Config) ShouldGenerate(str string) bool {
	for _, g := range t.Generate {
		if g == str {
			return true
		}
	}
	return false
}

type Model struct {
	Name        string  `yaml:"name"`
	Internal    string  `yaml:"internal,omitempty"`
	Description string  `yaml:"description,omitempty"`
	Fields      []Field `yaml:"fields"`
	Deprecated  string  `yaml:"deprecated,omitempty"`
}

type Field struct {
	Name         string   `yaml:"name"`
	Internal     string   `yaml:"internal,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Type         string   `yaml:"type"`
	Indexes      []string `yaml:"indexes,omitempty"`
	Deprecated   string   `yaml:"deprecated,omitempty"`
	Relationship string   `yaml:"relationship,omitempty"`
}

type Type struct {
	Name      string `yaml:"name"`
	Primative string `yaml:"primative"`
}

func (t Config) String() string {
	buf := &bytes.Buffer{}
	if err := yaml.NewEncoder(buf).Encode(t); err != nil {
		return err.Error()
	}
	return buf.String()
}

type GraphQLModule struct {
	Package string `yaml:"package"`
}

type ResolversModule struct {
	Package string `yaml:"package"`
}

type SQLModule struct {
	Package string `yaml:"package"`
	Dialect string `yaml:"dialect"`
}
