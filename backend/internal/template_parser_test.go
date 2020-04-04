package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingAttribute(t *testing.T) {
	{
		data := `
		<content>
			<radio group="animal" value="cat" vizual="Cat" />
			<radio group="animal" value="dog" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [radio] missing the attribute [vizual]", err.Error())
	}

	{
		data := `
		<content>
			<radio group="animal" />
			<radio group="animal" value="dog" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [radio] missing the attribute [value]", err.Error())
	}

	{
		data := `
		<content>
			<checkbox />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [checkbox] missing the attribute [group]", err.Error())
	}

}

func TestUnsupportedElement(t *testing.T) {
	data := `
	<content>
		<radio group="animal" value="cat" vizual="Cat" />
		<radio group="animal" value="dog" vizual="Dog" />

		<kek group="color" value="black" vizual="Black" />
		<kek group="color" value="white" vizual="White" />
		<kek group="color" value="pink" vizual="Pink" />
	</content>
	`
	_, err := XMLToTemplate(data)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported element [kek]", err.Error())
}

func TestXMLToTemplateRadios(t *testing.T) {
	{
		data := `
		<content>
			<radio group="animal" value="cat" vizual="Cat" />
			<radio group="animal" value="dog" vizual="Dog" />

			<checkbox group="color" value="black" vizual="Black" />
			<checkbox group="color" value="white" vizual="White" />
			<checkbox group="color" value="pink" vizual="Pink" />
		</content>
		`
		tt, err := XMLToTemplate(data)
		assert.Nil(t, err)

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

	{
		data := `
		<content>
			<radio group="animal" value="cat" vizual="Cat" />
			<radio group="animal" value="dog" vizual="Dog" />

			<checkbox group="animal" value="black" vizual="Black" />
			<checkbox group="animal" value="white" vizual="White" />
			<checkbox group="color" value="pink" vizual="Pink" />
		</content>
		`
		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate groups: animal", err.Error())
	}

	{
		data := `
		<content>
			<radio group="animal" value="cat" vizual="Cat" />
			<radio group="animal" value="cat" vizual="Dog" />

			<checkbox group="color" value="black" vizual="Black" />
			<checkbox group="color" value="white" vizual="White" />
			<checkbox group="color" value="white" vizual="white" />
		</content>
		`
		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate labels: group [animal] labels [cat], group [color] labels [white]", err.Error())
	}
}

// fixme test empty labels
func TestDuplicatedLabels(t *testing.T) {
	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:   "radio",
				Group:  "animal",
				Labels: []ValueWithVizual{{"dog", "Dog"}},
			}},
			Checkboxes: []CheckboxField{{
				Type:   "checkbox",
				Group:  "color",
				Labels: []ValueWithVizual{{"white", "White"}},
			}},
		}

		dups := duplicatedLabels(tmplt)
		assert.Empty(t, dups)
	}

	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:   "radio",
				Group:  "animal",
				Labels: []ValueWithVizual{{"dog", "Dog"}, {"dog", "Dogh"}},
			}},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"animal": {"dog"},
		}, dups)
	}

	{
		tmplt := Template{
			Checkboxes: []CheckboxField{{
				Type:  "checkbox",
				Group: "color",
				Labels: []ValueWithVizual{
					{"white", "White"},
					{"white", "White"},
					{"black", "Black"},
					{"black", "Black"},
				},
			}},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"color": {"white", "black"},
		}, dups)
	}

	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:   "radio",
				Group:  "animal",
				Labels: []ValueWithVizual{{"dog", "Dog"}, {"dog", "Dogh"}},
			}},
			Checkboxes: []CheckboxField{{
				Type:  "checkbox",
				Group: "color",
				Labels: []ValueWithVizual{
					{"white", "White"},
					{"white", "White"},
					{"black", "Black"},
					{"black", "Black"},
				},
			}},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"animal": {"dog"},
			"color":  {"white", "black"},
		}, dups)
	}
}

func TestDuplicatedGroups(t *testing.T) {
	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:  "radio",
				Group: "animal",
			}},
			Checkboxes: []CheckboxField{{
				Type:  "checkbox",
				Group: "color",
			}},
		}

		dups := duplicatedGroups(tmplt)
		assert.Empty(t, dups)
	}

	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:  "radio",
				Group: "animal",
			}},
			Checkboxes: []CheckboxField{{
				Type:  "checkbox",
				Group: "animal",
			}},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{"animal"}, dups)
	}

	{
		tmplt := Template{
			Radios: []RadioField{{
				Type:  "radio",
				Group: "animal",
			}},
			Checkboxes: []CheckboxField{
				{
					Type:  "checkbox",
					Group: "color",
				},
				{
					Type:  "checkbox",
					Group: "animal",
				}},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{"animal"}, dups)
	}

	{
		tmplt := Template{
			Radios:     []RadioField{},
			Checkboxes: []CheckboxField{},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{}, dups)
	}

	{
		tmplt := Template{
			Radios: []RadioField{},
			Checkboxes: []CheckboxField{{
				Type:  "checkbox",
				Group: "color",
			}},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{}, dups)
	}
}
