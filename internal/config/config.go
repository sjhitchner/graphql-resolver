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
	//"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
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

	validateConfig(&m)
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

func validateConfig(c *Config) {
	for i := range c.Models {
		validateModel(&c.Models[i])
	}
}

func (t Config) FindModelByName(name string) Model {
	for _, model := range t.Models {
		if model.Name == name {
			return model
		}
	}
	panic("Model Name " + name + " Not Found")
}

/*
func (t Config) FindModelByInternal(internal string) Model {
	for _, model := range t.Models {
		if model.Internal == internal {
			return model
		}
	}
	panic("Model Internal " + internal + " Not Found")
}
*/

// type and import
func (t Config) TypePrimative(base string) (string, string) {
	switch base {
	case "integer":
		return "int64", ""
	case "float":
		return "float64", ""
	case "boolean":
		return "bool", ""
	case "string":
		return "string", ""
	case "timestamp":
		return "time.Time", "time"
	//case "time.Time":
	//	return "time.Time", "time"
	default:
		for _, b := range t.Types {
			if base == b.Name {
				return b.Primative, ""
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
	Name        string   `yaml:"name"`
	Internal    string   `yaml:"internal,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Fields      []Field  `yaml:"fields"`
	Deprecated  string   `yaml:"deprecated,omitempty"`
	Actions     []string `yaml:"actions,omitempty"`
	Queries     []Query  `yaml:"queries,omitempty"`
}

func validateModel(m *Model) {
	if m.Internal == "" {
		m.Internal = m.Name
	}

	for i := range m.Fields {
		validateField(&m.Fields[i])
	}
}

type Field struct {
	Name         string       `yaml:"name"`
	Internal     string       `yaml:"internal,omitempty"`
	Description  string       `yaml:"description,omitempty"`
	Expose       bool         `yaml:"expose,omitempty"`
	Type         string       `yaml:"type"`
	Indexes      []string     `yaml:"indexes,omitempty"`
	Deprecated   string       `yaml:"deprecated,omitempty"`
	Relationship Relationship `yaml:"relationship,omitempty"`
}

func validateField(f *Field) {
	if f.Internal == "" {
		f.Internal = f.Name
	}
}

type Relationship struct {
	To      string `yaml:"to,omitempty"`
	Through string `yaml:"through,omitempty"`
	Field   string `yaml:"field,omitempty"`
	Type    string `yaml:"type,omitempty"`
}

type Query struct {
	Name    string   `yaml:"name,omitempty"`
	Args    []string `yaml:"args,omitempty"`
	Returns string   `yaml:"returns,omitempty"`
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

func (t Model) FindFieldByName(name string) Field {
	for _, field := range t.Fields {
		if field.Name == name {
			return field
		}
	}
	panic("No field " + name + " in model " + t.Name)
}

func (t Model) FindFieldByInternal(internal string) Field {
	for _, field := range t.Fields {
		if field.Internal == internal {
			return field
		}
	}
	panic("No field " + internal + " in model " + t.Internal)
}
