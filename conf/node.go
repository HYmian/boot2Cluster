package conf

import (
	"sort"
)

type Node map[string]string

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n Nodes) Less(i, j int) bool {
	return n[i]["Host"] < n[j]["Host"]
}

func (b *Boot) AddNode(node Node) error {
	if len(b.Nodes) == 0 {
		b.Nodes = append(b.Nodes, node)
		b.Entry()
	}

	i := sort.Search(len(b.Nodes), func(i int) bool {
		return b.Nodes[i]["Host"] >= node["Host"]
	})
	if i < len(b.Nodes) && b.Nodes[i]["Host"] == node["Host"] {
		b.Nodes[i]["IP"] = node["IP"]
	} else {
		b.Nodes = append(b.Nodes, node)
		sort.Sort(Nodes(b.Nodes))
	}

	return b.Entry()
}
