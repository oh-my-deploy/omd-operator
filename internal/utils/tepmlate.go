package utils

import (
	"bytes"
	"errors"
	"html/template"
)

func ParsingPreviewTemplate(format string, data map[string]string) (string, error) {
	tmpl, err := template.New("preview").Parse(format)
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing preview template"))
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", errors.Join(err, errors.New("error excuting preview template"))
	}
	return tpl.String(), nil
}
