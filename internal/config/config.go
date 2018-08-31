package config

import ()

/*
// Types
id
string
int
float
email
password
*/

type Config struct {
	Generate []string `yaml:"generate"`

	Models []Model `yaml:"models"`

	GraphQL   *GraphQLModule   `yaml:"graphql,omitempty"`
	Resolvers *ResolversModule `yaml:"resolvers,omitempty"`
	SQL       *SQLModule       `yaml:"sql,omitempty"`
}

type Model struct {
	Name        string  `yaml:"name"`
	Internal    string  `yaml:"internal,omitempty"`
	Description string  `yaml:"description,omitempty"`
	Fields      []Field `yaml:"fields"`
	Deprecated  string  `yaml:"deprecated"`
}

type Field struct {
	Name         string   `yaml:"name"`
	Internal     string   `yaml:"internal,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Type         Type     `yaml:"type"`
	Indexes      []string `yaml:"indexes,omitempty"`
	Deprecated   string   `yaml:"deprecated"`
	Relationship string   `yaml:"relationship"`
}
