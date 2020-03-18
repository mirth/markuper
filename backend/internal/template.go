package internal

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
)

type AnswerField = Jsonable

type TaskAnswer struct {
}

type ValueWithVizual struct {
	Value  string `json:"value"`
	Vizual string `json:"vizual"`
}

type RadioField struct {
	Type   string            `json:"type"`
	Group  string            `json:"group"`
	Labels []ValueWithVizual `json:"labels"`
}

func NewRadioField(group string) RadioField {
	return RadioField{
		Group:  group,
		Type:   "radio",
		Labels: make([]ValueWithVizual, 0),
	}
}

func (f RadioField) JSON() ([]byte, error) {
	return json.Marshal(f)
}

type Template struct {
	Radios []RadioField `json:"radios"`
	// FieldsOrder []string `json:"fields_order"`
}

type TemplateXML struct {
	Task string `json:"task"`
	XML  string `json:"xml"`
}

type TemplateList struct {
	Templates []TemplateXML `json:"templates"`
}

type TemplateService interface {
	ListTemplates() (TemplateList, error)
}

type TemplateServiceImpl struct {
}

var DEFAULT_CLASSIFICATION_TEMPLATE = TemplateXML{
	Task: "classification",
	XML: `<content>
	<radio group="animal" value="cat" vizual="Cat" />
	<radio group="animal" value="dog" vizual="Dog" />
</content>
`,
}

func (_ *TemplateServiceImpl) ListTemplates() (TemplateList, error) {
	return TemplateList{
		Templates: []TemplateXML{DEFAULT_CLASSIFICATION_TEMPLATE},
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
