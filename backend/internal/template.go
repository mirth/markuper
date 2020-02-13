package internal

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
)

type AnswerField = Jsonable

type TaskAnswer struct {
}

type StringWithName struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type RadioField struct {
	ID     string           `json:"id"`
	Type   string           `json:"type"`
	Labels []StringWithName `json:"labels"`
}

func NewRadioField(id string) RadioField {
	return RadioField{
		ID:     id,
		Type:   "radio",
		Labels: make([]StringWithName, 1), // fixme wtf 1 expected: internal.Template{Task:"classification", Radios:[]internal.RadioField{internal.RadioField{ID:"1", Type:"radio", Labels:[]internal.StringWithName{}}}} actual  : internal.Template{Task:"classification", Radios:[]internal.RadioField{internal.RadioField{ID:"1", Type:"radio", Labels:[]internal.StringWithName(nil)}}
	}
}

type RadioAnswer struct {
	FieldType string `json:"type"`
	Value     string `json:"label"`
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
		NewRadioField("1"),
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
