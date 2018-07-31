package domain

import (
	"fmt"
	//	"github.com/graph-gophers/graphql-go/introspection"
)

/*
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
*/

type Type struct {
	Base    BaseType
	Values  []string
	Indexed bool
}

type BaseType int

const (
	Null BaseType = iota
	ID
	Integer
	Float
	Boolean
	String
	Time
	DateTime
	List
	Object
	Interface
	ENUM
	Union
	InputObject
	NonNull
	UserType
)

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
