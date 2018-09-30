package generators

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"text/template"

	//"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

var tmpl *template.Template

type Generator interface {
	Generate(config *config.Config) error
}

func init() {
	tmpl = template.Must(
		template.New("").Funcs(
			template.FuncMap{
				"args":      Args,
				"snake":     SnakeCase,
				"camel":     CamelCase,
				"lcamel":    LowerCamelCase,
				"comment":   Comment,
				"safe":      Safe,
				"join":      Join,
				"unique":    Unique,
				"add":       Add,
				"sub":       Subtract,
				"mul":       Multiply,
				"div":       Divide,
				"now":       Now,
				"gotype":    GoType,
				"gqltype":   GraphQLType,
				"find":      Find,
				"return":    MethodReturn,
				"many2many": IsMany2Many,
				/*
					"allFields":    AllFields,
					"typeName":     TypeName,
					"isQuery":      IsQuery,
					"isPageInfo":   IsPageInfo,
					"isEdge":       IsEdge,
					"isConnection": IsConnection,
				*/
			},
		).ParseGlob("internal/templates/*.tmpl"),
	)
}

func TemplatePath(path, module, name string) string {
	return templatePath(path, module, name, ".go")
}

func SchemaPath(path, module, name string) string {
	return templatePath(path, module, name, ".gql")
}

func templatePath(path, module, name, suffix string) string {
	if err := os.MkdirAll(filepath.Join(path, module), os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(path, module, strcase.SnakeCase(name)+suffix)
}

func GenerateGoFile(filename, template string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	if err := ExecuteTemplate(buf, template, data); err != nil {
		return err
	}

	fset := token.NewFileSet()
	out, err := parser.ParseFile(
		fset,
		"",
		buf,
		parser.AllErrors|parser.ParseComments,
	)
	if err != nil {
		return err
	}

	return format.Node(f, fset, out)
}

func GenerateFile(filename, template string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := ExecuteTemplate(f, template, data); err != nil {
		return err
	}

	return nil
}

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		return err
	}

	return nil
}
