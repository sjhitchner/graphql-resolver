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
	Name   string
	Args   []Arg
	Return Return
}

type Arg struct {
	Name string
	Type string
}

type Return struct {
	Type  string
	Multi bool
}

type Model struct {
	Name        string
	Description string
	Fields      []Field
	Repo        Repo
	//	Methods     []Method
	Imports []string
}

type Relationship struct {
	Name        string
	Description string
	Fields      []Field
	RepoName    string
	Methods     []Method
	Through     string
	Models      []Model
}

type Field struct {
	Name         string
	Type         string
	Primative    string
	Internal     string // Snake Case
	Relationship string
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
