package conf

import (
	"sort"
)

type Node struct {
	Host     string
	IP       string
	Sequence uint
}

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Host < n[j].Host
}

func (b *Boot) AddNode(host, ip string, sequence uint) error {
	if len(b.Nodes) == 0 {
		b.Nodes = append(b.Nodes, Node{Host: host, IP: ip, Sequence: sequence})
		b.Entry()
	}

	i := sort.Search(len(b.Nodes), func(i int) bool {
		return b.Nodes[i].Host >= host
	})
	if i < len(b.Nodes) && b.Nodes[i].Host == host {
		b.Nodes[i].IP = ip
	} else {
		b.Nodes = append(b.Nodes, Node{Host: host, IP: ip, Sequence: sequence})
		sort.Sort(Nodes(b.Nodes))
	}

	return b.Entry()
}
