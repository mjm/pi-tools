package apiservice

import (
	"github.com/mjm/graphql-go"
)

type Node struct {
}

func (n *Node) ID() graphql.ID {
	return ""
}

func (n *Node) ToTrip() (*Trip, bool) {
	return nil, false
}

func (n *Node) ToLink() (*Link, bool) {
	return nil, false
}

func (n *Node) ToArchive() (*Archive, bool) {
	return nil, false
}

func (n *Node) ToPaperlessDocument() (*PaperlessDocument, bool) {
	return nil, false
}
