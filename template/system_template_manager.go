package template

import (
	"bytes"
	ht "html/template"
	tt "text/template"
)

type systemTemplateManager struct{}

func (m systemTemplateManager) CreateTemplate(ctx string, t string, data interface{}, html bool) (string, error) {
	var templateWithData bytes.Buffer
	if !html {
		ttp, err := tt.New(ctx).Parse(t)
		if err != nil {
			return "", err
		}
		err = ttp.ExecuteTemplate(&templateWithData, ctx, data)
		if err != nil {
			return "", err
		}
		return templateWithData.String(), nil
	}

	tp, err := ht.New(ctx).Parse(t)
	if err != nil {
		return "", err
	}
	err = tp.ExecuteTemplate(&templateWithData, ctx, data)
	if err != nil {
		return "", err
	}

	return templateWithData.String(), nil
}
