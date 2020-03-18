package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLToTemplateRadios(t *testing.T) {
	data := `
	<content>
		<radio group="animal" value="cat" vizual="Cat" />
		<radio group="animal" value="dog" vizual="Dog" />
	</content>
	`
	tt, _ := XMLToTemplate(data)

	assert.Equal(t, []RadioField{{
		Type:  "radio",
		Group: "animal",
		Labels: []ValueWithVizual{
			{Vizual: "Cat", Value: "cat"},
			{Vizual: "Dog", Value: "dog"},
		},
	}}, tt.Radios)
}
