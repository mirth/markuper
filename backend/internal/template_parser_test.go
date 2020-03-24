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

		<checkbox group="color" value="black" vizual="Black" />
		<checkbox group="color" value="white" vizual="White" />
		<checkbox group="color" value="pink" vizual="Pink" />
	</content>
	`
	tt, _ := XMLToTemplate(data)

	assert.Equal(t, []string{"animal", "color"}, tt.FieldsOrder)
	assert.ElementsMatch(t, []RadioField{{
		Type:  "radio",
		Group: "animal",
		Labels: []ValueWithVizual{
			{Vizual: "Cat", Value: "cat"},
			{Vizual: "Dog", Value: "dog"},
		},
	}}, tt.Radios)
	assert.ElementsMatch(t, []CheckboxField{{
		Type:  "checkbox",
		Group: "color",
		Labels: []ValueWithVizual{
			{Vizual: "Black", Value: "black"},
			{Vizual: "White", Value: "white"},
			{Vizual: "Pink", Value: "pink"},
		},
	}}, tt.Checkboxes)
}
