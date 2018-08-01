package domain

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	//"fmt"
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
	Name        string  `yaml:"name"`
	Table       string  `yaml:"table,omitempty"`
	Description string  `yaml:"description,omitempty"`
	Fields      []Field `yaml:"fields"`
	// Module TODO what modules to generate use PROTOBUFS/SQL etc
}

func (t Model) String() string {
	b, err := yaml.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type Field struct {
	Name        string `yaml:"name"`
	Internal    string `yaml:"internal,omitempty"`
	Description string `yaml:"description,omitempty"`
	Type        Type   `yaml:"type"`
}

type Type struct {
	Base  BaseType `yaml:"base"`
	Index string   `yaml:"index,omitempty"`
}

type BaseType int

const (
	ID BaseType = iota
	Integer
	Float
	Boolean
	String
	Time
	DateTime
	Enum

	IDS       = "ID"
	IntegerS  = "integer"
	FloatS    = "float64"
	BooleanS  = "bool"
	StringS   = "string"
	TimeS     = "time"
	DateTimeS = "datetime"
)

func (t BaseType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}

	switch value {
	case IDS:
		t = ID
	case IntegerS:
		t = Integer
	case FloatS:
		t = Float
	case BooleanS:
		t = Boolean
	case StringS:
		t = String
	case TimeS:
		t = Time
	case DateTimeS:
		t = DateTime
	default:
		return errors.Errorf("Invalid Base Type %s", value)
	}

	return nil
}

func (t BaseType) MarshalYAML() (interface{}, error) {
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
		return 0, errors.Errorf("Invalid Base Type %d", t)
	}
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
