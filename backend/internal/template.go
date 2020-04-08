package internal

import (
	"backend/pkg/utils"
	"context"
	"fmt"

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
	*ClassificationComponents

	Type  string `json:"type"`
	Group string `json:"group"`
}

type Field interface {
	GetType() string
}

func (f ClassificationField) GetType() string {
	return f.Type
}

func (f BoundingBoxField) GetType() string {
	return f.Type
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
		Type:  "bounding_box",
		Group: group,
		ClassificationComponents: &ClassificationComponents{
			Radios:     make([]*RadioField, 0),
			Checkboxes: make([]*CheckboxField, 0),
		},
	}
}

type ClassificationComponents struct {
	Radios      []*RadioField    `json:"radios"`
	Checkboxes  []*CheckboxField `json:"checkboxes"`
	FieldsOrder []string         `json:"fields_order"`
}

type Template struct {
	*ClassificationComponents
	BoundingBoxes []*BoundingBoxField `json:"bounding_boxes"`
}

func findClField(fields []*ClassificationField, group string) *ClassificationField {
	for _, iterField := range fields {
		if iterField.Group == group {
			return iterField
		}
	}

	return nil
}

func findBBoxField(fields []*BoundingBoxField, group string) *BoundingBoxField {
	for _, iterField := range fields {
		if iterField.Group == group {
			return iterField
		}
	}

	return nil
}

func appendIfNotExists(s *[]string, g string) {
	if len(*s) == 0 || !utils.Contains(*s, g) {
		*s = append(*s, g)
	}
}

func isClField(fieldName string) bool {
	if fieldName == "radio" {
		return true
	}

	if fieldName == "checkbox" {
		return true
	}

	return false
}

func (t *ClassificationComponents) CreateOrUpdateClFieldFor(n Node) {
	g := getGroup(n)
	var f *ClassificationField

	switch n.XMLName.Local {
	case "radio":
		f = findClField(t.Radios, g)
		if f == nil {
			f = NewRadioField(g)
			t.Radios = append(t.Radios, f)
		}
	case "checkbox":
		f = findClField(t.Checkboxes, g)
		if f == nil {
			f = NewCheckboxField(g)
			t.Checkboxes = append(t.Checkboxes, f)
		}
	}

	f.Labels = append(f.Labels, ValueWithVizual{
		Vizual: getVizual(n),
		Value:  getValue(n),
	})

	appendIfNotExists(&t.FieldsOrder, g)
}

func (t *Template) CreateOrUpdateBBoxFieldFor(n Node) error {
	g := getGroup(n)
	box := findBBoxField(t.BoundingBoxes, g)
	if box == nil {
		box = NewBoundingBoxField(g)
	}

	for _, iterNode := range n.Nodes {
		if !isClField(iterNode.XMLName.Local) {
			return NewBusinessError(
				fmt.Sprintf("Unsupported element [%s] in bounding_box field", n.XMLName.Local),
			)
		}
		box.CreateOrUpdateClFieldFor(iterNode)
	}

	appendIfNotExists(&t.FieldsOrder, g)

	return nil
}

func (t *Template) getClassificationFields() []*ClassificationField {
	fields := []*ClassificationField{}

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
