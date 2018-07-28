package generators

import (
	//"bytes"
	//"go/format"
	//"go/parser"
	//"go/token"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"

	"github.com/graph-gophers/graphql-go/introspection"
)

var tmpl *template.Template

type Generator interface {
	Generate(schema *introspection.Schema) error
}

/*
	template.ParseFiles(
		"templates/fragments.tmpl",
		"templates/resolver.tmpl",
	)
*/
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
			},
		).ParseGlob("templates/*.tmpl"),
	)
}

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {

	//buf := &bytes.Buffer{}

	//if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		return err
	}

	/*
		fset := token.NewFileSet()
		out, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
		if err != nil {
			return err
		}

		return format.Node(w, fset, out)
	*/
	return nil
}
