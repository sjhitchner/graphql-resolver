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
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func LoadConfigAtPath(str string) (*Config, error) {
	files, err := filepath.Glob(filepath.Join(str, "*.yml"))
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find *.yml config file in '%s'", str)
	}

	buf := &bytes.Buffer{}
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read config '%s'", f)
		}
		buf.Write(b)
	}

	return LoadConfig(buf)
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

const (
	ID           = "id"
	One2One      = "one2one"
	One2Many     = "one2many"
	Many2Many    = "many2many"
	UniqueIndex  = "unique"
	PrimaryIndex = "primary"
	MultiIndex   = "index"

	Integer   = "integer"
	Float     = "float"
	String    = "string"
	Boolean   = "boolean"
	Timestamp = "timestamp"

	Link   = "link"
	Domain = "domain"

	Create = "create"
	Update = "update"
	Delete = "delete"
)

type Config struct {
	Generate []string `yaml:"generate"`
	Models   []Model  `yaml:"models"`
	Types    []Type   `yaml:"custom_types"`

	GraphQL   *GraphQLModule   `yaml:"graphql,omitempty"`
	Resolvers *ResolversModule `yaml:"resolvers,omitempty"`
	SQL       *SQLModule       `yaml:"sql,omitempty"`

	BaseImport string
}

func validateConfig(c *Config) {
	for i := range c.Models {
		validateModel(&c.Models[i])
	}
}

func (t Config) FindModelByName(name string) *Model {
	for _, model := range t.Models {
		if model.Name == name {
			return &model
		}
	}
	panic("Model Name '" + name + "' Not Found")
}

// type and import
func (t Config) Primative(base string) (string, string) {
	switch base {
	case Integer:
		return "int64", ""
	case Float:
		return "float64", ""
	case Boolean:
		return "bool", ""
	case String:
		return "string", ""
	case Timestamp:
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
	panic("No type definition for '" + base + "'")
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
	Name        string     `yaml:"name"`
	Plural      string     `yaml:"plural"`
	Type        string     `yaml:"type"`
	Internal    string     `yaml:"internal,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Fields      []Field    `yaml:"fields"`
	Deprecated  string     `yaml:"deprecated,omitempty"`
	Mutations   []Mutation `yaml:"mutations,omitempty"`
}

func validateModel(m *Model) {
	if m.Plural == "" {
		m.Plural = m.Name + "s"
	}

	if m.Internal == "" {
		m.Internal = m.Name
	}

	if m.Type == "" {
		m.Type = Domain
	}

	for i := range m.Fields {
		validateField(&m.Fields[i])
	}

	for i := range m.Mutations {
		validateMutation(&m.Mutations[i])
	}
}

type Field struct {
	Name         string        `yaml:"name"`
	Internal     string        `yaml:"internal,omitempty"`
	Description  string        `yaml:"description,omitempty"`
	Expose       *bool         `yaml:"expose,omitempty"`
	Type         string        `yaml:"type"`
	Indexes      []string      `yaml:"indexes,omitempty"`
	Deprecated   string        `yaml:"deprecated,omitempty"`
	Relationship *Relationship `yaml:"relationship,omitempty"`
}

func validateField(f *Field) {
	if f.Type == "" {
		f.Type = ID
	}

	if f.Internal == "" {
		f.Internal = f.Name
		if f.Relationship != nil {
			if f.Internal != "id" && !strings.HasSuffix(f.Internal, "_id") {
				f.Internal = f.Name + "_id"
			}
		}
	}

	if f.Expose == nil {
		f.Expose = func() *bool {
			b := true
			return &b
		}()
	}

	if f.Relationship != nil {
		switch f.Relationship.Type {
		case Many2Many:
			f.Relationship.Type = Many2Many

		case One2Many:
			f.Relationship.Type = One2Many

		case One2One:
			f.Relationship.Type = One2One

		default:
			panic("Invalid '" + f.Relationship.Type + "' Relationship Type")
		}
	}
}

type Mutation struct {
	Name   string   `yaml:"name"`
	Type   string   `yaml:"type"`
	Fields []string `yaml:"fields"`
	Key    string   `yaml:"key"`
	//OutputFields []string `yaml:"outputs"`
}

func validateMutation(m *Mutation) {
	if m.Key == "" {
		m.Key = "id"
	}

	if !strings.Contains("insert|update|delete", m.Type) {
		m.Type = "insert"
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

type Index struct {
	Name   string
	Type   string
	Fields []string
}

func (t Index) NameWithIds() string {
	s := strings.Split(t.Name, "_")
	return strings.Join(s, "_")
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
	panic("No field named " + name + " in model " + t.Name)
}

func (t Model) FindFieldByInternal(internal string) Field {
	for _, field := range t.Fields {
		if field.Internal == internal {
			return field
		}
	}
	panic("No field internal " + internal + " in model " + t.Name)
}

func (t Model) Indexes() []Index {
	indexMap := make(map[string][]string)

	for _, field := range t.Fields {
		for _, index := range field.Indexes {
			_, found := indexMap[index]
			if !found {
				indexMap[index] = make([]string, 0, 10)
			}
			indexMap[index] = append(indexMap[index], field.Internal)
		}
	}

	indexes := make([]Index, 0, len(indexMap))
	for n, fields := range indexMap {
		name, typ := ParseIndex(n)
		typ = Validate(typ, MultiIndex, UniqueIndex, PrimaryIndex)
		indexes = append(indexes, Index{
			Name:   name,
			Type:   typ,
			Fields: fields,
		})
	}
	return indexes
}

func ParseIndex(index string) (string, string) {
	s := strings.Split(index, "_")

	if len(s) == 1 {
		return ID, PrimaryIndex
	}

	return strings.Join(s[0:len(s)-1], "_"), s[len(s)-1]
}

func Validate(input string, valid ...string) string {
	for _, v := range valid {
		if v == input {
			return input
		}
	}
	panic(fmt.Sprintf("Invalid Input Not Valid '%s' %v", input, valid))
}
