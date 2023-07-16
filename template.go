package donut

import (
	"strings"
	"text/template"
)

type templateParams struct {
	Source      string
	Destination string
}

var tmpl = template.New("")

func createTemplateMap(tmap map[string][]string) error {
	for name, args := range tmap {
		if err := createTemplate(name, args...); err != nil {
			return err
		}
	}
	return nil
}

func createTemplate(name string, args ...string) error {
	joined := strings.Join(args, " ")
	if _, err := tmpl.New("diff").Parse(joined); err != nil {
		return err
	}
	return nil
}
