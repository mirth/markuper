package internal

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type ProjectTemplate struct {
	Task string `json:"task"`
}

type ProjectTemplateList struct {
	Templates []ProjectTemplate `json:"templates"`
}

type ProjectTemplateService interface {
	ListProjectTemplates() (ProjectTemplateList, error)
}

type ProjectTemplateServiceImpl struct {
}

func (_ *ProjectTemplateServiceImpl) ListProjectTemplates() (ProjectTemplateList, error) {
	return ProjectTemplateList{
		Templates: []ProjectTemplate{
			{Task: "classification"},
		},
	}, nil
}

func NewProjectTemplateService() ProjectTemplateService {
	return &ProjectTemplateServiceImpl{}
}

func ListProjectTemplatesEndpoint(s ProjectTemplateService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.ListProjectTemplates()
	}
}
