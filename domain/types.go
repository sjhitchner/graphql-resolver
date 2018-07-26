package domain

import (
	"fmt"
	//	"github.com/graph-gophers/graphql-go/introspection"
)

type Query struct {
	Name         string
	Args         []Argument
	ResolverName string
	IsList       bool
}

func (t Query) Return() string {
	str := fmt.Sprintf("*%s", t.ResolverName)
	if t.IsList {
		return "[]" + str
	}
	return str
}

type Enum struct {
	Name        string
	Description string
	Values      []Field
}

type Model struct {
	Name        string
	Description string
	Implements  string
	Fields      []Field
}

type Field struct {
	Name        string
	Description string
	Args        []Argument
	Deprecated  string
	ReturnType  Type
}

type Argument struct {
	Name    string
	Type    Type
	Default string
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
