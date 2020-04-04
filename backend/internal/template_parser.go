package internal

import (
	"backend/pkg/utils"
	"bytes"
	"encoding/xml"
	"fmt"
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

	return Node{}, ""
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
		targetNode, missingAttr := missingAttribute(nodes)
		if len(missingAttr) > 0 {
			errMsg := fmt.Sprintf(
				"Element [%s] missing the attribute [%s]",
				targetNode.XMLName.Local,
				missingAttr,
			)
			return Template{}, NewBusinessError(errMsg)
		}
	}

	radiosByGroup := map[string]*RadioField{}
	checkboxesByGroup := map[string]*CheckboxField{}
	groupsByOrder := []string{}
	for _, n := range nodes {
		g := getGroup(n)
		groupsByOrder = append(groupsByOrder, g)

		appendLabels := func(labels []ValueWithVizual) []ValueWithVizual {
			return append(labels, ValueWithVizual{
				Vizual: getVizual(n),
				Value:  getValue(n),
			})
		}

		switch n.XMLName.Local {
		case "radio":
			f, ok := radiosByGroup[g]
			if !ok {
				f = NewRadioField(g)
				radiosByGroup[g] = f
			}

			f.Labels = appendLabels(f.Labels)
		case "checkbox":
			f, ok := checkboxesByGroup[g]
			if !ok {
				f = NewCheckboxField(g)
				checkboxesByGroup[g] = f
			}

			f.Labels = appendLabels(f.Labels)
		default:
			return Template{}, NewBusinessError(
				fmt.Sprintf("Unsupported element [%s]", n.XMLName.Local),
			)
		}
	}

	groupsByOrder = utils.Unique(groupsByOrder)

	radios := make([]RadioField, 0)
	checkboxes := make([]CheckboxField, 0)
	for _, v := range radiosByGroup {
		radios = append(radios, *v)
	}
	for _, v := range checkboxesByGroup {
		checkboxes = append(checkboxes, *v)
	}

	t := Template{
		Radios:      radios,
		Checkboxes:  checkboxes,
		FieldsOrder: groupsByOrder,
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

	return findCountsGt1(groupCount)
}

func findCountsGt1(m map[string]int) []string {
	dups := make([]string, 0)
	for key, c := range m {
		if c > 1 {
			dups = append(dups, key)
		}
	}

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

func walkFirstLevel(nodes []Node, f func(Node) bool) {
	if len(nodes) == 0 {
		return
	}

	for _, n := range nodes[0].Nodes {
		f(n)
	}
}
