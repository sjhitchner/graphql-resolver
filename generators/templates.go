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

	. "github.com/sjhitchner/graphql-resolver/domain"
)

var tmpl *template.Template

type Generator interface {
	Generate(models ...Model) error
}

func init() {
	tmpl = template.Must(
		template.New("").Funcs(
			template.FuncMap{
				"args":    Args,
				"snake":   SnakeCase,
				"camel":   CamelCase,
				"lcamel":  LowerCamelCase,
				"comment": Comment,
				"safe":    Safe,
				"join":    Join,
				"unique":  Unique,
				/*
					"allFields":    AllFields,
					"typeName":     TypeName,
					"isQuery":      IsQuery,
					"isPageInfo":   IsPageInfo,
					"isEdge":       IsEdge,
					"isConnection": IsConnection,
				*/
			},
		).ParseGlob("templates/*.tmpl"),
	)
}

func TemplatePath(path, module, name string) string {
	if err := os.MkdirAll(filepath.Join(path, module), os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(path, module, strcase.SnakeCase(name)+".go")
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
	out, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
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
