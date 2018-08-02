package domain

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	//"github.com/stoewer/go-strcase"
)

func ReadSchema(r io.Reader) ([]Model, error) {
	dec := yaml.NewDecoder(r)
	dec.SetStrict(true)

	var m []Model
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// TODO proper casing methods
type Model struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description,omitempty"`
	Fields      []Field        `yaml:"fields"`
	Generate    []string       `yaml:"generate"` // Modules to generate
	GraphQL     *GraphQLModule `yaml:"graphql,omitempty"`
	SQL         *SQLModule     `yaml:"sql,omitempty"`
	// Module TODO what modules to generate use PROTOBUFS/SQL etc
	// SQL fragments for active = true etc etc
	// Indexes method
}

func (t Model) String() string {
	b, err := yaml.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type Index struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

func (t Index) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t Model) Indexes() []Index {
	indexes := make(map[string][]Field)
	for _, field := range t.Fields {
		for _, index := range field.Indexes {
			_, found := indexes[index]
			if !found {
				indexes[index] = make([]Field, 0, 10)
			}
			indexes[index] = append(indexes[index], field)
		}
	}

	i := make([]Index, 0, len(indexes))
	for name, fields := range indexes {
		i = append(i, Index{
			Name:   name,
			Fields: fields,
		})
	}
	return i
}

type Field struct {
	Name        string   `yaml:"name"`
	Internal    string   `yaml:"internal,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Type        Type     `yaml:"type"`
	Indexes     []string `yaml:"indexes,omitempty"`
	Deprecated  string   `yaml:"deprecated"`
}

type Type int

const (
	ID Type = iota
	Integer
	Float
	Boolean
	String
	Time
	DateTime
	Enum

	IDS       = "ID"
	IntegerS  = "int"
	FloatS    = "float"
	BooleanS  = "bool"
	StringS   = "string"
	TimeS     = "time"
	DateTimeS = "datetime"
)

func (t *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}

	switch value {
	case IDS:
		*t = ID
	case IntegerS:
		*t = Integer
	case FloatS:
		*t = Float
	case BooleanS:
		*t = Boolean
	case StringS:
		*t = String
	case TimeS:
		*t = Time
	case DateTimeS:
		*t = DateTime
	default:
		return errors.Errorf("Invalid Type %s", value)
	}
	return nil
}

func (t Type) MarshalYAML() (interface{}, error) {
	switch t {
	case ID:
		return IDS, nil
	case Integer:
		return IntegerS, nil
	case Float:
		return FloatS, nil
	case Boolean:
		return BooleanS, nil
	case String:
		return StringS, nil
	case Time:
		return TimeS, nil
	case DateTime:
		return DateTimeS, nil
	default:
		return 0, errors.Errorf("Invalid Type %d", t)
	}
}

func (t Type) String() string {
	switch t {
	case ID:
		return IDS
	case Integer:
		return IntegerS
	case Float:
		return FloatS
	case Boolean:
		return BooleanS
	case String:
		return StringS
	case Time:
		return TimeS
	case DateTime:
		return DateTimeS
	default:
		return fmt.Sprintf("Invalid Type %d", t)
	}
}

const (
	SQL       = "sql"
	GraphQL   = "graphql"
	ProtoBufs = "protobufs"
	Thrift    = "thrift"
)

func (t Model) ShouldGenerate(module string) bool {
	for _, generate := range t.Generate {
		if generate == module {
			return true
		}
	}
	return false
}

type SQLModule struct {
	Package string `yaml:"package,omitempty"`
	Table   string `yaml:"table,omitempty"`
	Dialect string `yaml:"dialect"`
}

type GraphQLModule struct {
	Package string `yaml:"package,omitempty"`
}

type ProtoBufModule struct {
	Package string `yaml:"package,omitempty"`
}

/*
type Enum struct {
	Name        string
	Description string
	Values      []Field
}

type Model struct {
	Name        string // Constructed model name
	Plural      string
	Description string
	Fields      []Field
}

type Field struct {
	Name        string
	Description string
	Deprecated  string
	Type        Type
}

type Argument struct {
	Name        string
	Description string
	Type        Type
	Default     string
}

func (t Argument) DefaultType() string {
	return t.Default
}

func (t Argument) ToGraphQL() string {
	return fmt.Sprintf("%s %s%s", t.Name, t.Type, func() string {
		if t.Default != "" {
			switch t.Type.Base {
			case String:
				return fmt.Sprintf(` = "%s"`, t.Default)
			default:
				return fmt.Sprintf(` = %s`, t.Default)
			}
		}
		return ""
	}())
}
*/
