package templatekit

import (
	"bytes"
	"text/template"
)

func process(t *template.Template, vars any) string {
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

// ProcessString vars一般是map或struct; str包含{{.xx}}
func ProcessString(str string, vars any) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
