package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListProjectTemplates(t *testing.T) {
	svc := NewProjectTemplateService()

	l, err := svc.ListProjectTemplates()
	assert.Nil(t, err)

	{
		assert.ElementsMatch(t, []ProjectTemplate{
			{Task: "classification"},
		}, l.Templates)
	}
}
