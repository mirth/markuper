package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingAttribute(t *testing.T) {
	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="Cat" />
			<radio group="animal" value="dog" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [radio] missing the attribute [visual]", err.Error())
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

func TestEmptyAttribute(t *testing.T) {
	{
		data := `
		<content>
			<radio group="animal" value="" visual="Cat" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [radio] has an empty attribute [value]", err.Error())
	}

	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [radio] has an empty attribute [visual]", err.Error())
	}

	{
		data := `
		<content>
			<checkbox group="" value="cat" visual="Cat" />
		</content>
		`

		_, err := XMLToTemplate(data)
		assert.NotNil(t, err)
		assert.Equal(t, "Element [checkbox] has an empty attribute [group]", err.Error())
	}
}

func TestUnsupportedElement(t *testing.T) {
	data := `
	<content>
		<radio group="animal" value="cat" visual="Cat" />
		<radio group="animal" value="dog" visual="Dog" />

		<kek group="color" value="black" visual="Black" />
		<kek group="color" value="white" visual="White" />
		<kek group="color" value="pink" visual="Pink" />
	</content>
	`
	_, err := XMLToTemplate(data)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported element [kek]", err.Error())
}

func newValue(value, name string) ValueWithvisual {
	color := ""
	return NewValueWithvisualWithColor(value, name, color)
}

// fixme duplicatedGroups for bbox
func TestXMLToFields(t *testing.T) {
	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="Cat" />
			<radio group="animal" value="dog" visual="Dog" />

			<checkbox group="color" value="black" visual="Black" />
			<checkbox group="color" value="white" visual="White" />
			<checkbox group="color" value="pink" visual="Pink" />

			<bounding_box group="box">
				<radio group="box_kind" value="kitty" visual="Kitty" />
				<radio group="box_kind" value="doge" visual="Doge" />
			</bounding_box>
		</content>
		`
		tt, err := XMLToTemplate(data)
		assert.Nil(t, err)

		assert.Equal(t, []string{"animal", "color", "box"}, tt.FieldsOrder)
		assert.ElementsMatch(t, []*RadioField{{
			Type:  "radio",
			Group: "animal",
			Labels: []ValueWithvisual{
				newValue("cat", "Cat"),
				newValue("dog", "Dog"),
			},
		}}, tt.Radios)

		assert.ElementsMatch(t, []*CheckboxField{{
			Type:  "checkbox",
			Group: "color",
			Labels: []ValueWithvisual{
				newValue("black", "Black"),
				newValue("white", "White"),
				newValue("pink", "Pink"),
			},
		}}, tt.Checkboxes)

		assert.ElementsMatch(t, []*BoundingBoxField{{
			Type:  "bounding_box",
			Group: "box",
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "box_kind",
					Labels: []ValueWithvisual{
						newValue("kitty", "Kitty"),
						newValue("doge", "Doge"),
					},
				}},
				Checkboxes:  make([]*CheckboxField, 0),
				FieldsOrder: []string{"box_kind"},
			},
		}}, tt.BoundingBoxes)
	}

	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="Cat" />
			<radio group="animal" value="dog" visual="Dog" />

			<checkbox group="color" value="black" visual="Black" />
			<checkbox group="color" value="white" visual="White" />
			<checkbox group="color" value="pink" visual="Pink" />

			<bounding_box group="box">
				<radio group="animal" value="kitty" visual="Kitty" />
				<radio group="animal" value="doge" visual="Doge" />
			</bounding_box>
		</content>
		`
		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate groups: animal", err.Error())
	}

	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="Cat" />
			<radio group="animal" value="dog" visual="Dog" />

			<checkbox group="animal" value="black" visual="Black" />
			<checkbox group="animal" value="white" visual="White" />
			<checkbox group="color" value="pink" visual="Pink" />
		</content>
		`
		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate groups: animal", err.Error())
	}

	{
		data := `
		<content>
			<radio group="animal" value="cat" visual="Cat" />
			<radio group="animal" value="cat" visual="Dog" />

			<checkbox group="color" value="black" visual="Black" />
			<checkbox group="color" value="white" visual="White" />
			<checkbox group="color" value="white" visual="white" />
		</content>
		`
		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate labels: group [animal] labels [cat], group [color] labels [white]", err.Error())
	}

	{
		data := `
				<radio group="animal" value="cat" visual="Cat" />
				<radio group="animal" value="cat" visual="Dog" />
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Root node <content> not found", err.Error())
	}

	{
		data := `
		<content>
			<bounding_box group="box" value="cat" visual="Cat" />
		</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "bounding_box field should have some classification fields in it", err.Error())
	}

	{
		data := `
		<content>
		    <bounding_box group="box1">
		      <checkbox group="color1" value="black" visual="Black" />
				</bounding_box>
		    <bounding_box group="box2">
		      <checkbox group="color2" value="black" visual="Black" />
		    </bounding_box>
		</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Only one bounding_box field supported for now", err.Error())
	}

	{
		data := `
		<content>
		    <bounding_box group="box">
					<checkbox group="color" value="black" visual="Black" />
					<radio group="color" value="cat" visual="Dog" />
		    </bounding_box>
		</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Template has duplicate groups: color", err.Error())
	}

	{
		data := `
		<content>
			<checkbox group="box" value="black" visual="Black" />
		</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Group name [box] is reserved for bounding_box box markup", err.Error())
	}

	{
		data := `
		<content>
		    <bounding_box group="box1">
					<checkbox group="box" value="black" visual="Black" />
		    </bounding_box>
		</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Group name [box] is reserved for bounding_box box markup", err.Error())
	}
}

func TestDuplicatedLabels(t *testing.T) {
	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:   "radio",
					Group:  "animal",
					Labels: []ValueWithvisual{newValue("dog", "Dog")},
				}},
				Checkboxes: []*CheckboxField{{
					Type:   "checkbox",
					Group:  "color",
					Labels: []ValueWithvisual{newValue("white", "White")},
				}},
			},
		}

		dups := duplicatedLabels(tmplt)
		assert.Empty(t, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "animal",
					Labels: []ValueWithvisual{
						newValue("dog", "Dog"),
						newValue("dog", "Dogh"),
					},
				}},
			},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"animal": {"dog"},
		}, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Checkboxes: []*CheckboxField{{
					Type:  "checkbox",
					Group: "color",
					Labels: []ValueWithvisual{
						newValue("white", "White"),
						newValue("white", "White"),
						newValue("black", "Black"),
						newValue("black", "Black"),
					},
				}},
			},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"color": {"black", "white"},
		}, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "animal",
					Labels: []ValueWithvisual{
						newValue("dog", "Dog"),
						newValue("dog", "Dogh"),
					},
				}},
				Checkboxes: []*CheckboxField{{
					Type:  "checkbox",
					Group: "color",
					Labels: []ValueWithvisual{
						newValue("white", "White"),
						newValue("white", "White"),
						newValue("black", "Black"),
						newValue("black", "Black"),
					},
				}},
			},
		}

		dups := duplicatedLabels(tmplt)
		assert.Equal(t, map[string][]string{
			"animal": {"dog"},
			"color":  {"black", "white"},
		}, dups)
	}
}

func TestDuplicatedGroups(t *testing.T) {
	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "animal",
				}},
				Checkboxes: []*CheckboxField{{
					Type:  "checkbox",
					Group: "color",
				}},
			},
		}

		dups := duplicatedGroups(tmplt)
		assert.Empty(t, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "animal",
				}},
				Checkboxes: []*CheckboxField{{
					Type:  "checkbox",
					Group: "animal",
				}},
			},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{"animal"}, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{{
					Type:  "radio",
					Group: "animal",
				}},
				Checkboxes: []*CheckboxField{
					{
						Type:  "checkbox",
						Group: "color",
					},
					{
						Type:  "checkbox",
						Group: "animal",
					}},
			},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{"animal"}, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios:     []*RadioField{},
				Checkboxes: []*CheckboxField{},
			},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{}, dups)
	}

	{
		tmplt := Template{
			ClassificationComponents: &ClassificationComponents{
				Radios: []*RadioField{},
				Checkboxes: []*CheckboxField{{
					Type:  "checkbox",
					Group: "color",
				}},
			},
		}

		dups := duplicatedGroups(tmplt)
		assert.ElementsMatch(t, []string{}, dups)
	}
}

func TestClassificationFieldInsideClassificationFieldError(t *testing.T) {
	{
		data := `
			<content>
				<radio group="animal" value="cat" visual="Cat">
					<radio group="animal" value="cat" visual="Cat"/>
				</radio>
				<radio group="animal" value="cat" visual="Dog"/>
			</content>
		`

		_, err := XMLToTemplate(data)

		assert.Equal(t, "Classification field [radio] mustn't have children", err.Error())
	}
}
