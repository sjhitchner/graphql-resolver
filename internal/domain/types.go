package domain

import (
	"encoding/json"
	"fmt"

	"github.com/stoewer/go-strcase"
)

type Resolver struct {
	Name    string
	Methods []Method
}

type Repo struct {
	Name    string
	Methods []Method
}

type Method struct {
	Type         string
	Name         string
	Args         []Arg
	Relationship *Relationship
	Return       Return
}

type Arg struct {
	Name   string
	Parent string
	Type   string
	Deref  bool
}

type Return struct {
	Type  string
	Multi bool
}

type Model struct {
	Name        string
	Type        string
	Plural      string
	Description string
	Fields      []Field
	Internal    string
	Repo        Repo
	Imports     Imports
	Mutations   []Mutation
}

type Mutation struct {
	Name   string
	Type   string
	Fields []Field
	Key    string
}

type Relationship struct {
	To      string
	Through string
	Field   string
	Type    string
}

type Field struct {
	Name         string
	Type         string
	Primative    string
	Internal     string // Snake Case
	Relationship *Relationship
	ShouldExpose bool
}

type Type struct {
	Name string
	Type string
}

func FindModel(models []Model, name string) Model {
	for _, model := range models {
		if model.Name == name {
			return model
		}
	}
	panic("No model named " + name)
}

/*
func FindRelationship(relationships []Relationship, name string) Relationship {
	for _, relationship := range relationships {
		for _, model := range relationship.Models {
			if model.Name == name {
				return relationship
			}
		}
	}
	return Relationship{}
}
*/

func (t Model) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t Relationship) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func CamelCasef(f string, args ...interface{}) string {
	return strcase.UpperCamelCase(fmt.Sprintf(f, args...))
}

func LowerCamelCasef(f string, args ...interface{}) string {
	return strcase.LowerCamelCase(fmt.Sprintf(f, args...))
}
