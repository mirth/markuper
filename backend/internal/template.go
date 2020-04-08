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

type BoundingBoxField struct {
	Type   string            `json:"type"`
	Group  string            `json:"group"`
	Labels []ValueWithVizual `json:"labels"`
}

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

func NewBoundingBoxField(group string) *BoundingBoxField {
	return &BoundingBoxField{
		Type:   "bounding_box",
		Group:  group,
		Labels: make([]ValueWithVizual, 0),
	}
}

type Template struct {
	Radios        []RadioField       `json:"radios"`
	Checkboxes    []CheckboxField    `json:"checkboxes"`
	BoundingBoxes []BoundingBoxField `json:"bounding_boxes"`

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
	Task: "Classification",
	XML: `<content>
    <radio group="animal" value="cat" vizual="Cat" />
    <radio group="animal" value="dog" vizual="Dog" />
</content>
`,
}

var DEFAULT_MULTILABEL_CLASSIFICATION_TEMPLATE = TemplateXML{
	Task: "Multi-label classification",
	XML: `<content>
    <checkbox group="color" value="black" vizual="Black" />
    <checkbox group="color" value="white" vizual="White" />
</content>
`,
}

var DEFAULT_OBJECT_DETECTION_TEMPLATE = TemplateXML{
	Task: "Object detection",
	XML: `<content>
	<bounding_box group="animal" value="cat" vizual="Cat"/>
	<bounding_box group="animal" value="dog" vizual="Dog"/>
</content>
`,
}

func (_ *TemplateServiceImpl) ListTemplates() (TemplateList, error) {
	return TemplateList{
		Templates: []TemplateXML{
			DEFAULT_CLASSIFICATION_TEMPLATE,
			DEFAULT_MULTILABEL_CLASSIFICATION_TEMPLATE,
			DEFAULT_OBJECT_DETECTION_TEMPLATE,
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
