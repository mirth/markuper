package internal

import (
	"backend/pkg/utils"
	"bytes"
	"encoding/xml"

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
	walk([]Node{n}, func(n Node) bool {
		if n.XMLName.Local == "radio" {
			nodes = append(nodes, n)
		}
		return true
	})

	radiosByGroups := map[string]*RadioField{}
	groupsByOrder := []string{}
	for _, n := range nodes {
		g := getGroup(n)
		groupsByOrder = append(groupsByOrder, g)

		{
			r, ok := radiosByGroups[g]
			if !ok {
				newR := NewRadioField(g)
				r = &newR
				radiosByGroups[g] = r
			}

			r.Labels = append(r.Labels, ValueWithVizual{
				Vizual: getVizual(n),
				Value:  getValue(n),
			})
		}
	}

	groupsByOrder = utils.Unique(groupsByOrder)

	radios := []RadioField{}
	for _, group := range groupsByOrder {
		r, _ := radiosByGroups[group]
		// fixme assert ok
		radios = append(radios, *r)
	}

	return Template{
		Radios: radios,
	}, nil
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

func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		f(n)
		walk(n.Nodes, f)
	}
}
