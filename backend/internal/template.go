package internal

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Template struct {
	Task string `json:"task"`
}

type TemplateList struct {
	Templates []Template `json:"templates"`
}

type TemplateService interface {
	ListTemplates() (TemplateList, error)
}

type TemplateServiceImpl struct {
}

func (_ *TemplateServiceImpl) ListTemplates() (TemplateList, error) {
	return TemplateList{
		Templates: []Template{
			{Task: "classification"},
		},
	}, nil
}

func NewTemplateService() TemplateService {
	return &TemplateServiceImpl{}
}

func ListTemplatesEndpoint(s TemplateService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.ListTemplates()
	}
}
