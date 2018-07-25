package generate

import (
	//"bytes"
	//"go/format"
	//"go/parser"
	//"go/token"
	"io"
	"text/template"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(
		template.ParseFiles(
			"templates/fragments.tmpl",
			"templates/resolver.tmpl",
		),
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
