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

	duplicateGroups := duplicatedGroups(t)

	if len(duplicateGroups) > 0 {
		errMsg := fmt.Sprintf(
			"Template has duplicate groups: %s",
			strings.Join(duplicateGroups, ", "),
		)
		return Template{}, NewBusinessError(errMsg)
	}

	return t, nil
}

func duplicatedGroups(t Template) []string {
	groupCount := map[string]int{}

	for _, f := range t.Radios {
		groupCount[f.Group] += 1
	}

	for _, f := range t.Checkboxes {
		groupCount[f.Group] += 1
	}

	dups := make([]string, 0)
	for l, c := range groupCount {
		if c > 1 {
			dups = append(dups, l)
		}
	}

	return dups
}

func getVizual(n Node) string {
	return n.Attrs[2].Value
}

func getGroup(n Node) string {
	return n.Attrs[0].Value
}

func getValue(n Node) string {
	return n.Attrs[1].Value
}

func walkFirstLevel(nodes []Node, f func(Node) bool) {
	if len(nodes) == 0 {
		return
	}

	for _, n := range nodes[0].Nodes {
		f(n)
	}
}
