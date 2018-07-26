package generate

import (
	//"bytes"
	//"go/format"
	//"go/parser"
	//"go/token"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
)

var tmpl *template.Template

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
				"args": func(values ...interface{}) (map[string]interface{}, error) {
					if len(values)%2 != 0 {
						return nil, errors.New("invalid dict call")
					}
					dict := make(map[string]interface{}, len(values)/2)
					for i := 0; i < len(values); i += 2 {
						key, ok := values[i].(string)
						if !ok {
							return nil, errors.New("dict keys must be strings")
						}
						dict[key] = values[i+1]
					}
					return dict, nil
				},
				"snake": func(values ...interface{}) (string, error) {
					s, ok := values[0].(string)
					if !ok {
						return "", errors.Errorf("Invalud argument '%s'", values[0])
					}
					return strcase.SnakeCase(s), nil
				},
				"camel": func(values ...interface{}) (string, error) {
					s, ok := values[0].(string)
					if !ok {
						return "", errors.Errorf("Invalud argument '%s'", values[0])
					}
					return strcase.UpperCamelCase(strcase.SnakeCase(s)), nil
				},
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
