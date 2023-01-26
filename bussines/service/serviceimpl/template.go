package serviceimpl

import (
	"bytes"
	"errors"
	"github.com/greenfield0000/microcore/bussines/service"
	"html/template"
)

type HTMLTemplateService struct {
}

func NewHTMLTemplateService() service.TemplateService {
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
