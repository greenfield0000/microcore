package commonservice

import (
	"bytes"
	"errors"
	"html/template"
)

type TemplateService interface {
	Render(path string, data interface{}) (string, error)
}

type HTMLTemplateService struct {
}

func NewHTMLTemplateService() TemplateService {
	return HTMLTemplateService{}
}

func (H HTMLTemplateService) Render(path string, data interface{}) (string, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return "", errors.New("не удалось прочитать шаблон письма")
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", errors.New("не удалось сформировать шаблон письма")
	}
	return buff.String(), nil
}
