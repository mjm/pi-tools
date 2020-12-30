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
