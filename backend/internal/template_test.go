package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTemplates(t *testing.T) {
	svc := NewTemplateService()

	l, err := svc.ListTemplates()
	assert.Nil(t, err)

	{
		assert.ElementsMatch(t, []TemplateXML{
			DEFAULT_CLASSIFICATION_TEMPLATE,
			DEFAULT_MULTILABEL_CLASSIFICATION_TEMPLATE,
			DEFAULT_OBJECT_DETECTION_TEMPLATE,
		}, l.Templates)
	}
}
