package template

import (
	"bytes"
	"html/template"
)

type systemTemplateManager struct{}

func (m systemTemplateManager) CreateTemplate(ctx string, t string, data interface{}, html bool) (string, error) {
	tp, err := template.New(ctx).Parse(t)
	if err != nil {
		return "", err
	}

	var templateWithData bytes.Buffer
	if html {
		err = tp.ExecuteTemplate(&templateWithData, ctx, data)
	} else {
		err = tp.Execute(&templateWithData, data)
	}

	if err != nil {
		return "", err
	}

	return templateWithData.String(), nil
}
