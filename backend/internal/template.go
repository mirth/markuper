package internal

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
)

type AnswerField = Jsonable

type TaskAnswer struct {
}

type ValueWithName struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type RadioField struct {
	ID     string          `json:"id"` // fixme keep only Name?
	Type   string          `json:"type"`
	Name   ValueWithName   `json:"name"`
	Labels []ValueWithName `json:"labels"`
}

func NewRadioField(id string, name ValueWithName) RadioField {
	return RadioField{
		ID:   id,
		Name: name,
		Type: "radio",
		Labels: []ValueWithName{
			{Value: "cat", Name: "cat"},
			{Value: "dog", Name: "dog"},
		},
	}
}

func (f RadioField) JSON() ([]byte, error) {
	return json.Marshal(f)
}

type Template struct {
	Task   string       `json:"task"`
	Radios []RadioField `json:"radios"`
	// FieldsOrder []string `json:"fields_order"`
}

type TemplateList struct {
	Templates []Template `json:"templates"`
}

type TemplateService interface {
	ListTemplates() (TemplateList, error)
}

type TemplateServiceImpl struct {
}

var DEFAULT_CLASSIFICATION_TEMPLATE = Template{
	Task: "classification",
	Radios: []RadioField{
		NewRadioField("1", ValueWithName{
			Name:  "class",
			Value: "class",
		}),
	},
}

func (_ *TemplateServiceImpl) ListTemplates() (TemplateList, error) {
	return TemplateList{
		Templates: []Template{DEFAULT_CLASSIFICATION_TEMPLATE},
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
