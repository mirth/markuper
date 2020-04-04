package internal

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AnswerField = Jsonable

type TaskAnswer struct {
}

type ValueWithVizual struct {
	Value  string `json:"value"`
	Vizual string `json:"vizual"`
}

type ClassificationField struct {
	Type   string            `json:"type"`
	Group  string            `json:"group"`
	Labels []ValueWithVizual `json:"labels"`
}

type RadioField = ClassificationField
type CheckboxField = ClassificationField

func NewRadioField(group string) *RadioField {
	return &RadioField{
		Group:  group,
		Type:   "radio",
		Labels: make([]ValueWithVizual, 0),
	}
}

func NewCheckboxField(group string) *CheckboxField {
	return &CheckboxField{
		Group:  group,
		Type:   "checkbox",
		Labels: make([]ValueWithVizual, 0),
	}
}

type Template struct {
	Radios     []RadioField    `json:"radios"`
	Checkboxes []CheckboxField `json:"checkboxes"`

	FieldsOrder []string `json:"fields_order"`
}

func (t *Template) getClassificationFields() []ClassificationField {
	fields := []ClassificationField{}

	for _, f := range t.Radios {
		fields = append(fields, f)
	}

	for _, f := range t.Checkboxes {
		fields = append(fields, f)
	}

	return fields
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
