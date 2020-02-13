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
		assert.ElementsMatch(t, []Template{
			DEFAULT_CLASSIFICATION_TEMPLATE,
		}, l.Templates)
	}
}
