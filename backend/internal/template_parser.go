package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Node struct {
	XMLName xml.Name
	Content []byte     `xml:",innerxml"`
	Attrs   []xml.Attr `xml:"-"`
	Nodes   []Node     `xml:",any"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func missingAttribute(nodes []Node) (Node, string) {
	for _, n := range nodes {
		{
			_, ok := getAttrByName(n, "group")
			if !ok {
				return n, "group"
			}
		}

		if isClField(n.XMLName.Local) {
			{
				_, ok := getAttrByName(n, "value")
				if !ok {
					return n, "value"
				}
			}

			{
				_, ok := getAttrByName(n, "vizual")
				if !ok {
					return n, "vizual"
				}
			}
		}
	}

	return Node{}, ""
}

func emptyAttribute(nodes []Node) (Node, string) {
	for _, n := range nodes {
		{
			a, _ := getAttrByName(n, "group")
			if len(a.Value) == 0 {
				return n, "group"
			}
		}

		if isClField(n.XMLName.Local) {
			{
				a, _ := getAttrByName(n, "value")
				if len(a.Value) == 0 {
					return n, "value"
				}
			}

			{
				a, _ := getAttrByName(n, "vizual")
				if len(a.Value) == 0 {
					return n, "vizual"
				}
			}
		}
	}

	return Node{}, ""
}

func nodeToUpdatedField(t *Template, n Node) error {
	switch n.XMLName.Local {
	case "radio":
		t.CreateOrUpdateClFieldFor(n)
	case "checkbox":
		t.CreateOrUpdateClFieldFor(n)
	case "bounding_box":
		err := t.CreateOrUpdateBBoxFieldFor(n)
		if err != nil {
			return err
		}
	default:
		return NewBusinessError(
			fmt.Sprintf("Unsupported element [%s]", n.XMLName.Local),
		)
	}

	return nil
}

func NewTemplate() Template {
	return Template{
		ClassificationComponents: &ClassificationComponents{
			Radios:      make([]*RadioField, 0),
			Checkboxes:  make([]*CheckboxField, 0),
			FieldsOrder: make([]string, 0),
		},
		BoundingBoxes: make([]*BoundingBoxField, 0),
	}
}

func nodesToTemplate(nodes []Node) (Template, error) {
	t := NewTemplate()

	for _, n := range nodes {
		err := nodeToUpdatedField(&t, n)
		if err != nil {
			return Template{}, err
		}
	}

	return t, nil
}

func XMLToTemplate(s string) (Template, error) {
	buf := bytes.NewBuffer([]byte(s))
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		return Template{}, errors.WithStack(err)
	}

	nodes := []Node{}
	walkFirstLevel([]Node{n}, func(n Node) bool {
		nodes = append(nodes, n)
		return true
	})

	{
		targetNode, attr := missingAttribute(nodes)
		if len(attr) > 0 {
			errMsg := fmt.Sprintf(
				"Element [%s] missing the attribute [%s]",
				targetNode.XMLName.Local,
				attr,
			)
			return Template{}, NewBusinessError(errMsg)
		}
	}

	{
		targetNode, attr := emptyAttribute(nodes)
		if len(attr) > 0 {
			errMsg := fmt.Sprintf(
				"Element [%s] has an empty attribute [%s]",
				targetNode.XMLName.Local,
				attr,
			)
			return Template{}, NewBusinessError(errMsg)
		}
	}

	t, err := nodesToTemplate(nodes)

	if err != nil {
		return Template{}, err
	}

	{
		dups := duplicatedGroups(t)

		if len(dups) > 0 {
			errMsg := fmt.Sprintf(
				"Template has duplicate groups: %s",
				strings.Join(dups, ", "),
			)
			return Template{}, NewBusinessError(errMsg)
		}
	}
	{
		dups := []string{}
		for g, l := range duplicatedLabels(t) {
			dl := strings.Join(l, ", ")
			dups = append(dups, fmt.Sprintf("group [%s] labels [%s]", g, dl))
		}

		if len(dups) > 0 {
			sort.Strings(dups)
			errMsg := fmt.Sprintf(
				"Template has duplicate labels: %s",
				strings.Join(dups, ", "),
			)
			return Template{}, NewBusinessError(errMsg)
		}
	}

	return t, nil
}

func duplicatedGroups(t Template) []string {
	groupCount := map[string]int{}
	for _, f := range t.getClassificationFields() {
		groupCount[f.Group] += 1
	}

	for _, b := range t.BoundingBoxes {
		for _, f := range b.getClassificationFields() {
			groupCount[f.Group] += 1
		}
		groupCount[b.Group] += 1
	}

	return findCountsGt1(groupCount)
}

func findCountsGt1(m map[string]int) []string {
	dups := make([]string, 0)
	for key, c := range m {
		if c > 1 {
			dups = append(dups, key)
		}
	}

	sort.Strings(dups)

	return dups
}

func duplicatedLabels(t Template) map[string][]string {
	dups := map[string][]string{}
	for _, f := range t.getClassificationFields() {
		labelCount := map[string]int{}
		for _, l := range f.Labels {
			labelCount[l.Value] += 1
		}

		d := findCountsGt1(labelCount)
		sort.Strings(d)
		if len(d) > 0 {
			dups[f.Group] = d
		}
	}

	return dups
}

func getAttrByName(n Node, name string) (xml.Attr, bool) {
	for _, a := range n.Attrs {
		if name == a.Name.Local {
			return a, true
		}
	}

	return xml.Attr{}, false
}

func getVizual(n Node) string {
	a, _ := getAttrByName(n, "vizual")
	return a.Value
}

func getGroup(n Node) string {
	a, _ := getAttrByName(n, "group")
	return a.Value
}

func getValue(n Node) string {
	a, _ := getAttrByName(n, "value")
	return a.Value
}

func getColor(n Node) string {
	a, _ := getAttrByName(n, "color")
	return a.Value
}

func walkFirstLevel(nodes []Node, f func(Node) bool) {
	if len(nodes) == 0 {
		return
	}

	for _, n := range nodes[0].Nodes {
		f(n)
	}
}
